package ecs

import (
	"testing"

	rockmem "github.com/zllangct/rockmem/golang"
)

// 游牧组件：每帧末自动清空
type nomadicComp struct {
	Seq int32
}

func (d *nomadicComp) NewComponentSet() ComponentSet {
	return NewCSet[nomadicComp]()
}

func (d *nomadicComp) PacketIdentifier() rockmem.PacketIdentifier {
	return 65528
}

func (d *nomadicComp) IsNomadic() bool {
	return true
}

func (d *nomadicComp) IsDisposable() bool {
	return false
}

// D24: DestroyEntity 延迟到帧同步点，回收 ID、清理所有组件与 compound
func TestDestroyEntity_RecyclesIDAndComponents(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity(WithComponents(&dummyComponent{Seq: 1}))
	w.Update() // flush，组件生效

	ww := w
	w.DestroyEntity(e)
	// 帧同步点前尚未生效
	if _, ok := ww.getEntityInfo(e); !ok {
		t.Fatal("entity should still exist before sync point")
	}
	w.Update() // 帧同步点销毁

	if _, ok := ww.getEntityInfo(e); ok {
		t.Error("entity should be destroyed")
	}
	set, ok := ww.getComponentSet(GetIntTypeByComp(&dummyComponent{}))
	if ok && set.Get(e) != nil {
		t.Error("component should be removed")
	}
	// ID 复用需凑满 delayCap 才归并 freelist：补满 10 个释放触发 flush
	var extra []Entity
	for i := 0; i < 9; i++ {
		extra = append(extra, w.NewEntity())
	}
	for _, x := range extra {
		w.DestroyEntity(x)
	}
	w.Update() // 帧同步点销毁 9 个 + 触发 delayFlush（1+9=10=delayCap）

	e2 := w.NewEntity()
	if e2.Index() != e.Index() {
		t.Errorf("index not recycled: want %d, got %d", e.Index(), e2.Index())
	}
	if e2.toReuseID().reuse != e.toReuseID().reuse+1 {
		t.Errorf("reuse should increment: want %d, got %d", e.toReuseID().reuse+1, e2.toReuseID().reuse)
	}
	// 旧句柄不得再命中
	if _, ok := ww.getEntityInfo(e); ok {
		t.Error("stale handle should not resolve after recycle")
	}
}

// D24: EntityInfo.Remove 移除组件
func TestEntityInfoRemove_Component(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity(WithComponents(&dummyComponent{Seq: 3}))
	w.Update()

	ww := w
	info, _ := ww.getEntityInfo(e)
	info.Remove(&dummyComponent{})
	w.Update() // flush

	set, _ := ww.getComponentSet(GetIntTypeByComp(&dummyComponent{}))
	if set.Get(e) != nil {
		t.Error("component should be removed")
	}
	info, _ = ww.getEntityInfo(e)
	if info.compound.Exist(GetIntTypeByComp(&dummyComponent{})) {
		t.Error("compound should not contain removed type")
	}
}

// D23: World.Destroy 触发 Destroy stage 链并置停
func TestWorldDestroy_TriggersDestroyStages(t *testing.T) {
	w := NewWorld()
	destroyCalled := false
	syncAfterDestroyOk := false
	destroyStageOnDestroy = func() { destroyCalled = true }
	destroyStageOnSyncAfterDestroy = func() { syncAfterDestroyOk = true }
	defer func() {
		destroyStageOnDestroy = nil
		destroyStageOnSyncAfterDestroy = nil
	}()
	if err := w.Register[destroyStageSys](); err != nil {
		t.Fatal(err)
	}
	if err := w.Destroy(); err != nil {
		t.Fatal(err)
	}
	if !destroyCalled {
		t.Error("Destroy stage not called")
	}
	if !syncAfterDestroyOk {
		t.Error("SyncAfterDestroy stage not called")
	}
	if w.getStatus() != WorldStatusStop {
		t.Error("world status should be Stop")
	}
	if err := w.Update(); err == nil {
		t.Error("Update after Destroy should return error")
	}
}

// D24: nomadic 组件每帧末清空
func TestClearNomadic_PerFrame(t *testing.T) {
	w := NewWorld()
	w.AddNomadic(&nomadicComp{Seq: 1})
	w.Update()
	w.Update() // 第二帧末应已清空

	ww := w
	set, ok := ww.getComponentSet(GetIntTypeByComp(&nomadicComp{}))
	if !ok {
		t.Fatal("nomadic set should exist")
	}
	if set.Len() != 0 {
		t.Errorf("nomadic set should be cleared per frame, len=%d", set.Len())
	}
}
