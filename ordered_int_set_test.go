package ecs

import "testing"

func TestOrderedIntSet_Add(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4}
	for _, it := range insert {
		c.Add(it)
	}

	want := []uint16{2, 3, 4, 6, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}
}

func TestOrderedIntSet_Remove(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4}
	for _, it := range insert {
		c.Add(it)
	}

	c.Remove(3)

	c.Add(1)

	want := []uint16{1, 2, 4, 6, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}
}

func TestOrderedIntSet_InsertIndex(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4}
	for _, it := range insert {
		c.Add(it)
	}

	want := []uint16{2, 3, 4, 6, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}

	wantIndex := 3
	if got := c.FindIndexToInsert(5, 0); got != wantIndex {
		t.Errorf("insertIndex() = %v, want %v", got, wantIndex)
	}
}

// 回归：重复元素位于二分循环结束后的 l==r 位置时（如向 {1,2} 再 Add 2），
// FindIndexToInsert 必须返回 -1，不得返回插入位导致乱序重复。
func TestOrderedIntSet_AddDuplicateTail(t *testing.T) {
	c := OrderedIntSet[uint16]{1, 2}
	if got := c.FindIndexToInsert(2, 0); got != -1 {
		t.Errorf("FindIndexToInsert(dup tail) = %v, want -1", got)
	}
	c2 := OrderedIntSet[uint16]{1, 2}
	if c2.Add(2) {
		t.Errorf("Add(dup tail) should return false, got %v", c2)
	}
	if len(c2) != 2 || c2[0] != 1 || c2[1] != 2 {
		t.Errorf("set corrupted: %v", c2)
	}
	// 头部/中部重复也不得插入
	c3 := OrderedIntSet[uint16]{1, 2, 3}
	if c3.Add(1) || c3.Add(2) || c3.Add(3) {
		t.Errorf("dup adds should all fail, got %v", c3)
	}
}

func TestOrderedIntSet_Find(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4}
	for _, it := range insert {
		c.Add(it)
	}

	want := []uint16{2, 3, 4, 6, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}

	wantIndex := 4
	if got, _ := c.Find(7); got != wantIndex {
		t.Errorf("Find() = %v, want %v", got, wantIndex)
	}
}

func TestOrderedIntSet_IsSubSet(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4}
	for _, it := range insert {
		c.Add(it)
	}

	want := []uint16{2, 3, 4, 6, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}

	subSet := []uint16{3, 4, 6}
	wantBool := true
	if got := c.IsSubSet(subSet); got != wantBool {
		t.Errorf("IsSubSet() = %v, want %v", got, wantBool)
	}

	subSet = []uint16{2, 3, 4, 6, 7, 9}
	wantBool = true
	if got := c.IsSubSet(subSet); got != wantBool {
		t.Errorf("IsSubSet() = %v, want %v", got, wantBool)
	}

	subSet = []uint16{3, 4, 8}
	wantBool = false
	if got := c.IsSubSet(subSet); got != wantBool {
		t.Errorf("IsSubSet() = %v, want %v", got, wantBool)
	}
}

func TestOrderedIntSet_Merge(t *testing.T) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3}
	for _, it := range insert {
		c.Add(it)
	}

	c2 := OrderedIntSet[uint16]{}
	insert2 := []uint16{4, 9}
	for _, it := range insert2 {
		c2.Add(it)
	}

	c.Merge(c2)

	want := [4]uint16{3, 4, 7, 9}
	for i := 0; i < len(c); i++ {
		if c[i] != want[i] {
			t.Errorf("c[%d] = %d, want %d", i, c[i], want[i])
		}
	}
}

func BenchmarkSubSet(b *testing.B) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	for _, it := range insert {
		c.Add(it)
	}
	subSet := []uint16{3, 4, 10}
	b.Run("1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			c.IsSubSet(subSet)
		}
	})
}

func BenchmarkFind(b *testing.B) {
	c := OrderedIntSet[uint16]{}
	insert := []uint16{7, 3, 6, 2, 9, 4, 11, 12, 17, 18, 19, 26, 28, 30}
	for _, it := range insert {
		c.Add(it)
	}

	upper := 50
	b.Run("0", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			c.Find(uint16(i % upper))
		}
	})
}
