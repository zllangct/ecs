package ecs

import (
	"bytes"
	"slices"
	"testing"
	"testing/synctest"
	"time"

	rockmem "github.com/zllangct/rockmem/golang"
)

func TestRestoreKeepsRegisteredCommandOrder(t *testing.T) {
	source := NewWorld()
	entity := source.NewEntity()
	snapshot, err := source.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []WorldOption{WithWorldSyncMode(), WithWorldASyncMode()} {
		w := NewWorld(mode)
		// 两个回调按 Order 反向执行，同一组件的延迟 Add 仍应按注册序合并。
		for i, value := range []int32{11, 22} {
			if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
				ctx.AddComponents(entity, &p2CompTag{value})
				return nil
			}, WithOrder(Order(2-i)), WithDep[p2CompTag]()); err != nil {
				t.Fatal(err)
			}
		}
		w.NewEntity(WithComponents(&p2CompA{99}))
		if err := w.Unmarshal(snapshot); err != nil {
			t.Fatal(err)
		}
		if err := w.Step(time.Millisecond); err != nil {
			t.Fatal(err)
		}
		tag, ok := w.GetComponent[p2CompTag](entity)
		if !ok || tag.V != 11 {
			t.Fatalf("first registered writer should win after restore, got %v", tag)
		}
		if _, ok := w.GetComponent[p2CompA](entity); ok {
			t.Fatal("pre-restore pending command survived")
		}
	}
}

func TestFixedStepSkipsWallClockStorageOptimization(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		iterate := func(fixed bool, pause time.Duration) []EntityIndex {
			w := NewWorld(WithWorldAutoOptimize())
			var order []EntityIndex
			err := w.RegisterLight(func(ctx *SystemContext, e Event) error {
				if e.Frame == 2 {
					time.Sleep(pause)
				}
				if e.Frame == 3 {
					for idx := range ctx.GetComponents[p2CompA]() {
						order = append(order, idx)
					}
				}
				return nil
			}, WithDep[p2CompA]())
			if err != nil {
				t.Fatal(err)
			}
			step := w.Update
			if fixed {
				step = func() error { return w.Step(time.Millisecond) }
			}
			var ids []Entity
			for range 4 {
				ids = append(ids, w.NewEntity(WithComponents(&p2CompA{})))
			}
			if err := step(); err != nil {
				t.Fatal(err)
			}
			w.DestroyEntity(ids[1])
			if err := step(); err != nil {
				t.Fatal(err)
			}
			if err := step(); err != nil {
				t.Fatal(err)
			}
			return order
		}
		a, b := iterate(true, 0), iterate(true, time.Second)
		if !slices.Equal(a, []EntityIndex{1, 4, 3}) || !slices.Equal(a, b) {
			t.Fatalf("fixed-step iteration changed with wall-clock delay: %v and %v", a, b)
		}
		if order := iterate(false, 0); !slices.Equal(order, []EntityIndex{1, 3, 4}) {
			t.Fatalf("Update should still honor its optimization budget, got %v", order)
		}
	})
}

func TestCheckpointKeepsMergedGroupAndResidualOrder(t *testing.T) {
	const delta = 16 * time.Millisecond
	register := func(w *World, rows *[]int32) {
		t.Helper()
		err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
			*rows = (*rows)[:0]
			for a, b := range View2[p2CompA, p2CompB](ctx) {
				*rows = append(*rows, a.V)
				b.V += int32(len(*rows))
			}
			return nil
		}, WithDeps(Dep[p2CompA](), Dep[p2CompB]()),
			WithGroup(Dep[p2CompA](), Dep[p2CompB]()),
			WithGroup(Dep[p2CompB](), Dep[p2CompC]()))
		if err != nil {
			t.Fatal(err)
		}
	}
	encode := func(w *World) []byte {
		t.Helper()
		b := rockmem.NewWriter()
		if _, err := w.MarshalTo(b); err != nil {
			t.Fatal(err)
		}
		return b.Bytes()
	}
	original := NewWorld()
	var originalRows, restoredRows []int32
	register(original, &originalRows)
	var ids []Entity
	for i := int32(1); i <= 10; i++ {
		comps := []Component{&p2CompA{i}}
		if i != 9 {
			comps = append(comps, &p2CompB{i})
		}
		if i <= 5 || i == 9 {
			comps = append(comps, &p2CompC{i})
		}
		ids = append(ids, original.NewEntity(WithComponents(comps...)))
	}
	if err := original.Step(delta); err != nil {
		t.Fatal(err)
	}
	original.DestroyEntity(ids[1])
	if err := original.Step(delta); err != nil {
		t.Fatal(err)
	}
	if len(originalRows) != 8 || slices.IsSorted(originalRows[:4]) {
		t.Fatalf("expected four out-of-order group rows and four residual matches, got %v", originalRows)
	}

	checkpoint := encode(original)
	t.Run("regroup", func(t *testing.T) {
		other := NewWorld()
		if err := other.UnmarshalFrom(rockmem.NewReader(checkpoint)); err != nil {
			t.Fatal(err)
		}
		if err := other.RegisterLight(groupProbeLight, WithGroup(Dep[p2CompA](), Dep[p2CompC]())); err != nil {
			t.Fatal(err)
		}
		if err := other.Step(delta); err != nil {
			t.Fatal(err)
		}
		// 9 号实体原来缺 B 而位于残段，现在具备新组 {A,C} 的全部组件。
		g, _ := other.archetypes.owner(GetIntType[p2CompA, *p2CompA]())
		if _, inGroup := g.rowOf(ids[8].Index()); !inGroup {
			t.Fatal("restored entity failed to join the new group")
		}
		a, hasA := other.GetComponent[p2CompA](ids[8])
		c, hasC := other.GetComponent[p2CompC](ids[8])
		if !hasA || !hasC || a.V != 9 || c.V != 9 {
			t.Fatal("regrouped component values changed")
		}
	})

	restored := NewWorld()
	if err := restored.UnmarshalFrom(rockmem.NewReader(checkpoint)); err != nil {
		t.Fatal(err)
	}
	register(restored, &restoredRows)
	for frame := 0; frame < 3; frame++ {
		for _, w := range []*World{original, restored} {
			if frame == 1 {
				// 将一个残段提升为组行，会改变残段和行序列。
				e, _ := w.GetEntityInfo(ids[6])
				e.Add(&p2CompC{7})
			}
			if frame == 2 {
				for _, id := range ids {
					if e, ok := w.GetEntityInfo(id); ok {
						e.Remove(&p2CompC{})
					}
				}
			}
			if err := w.Step(delta); err != nil {
				t.Fatal(err)
			}
		}
		if !slices.Equal(originalRows, restoredRows) {
			t.Fatalf("frame %d: view order changed after restore: %v vs %v", frame, originalRows, restoredRows)
		}
		if !bytes.Equal(encode(original), encode(restored)) {
			t.Fatalf("frame %d: resumed values or serialized state diverged", frame)
		}
	}
}

// 恢复一份含 7500 个完整组行及 2500 个 A 残段的 V1 快照，计入解码、系统装配与首帧重建。
func BenchmarkResumeGroupedFirstStep(b *testing.B) {
	const delta = 16 * time.Millisecond
	register := func(w *World) {
		if err := w.RegisterLight(groupProbeLight, WithGroup(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
			b.Fatal(err)
		}
	}
	source := NewWorld()
	register(source)
	for i := 1; i <= 10000; i++ {
		comps := []Component{newGroupCompA(int32(i))}
		if i%4 != 0 {
			comps = append(comps, newGroupCompB(int32(-i)))
		}
		source.NewEntity(WithComponents(comps...))
	}
	if err := source.Step(delta); err != nil {
		b.Fatal(err)
	}
	source.DestroyEntity(2)
	source.NewEntity(WithComponents(newGroupCompA(10001), newGroupCompB(-10001)))
	if err := source.Step(delta); err != nil {
		b.Fatal(err)
	}
	buf := rockmem.NewWriter()
	if _, err := source.MarshalTo(buf); err != nil {
		b.Fatal(err)
	}
	blob := buf.Bytes()
	var restored *World
	b.ReportAllocs()
	b.SetBytes(int64(len(blob)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		restored = NewWorld()
		if err := restored.UnmarshalFrom(rockmem.NewReader(blob)); err != nil {
			b.Fatal(err)
		}
		register(restored)
		if err := restored.Step(delta); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	it := GetIntType[groupCompA, *groupCompA]()
	before, _ := source.archetypes.owner(it)
	after, _ := restored.archetypes.owner(it)
	if !slices.Equal(before.entities, after.entities) ||
		!slices.Equal(source.components[it].EntityIndexes(), restored.components[it].EntityIndexes()) {
		b.Fatal("restored component iteration order differs")
	}
}
