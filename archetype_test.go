package ecs

import (
	"testing"
	"unsafe"
)

func newTestArchetype() *Archetype {
	types := []ComponentIntType{1001, 1002}
	sizes := []int{16, 8}
	return NewArchetype(types, sizes)
}

func TestArchetype_RowLifecycle(t *testing.T) {
	a := newTestArchetype()
	r0 := a.addRow(10)
	r1 := a.addRow(20)
	if r0 != 0 || r1 != 1 || a.Len() != 2 {
		t.Fatalf("rows: %d %d len=%d", r0, r1, a.Len())
	}
	*(*float32)(a.cellPtr(0, r0)) = 1.5
	if got := *(*float32)(a.cellPtr(0, r0)); got != 1.5 {
		t.Fatalf("cell = %v", got)
	}
	if _, ok := a.rowOf(10); !ok {
		t.Fatal("rowOf(10) missing")
	}
	// swap-remove 行0：实体20 顶上来
	a.removeRow(10)
	if _, ok := a.rowOf(10); ok {
		t.Fatal("rowOf(10) should be gone")
	}
	row, ok := a.rowOf(20)
	if !ok || row != 0 {
		t.Fatalf("entity 20 should move to row 0, got %d %v", row, ok)
	}
	if a.entities[0] != 20 || a.Len() != 1 {
		t.Fatal("entities not maintained")
	}
}

func TestArchetype_RemoveLastRow(t *testing.T) {
	a := newTestArchetype()
	a.addRow(10)
	a.addRow(20)
	a.removeRow(20) // 删末行
	if a.Len() != 1 {
		t.Fatalf("len = %d, want 1", a.Len())
	}
	if _, ok := a.rowOf(20); ok {
		t.Fatal("rowOf(20) should be gone")
	}
	if row, ok := a.rowOf(10); !ok || row != 0 {
		t.Fatalf("entity 10 row = %d %v", row, ok)
	}
	if len(a.columns[0]) != 16 || len(a.columns[1]) != 8 {
		t.Fatalf("columns not shrunk: %d %d", len(a.columns[0]), len(a.columns[1]))
	}
}

func TestArchetype_RemoveAbsent(t *testing.T) {
	a := newTestArchetype()
	a.addRow(10)
	if got := a.removeRow(999); got != -1 {
		t.Fatalf("removeRow(999) = %d, want -1", got)
	}
	if a.Len() != 1 {
		t.Fatal("len changed")
	}
}

func TestArchetype_CellPtrAlignment(t *testing.T) {
	a := newTestArchetype()
	for i := 0; i < 100; i++ {
		a.addRow(EntityIndex(i + 1))
	}
	for row := 0; row < a.Len(); row++ {
		p := uintptr(a.cellPtr(0, row))
		if p%unsafe.Alignof(float64(0)) != 0 {
			t.Fatalf("row %d misaligned: %x", row, p)
		}
	}
}

func TestArchetype_CompoundAndColOf(t *testing.T) {
	a := newTestArchetype()
	c := a.compound()
	if len(c) != 2 || !c.Exist(1001) || !c.Exist(1002) {
		t.Fatalf("compound = %v", c)
	}
	if col, ok := a.colOf(1002); !ok || col != 1 {
		t.Fatalf("colOf(1002) = %d %v", col, ok)
	}
	if _, ok := a.colOf(9999); ok {
		t.Fatal("colOf(9999) should miss")
	}
}
