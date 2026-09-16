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
