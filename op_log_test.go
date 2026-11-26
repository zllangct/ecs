package ecs

import (
	"testing"
	"testing/synctest"

	rockmem "github.com/zllangct/rockmem/golang"
)

type dummyComponent2 struct {
	Seq int32
}

func (d *dummyComponent2) NewComponentSet() ComponentSet {
	return NewCSet[dummyComponent2]()
}

func (d *dummyComponent2) PacketIdentifier() rockmem.PacketIdentifier {
	return 65529
}

func (d *dummyComponent2) IsNomadic() bool {
	return false
}

func (d *dummyComponent2) IsDisposable() bool {
	return false
}

// C1: system 内通过 ctx 提交组件操作，flush 后组件与 compound 均生效
func TestOpLog_SystemCtxAddComponents(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity()

	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		ctx.AddComponents(e, &dummyComponent{Seq: 42})
		return nil
	})
	if err := w.RegisterLight(sys, WithDep[dummyComponent]()); err != nil {
		t.Fatal(err)
	}
	// 第一帧：flush 创建组件集；第二帧：system 提交，下一帧生效
	w.Update()
	w.Update()
	w.Update()

	set, ok := w.getComponentSet(GetIntTypeByComp(&dummyComponent{}))
	if !ok {
		t.Fatal("component set not created")
	}
	comp := set.Get(e)
	if comp == nil || comp.(*dummyComponent).Seq != 42 {
		t.Fatalf("component not added, got %v", comp)
	}
	info, _ := w.GetEntityInfo(e)
	if !info.compound.Exist(GetIntTypeByComp(&dummyComponent{})) {
		t.Fatal("compound not updated")
	}
}

func TestOpLog_OpOnDestroyedEntityDropped(t *testing.T) {
	w := NewWorld()
	if err := w.RegisterLight(groupProbeLight, WithDeps(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
		t.Fatal(err)
	}
	destroyed := w.NewEntity(WithComponents(newGroupCompA(1)))
	alive := w.NewEntity(WithComponents(newGroupCompA(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	for _, e := range []Entity{destroyed, alive} {
		info, ok := w.GetEntityInfo(e)
		if !ok {
			t.Fatalf("entity %v missing before flush", e)
		}
		info.Add(newGroupCompB(514))
	}
	w.DestroyEntity(destroyed)
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if got, ok := w.GetComponent[groupCompB](alive); !ok || got.V != 514 {
		t.Fatalf("live entity component = %v, want 514", got)
	}
	// An entity lookup alone would hide an orphaned row behind the stale handle.
	pool, ok := w.getComponentSet(GetIntType[groupCompB]())
	if !ok || pool.Exist(destroyed.Index()) {
		t.Fatal("component pool missing or retains a row for the destroyed entity")
	}
}

// Main-queue writers can submit more work after the synchronization point
// has detached a batch and before its cleanup returns the old nodes.
func TestOpLog_ReactivateBucketWhileBatchIsDetached(t *testing.T) {
	w := NewWorld()
	first, second := w.NewEntity(), w.NewEntity()
	queue := w.opLog
	add := func(e Entity, comp Component) {
		queue.operate(Operate{Entity: e, Comp: comp, Op: ComponentOperateAdd})
	}
	add(first, &groupCompA{1})
	add(first, &groupCompB{2})
	tasks, syncTasks, clean := queue.getOpTasks()
	if queue.HasPending() {
		t.Fatal("detached batch still pending")
	}
	add(second, &groupCompB{3})
	add(second, &groupCompA{4})
	for _, run := range tasks {
		run()
	}
	for _, run := range syncTasks {
		run()
	}
	clean()
	if !queue.HasPending() {
		t.Fatal("lost commands submitted during flush")
	}
	if _, ok := w.GetComponent[groupCompA](second); ok {
		t.Fatal("new command applied before its synchronization point")
	}
	if err := w.Step(1); err != nil {
		t.Fatal(err)
	}
	for i, e := range []Entity{first, second} {
		a, okA := w.GetComponent[groupCompA](e)
		b, okB := w.GetComponent[groupCompB](e)
		if !okA || !okB || a.V != []int32{1, 4}[i] || b.V != []int32{2, 3}[i] {
			t.Fatalf("entity %v components = %v, %v", e, a, b)
		}
	}
	if queue.HasPending() {
		t.Fatal("drained queue still pending")
	}
}

// C1: 并行模式下多 system 并发经各自队列提交，synctest 提供确定性调度
func TestOpLog_ParallelSubmitNoRace(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		w := NewWorld(WithWorldASyncMode())
		var entities []Entity
		for i := 0; i < 8; i++ {
			entities = append(entities, w.NewEntity())
		}

		sysA := LightSystem(func(ctx *SystemContext, event Event) error {
			for _, e := range entities {
				ctx.AddComponents(e, &dummyComponent{Seq: 1})
			}
			return nil
		})
		sysB := LightSystem(func(ctx *SystemContext, event Event) error {
			for _, e := range entities {
				ctx.AddComponents(e, &dummyComponent2{Seq: 2})
			}
			return nil
		})
		if err := w.RegisterLight(sysA, WithDep[dummyComponent]()); err != nil {
			t.Fatal(err)
		}
		if err := w.RegisterLight(sysB, WithDep[dummyComponent2]()); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 3; i++ {
			if err := w.Update(); err != nil {
				t.Fatal(err)
			}
		}
		synctest.Wait()

		ww := w
		for _, e := range entities {
			info, ok := ww.getEntityInfo(e)
			if !ok {
				t.Fatalf("entity %d not found", e.Index())
			}
			if !info.compound.Exist(GetIntTypeByComp(&dummyComponent{})) {
				t.Fatalf("entity %d missing dummyComponent", e.Index())
			}
			if !info.compound.Exist(GetIntTypeByComp(&dummyComponent2{})) {
				t.Fatalf("entity %d missing dummyComponent2", e.Index())
			}
		}
	})
}
