package ecs

import (
	"iter"
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// groupCompBView 手写只读视图（codegen 视图的等价物）
type groupCompBView struct{ b *groupCompB }

func (v groupCompBView) ComponentPacketIdentifier() rockmem.PacketIdentifier {
	return 99999999991002
}
func (v groupCompBView) FromPtr(p unsafe.Pointer) groupCompBView {
	return groupCompBView{b: (*groupCompB)(p)}
}
func (v groupCompBView) V() int32 { return v.b.V }

type groupCompBSnapshot struct {
	groupCompBView
	doubled int32
}

func (groupCompBSnapshot) FromPtr(p unsafe.Pointer) groupCompBSnapshot {
	view := groupCompBView{}.FromPtr(p)
	return groupCompBSnapshot{groupCompBView: view, doubled: 2 * view.V()}
}

func TestView2ROExactGroupIterator(t *testing.T) {
	// The two columns have different strides (16 and 4 bytes).
	RegisterComponent[PositionTestComp]("test_ext")
	w := NewWorld()
	var savedCtx *SystemContext
	var denied iter.Seq2[*PositionTestComp, groupCompBSnapshot]
	var entities []Entity
	if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		savedCtx = ctx
		if denied != nil {
			for range denied {
				t.Fatal("iterator created with an inactive context yielded after reactivation")
			}
		}

		rows := View2RO[PositionTestComp, groupCompBSnapshot](ctx)
		first := 0
		for range rows {
			first++
			break
		}
		if first != min(1, len(entities)) {
			t.Fatalf("early break visited %d rows for %d entities", first, len(entities))
		}
		for a, b := range rows {
			want := 10 * int32(a.X)
			if b.V() != want || b.doubled != 2*want {
				t.Fatalf("row X=%v: B=%v, doubled=%v, want %v and %v", a.X, b.V(), b.doubled, want, 2*want)
			}
			a.Y++
		}
		return nil
	}, WithDeps(Dep[PositionTestComp](), Dep[groupCompB](ReadOnly)),
		WithGroup(Dep[PositionTestComp](), Dep[groupCompB]())); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if savedCtx == nil {
		t.Fatal("system did not run with an empty group")
	}
	denied = View2RO[PositionTestComp, groupCompBSnapshot](savedCtx)
	unmatched := w.NewEntity(WithComponents(&PositionTestComp{X: 4}))
	for i := 1; i <= 3; i++ {
		entities = append(entities, w.NewEntity(WithComponents(&PositionTestComp{X: float32(i)}, newGroupCompB(10*int32(i)))))
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	for _, entity := range entities {
		if a, ok := w.GetComponent[PositionTestComp](entity); !ok || a.Y != 1 {
			t.Fatalf("reused iterator must write to entity %v once, got %v", entity, a)
		}
	}
	if a, ok := w.GetComponent[PositionTestComp](unmatched); !ok || a.Y != 0 {
		t.Fatalf("iterator wrote to an entity missing its peer component: %v", a)
	}
}

type viewRecorder struct {
	pair  [2]int32
	count int
}

func (rec *viewRecorder) update(ctx *SystemContext, _ Event) error {
	rec.count = 0
	for a, b := range View2[groupCompA, groupCompB](ctx) {
		rec.count++
		rec.pair = [2]int32{a.V, b.V}
		a.V++
	}
	return nil
}

// View2：只命中组表行（集齐 A,B 的实体），写穿透生效
func TestView2_Lockstep(t *testing.T) {
	w := NewWorld()
	rec := &viewRecorder{}
	err := w.RegisterLight(rec.update,
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2))) // 残段不命中
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 1 {
		t.Fatalf("View2 count = %d, want 1", rec.count)
	}
	if rec.pair != [2]int32{1, 11} {
		t.Fatalf("pair = %v", rec.pair)
	}
	// 写穿透：a.V++ 应落到表内存
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 2 {
		t.Fatalf("write-through failed: A.V = %d, want 2", got.V)
	}
}

// View2 权限：B 为只读依赖时，可写 View2 静默返回空
func TestView2_ReadonlyRejected(t *testing.T) {
	w := NewWorld()
	rec := &viewRecorder{}
	err := w.RegisterLight(rec.update,
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 0 {
		t.Fatalf("View2 with readonly B should be empty, got %d", rec.count)
	}
}

// View2 回退：未建组（WithoutGroups）时仍返回正确结果
func TestView2_FallbackNonGrouped(t *testing.T) {
	w := NewWorld(WithoutGroups())
	rec := &viewRecorder{}
	err := w.RegisterLight(rec.update,
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 1 {
		t.Fatalf("fallback count = %d, want 1", rec.count)
	}
	if rec.pair != [2]int32{1, 11} {
		t.Fatalf("fallback pair = %v", rec.pair)
	}
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 2 {
		t.Fatalf("fallback write-through failed: A.V = %d", got.V)
	}
}
