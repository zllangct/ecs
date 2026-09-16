package ecs

import "testing"

// 组表世界 Marshal → NewWorldFromData 恢复 → 重新注册系统并 Update 后：
// 组件值一致、组成员资格恢复。
func TestGroupSerialization_RoundTrip(t *testing.T) {
	w := newGroupWorld(t)
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11))) // 表内
	e2 := w.NewEntity(WithComponents(newGroupCompA(2)))                    // 残段
	e3 := w.NewEntity(WithComponents(newGroupCompB(33)))                   // 残段
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	a := groupOf(t, w)
	if _, ok := a.rowOf(e1.Index()); !ok {
		t.Fatal("e1 should be in table before marshal")
	}

	data := w.Marshal()

	// Marshal 后源世界组表应已恢复（scatter→marshal→absorb）
	if _, ok := a.rowOf(e1.Index()); !ok {
		t.Fatal("source world table should be restored after marshal")
	}

	w2 := NewWorldFromData(data)
	if err := w2.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
		t.Fatal(err)
	}
	if err := w2.Update(); err != nil {
		t.Fatal(err)
	}

	a2 := groupOf(t, w2)
	if _, ok := a2.rowOf(e1.Index()); !ok {
		t.Fatal("e1 should be re-gathered into table after unmarshal")
	}
	if _, ok := a2.rowOf(e2.Index()); ok {
		t.Fatal("e2 (partial) should stay out of table")
	}
	gotA, ok := w2.GetComponent[groupCompA, *groupCompA](e1)
	if !ok || gotA.V != 1 {
		t.Fatalf("e1.A = %+v %v", gotA, ok)
	}
	gotB, ok := w2.GetComponent[groupCompB, *groupCompB](e1)
	if !ok || gotB.V != 11 {
		t.Fatalf("e1.B = %+v %v", gotB, ok)
	}
	gotA2, ok := w2.GetComponent[groupCompA, *groupCompA](e2)
	if !ok || gotA2.V != 2 {
		t.Fatalf("e2.A = %+v %v", gotA2, ok)
	}
	gotB3, ok := w2.GetComponent[groupCompB, *groupCompB](e3)
	if !ok || gotB3.V != 33 {
		t.Fatalf("e3.B = %+v %v", gotB3, ok)
	}
}

// 关闭组表的世界也能消费组表世界的存档（导出本就是 CSet 形态）
func TestGroupSerialization_OffWorldReadsArchive(t *testing.T) {
	w := newGroupWorld(t)
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	data := w.Marshal()

	w2 := NewWorldFromData(data, WithoutGroups())
	if err := w2.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
		t.Fatal(err)
	}
	if err := w2.Update(); err != nil {
		t.Fatal(err)
	}
	gotA, ok := w2.GetComponent[groupCompA, *groupCompA](e1)
	if !ok || gotA.V != 1 {
		t.Fatalf("e1.A = %+v %v", gotA, ok)
	}
}
