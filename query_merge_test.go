package ecs

import "testing"

// 保序：按递增 key 添加不破坏 isKOrder；乱序添加/删除破坏
func TestSparseArray_OrderPreservation(t *testing.T) {
	s := NewSparseArray[EntityIndex, groupCompA]()
	v1, v2, v3 := &groupCompA{1}, &groupCompA{2}, &groupCompA{3}
	s.Add(1, v1)
	s.Add(3, v2)
	s.Add(7, v3)
	if !s.isKOrder {
		t.Fatal("ascending adds should keep order")
	}
	s.Add(2, &groupCompA{4}) // 乱序插入
	if s.isKOrder {
		t.Fatal("out-of-order add should break order flag")
	}
	s.Sort()
	if !s.isKOrder {
		t.Fatal("Sort should restore order flag")
	}
	s.Remove(3) // swap-remove 破坏顺序
	if s.isKOrder {
		t.Fatal("remove should break order flag")
	}
}

// merge-join：有序池查询结果正确（与 Exist 过滤路径等价）
func TestQueryMergeJoin_Ordered(t *testing.T) {
	w := NewWorld(WithoutGroups())
	rec := &queryAccessRecorder2{}
	err := w.RegisterLight(LightSystem(rec.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	// 递增创建，池天然有序：1:{A,B} 2:{A} 3:{A,B} 4:{B} 5:{A,B}
	w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(1)))
	w.NewEntity(WithComponents(newGroupCompA(2)))
	w.NewEntity(WithComponents(newGroupCompA(3), newGroupCompB(3)))
	w.NewEntity(WithComponents(newGroupCompB(4)))
	w.NewEntity(WithComponents(newGroupCompA(5), newGroupCompB(5)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	// 确认两池有序（merge 路径生效的前提）
	sa, _ := w.getComponentSet(GetIntType[groupCompA, *groupCompA]())
	sb, _ := w.getComponentSet(GetIntType[groupCompB, *groupCompB]())
	if !sa.IsKeyOrdered() || !sb.IsKeyOrdered() {
		t.Fatal("pools should be key-ordered")
	}
	if len(rec.hits) != 3 {
		t.Fatalf("hits = %v, want 3 entities", rec.hits)
	}
	for idx, v := range rec.hits {
		if v[0] != v[1] {
			t.Fatalf("entity %d: A=%d B=%d mismatch", idx, v[0], v[1])
		}
	}
}

// 乱序回退：池失序后查询结果仍正确（Exist 过滤路径）
func TestQueryMergeJoin_FallbackWhenUnordered(t *testing.T) {
	w := NewWorld(WithoutGroups())
	rec := &queryAccessRecorder2{}
	err := w.RegisterLight(LightSystem(rec.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(1)))
	w.NewEntity(WithComponents(newGroupCompA(2), newGroupCompB(2)))
	e3 := w.NewEntity(WithComponents(newGroupCompA(3), newGroupCompB(3)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	// 删除 e1 → swap-remove 破坏 A/B 池顺序 → 回退 Exist 过滤
	info, _ := w.GetEntityInfo(e1)
	info.Remove(&groupCompA{}, &groupCompB{})
	// e3 也删 B，结果集应只剩 e2
	info3, _ := w.GetEntityInfo(e3)
	info3.Remove(&groupCompB{})
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	sa, _ := w.getComponentSet(GetIntType[groupCompA, *groupCompA]())
	if sa.IsKeyOrdered() {
		t.Skip("pool still ordered, fallback path not exercised")
	}
	if len(rec.hits) != 1 {
		t.Fatalf("hits = %v, want 1 entity", rec.hits)
	}
	for _, v := range rec.hits {
		if v != [2]int32{2, 2} {
			t.Fatalf("hit = %v, want [2 2]", v)
		}
	}
}

type queryAccessRecorder2 struct {
	hits map[EntityIndex][2]int32
}

func (r *queryAccessRecorder2) Update(ctx *SystemContext, event Event) error {
	r.hits = map[EntityIndex][2]int32{}
	q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
	for idx := range q.Iter() {
		a, _ := ctx.GetBuddy[groupCompA, *groupCompA](idx)
		b, _ := ctx.GetBuddy[groupCompB, *groupCompB](idx)
		r.hits[idx] = [2]int32{a.V, b.V}
	}
	return nil
}
