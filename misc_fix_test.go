package ecs

import "testing"

// D15: FixedString.Set 超长输入必须截断而非越界
func TestFixedString_SetOverflowClamped(t *testing.T) {
	type smallBuf [8]byte
	var f FixedString[smallBuf]
	f.Set("0123456789abcdef") // 16 > 8
	got := f.String()
	if got != "01234567" {
		t.Errorf("want truncated %q, got %q", "01234567", got)
	}
	if f.Len() != 8 {
		t.Errorf("len should clamp to capacity 8, got %d", f.Len())
	}
}

// D16: 查询条件含不存在的组件集 → 零匹配
func TestNewQuery_MissingSetZeroMatch(t *testing.T) {
	w := NewWorld()
	count := 0
	sys := LightSystem(func(ctx *SystemContext, event Event) error {
		q := NewQuery(ctx, WithComp[dummyComponent](), WithComp[dummyComponent2]())
		for range q.Iter() {
			count++
		}
		return nil
	})
	if err := w.RegisterLight(sys, WithDep[dummyComponent](), WithDep[dummyComponent2]()); err != nil {
		t.Fatal(err)
	}
	w.NewEntity(WithComponents(&dummyComponent{Seq: 1}))
	w.Update() // flush，dummyComponent 集合存在，dummyComponent2 集合不存在
	if count != 0 {
		t.Errorf("query with missing component set should match 0, got %d", count)
	}
}

// D34: LocalUniqueID 必须唯一且单调不减
func TestLocalUniqueID_UniqueMonotonic(t *testing.T) {
	const n = 200000
	prev := uint64(0)
	seen := make(map[uint64]struct{}, n)
	for i := 0; i < n; i++ {
		id := LocalUniqueID()
		if id < prev {
			t.Fatalf("not monotonic at %d: %d < %d", i, id, prev)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id at %d: %d", i, id)
		}
		seen[id] = struct{}{}
		prev = id
	}
}
