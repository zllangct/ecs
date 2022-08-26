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
	if got := c.InsertIndex(5); got != wantIndex {
		t.Errorf("insertIndex() = %v, want %v", got, wantIndex)
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
	if got := c.Find(7); got != wantIndex {
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
