package ecs

import "testing"

func newGroupWorld(t *testing.T) *World {
	t.Helper()
	w := NewWorld()
	err := w.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

func groupOf(t *testing.T, w *World) *Archetype {
	t.Helper()
	a, ok := w.archetypes.owner(GetIntType[groupCompA, *groupCompA]())
	if !ok {
		t.Fatal("group not built")
	}
	return a
}

// 创建时集齐组组件 → 首次 flush 后入行，CSet 残段被收走
func TestGroupSync_GatherOnCreate(t *testing.T) {
	w := newGroupWorld(t)
	e := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	a := groupOf(t, w)
	if _, ok := a.rowOf(e.Index()); !ok {
		t.Fatal("entity should be in group table")
	}
	if s, ok := w.getComponentSet(GetIntType[groupCompA, *groupCompA]()); ok && s.Get(e) != nil {
		t.Fatal("A should be gathered out of CSet")
	}
	// 表内值正确
	got, ok := w.GetComponent[groupCompB, *groupCompB](e)
	if !ok || got.V != 2 {
		t.Fatalf("B = %+v %v", got, ok)
	}
}

// 部分组件 → 残段；集齐 → 入行；删除 → 拆行散回
func TestGroupSync_GatherAndScatter(t *testing.T) {
	w := newGroupWorld(t)
	e := w.NewEntity()
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	a := groupOf(t, w)

	info, _ := w.GetEntityInfo(e)
	info.Add(newGroupCompA(1)) // 只加 A → 残段，不入行
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if _, ok := a.rowOf(e.Index()); ok {
		t.Fatal("partial member should not be in table")
	}

	info.Add(newGroupCompB(2)) // 集齐 → 入行
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if _, ok := a.rowOf(e.Index()); !ok {
		t.Fatal("full member should be in table")
	}

	info.Remove(newGroupCompB(2)) // 拆行 → A 散回 CSet
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if _, ok := a.rowOf(e.Index()); ok {
		t.Fatal("should be out of table")
	}
	s, _ := w.getComponentSet(GetIntType[groupCompA, *groupCompA]())
	got := s.Get(e)
	if got == nil || got.(*groupCompA).V != 1 {
		t.Fatalf("A should be scattered back with value 1, got %v", got)
	}
}

// 销毁表内实体 → 行被回收
func TestGroupSync_DestroyInTable(t *testing.T) {
	w := newGroupWorld(t)
	e := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	a := groupOf(t, w)
	if a.Len() != 1 {
		t.Fatalf("rows = %d, want 1", a.Len())
	}
	w.DestroyEntity(e)
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if a.Len() != 0 {
		t.Fatalf("rows = %d, want 0", a.Len())
	}
}
