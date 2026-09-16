package ecs

import (
	"testing"
	"unsafe"
)

func TestCSet_AddRaw(t *testing.T) {
	set := NewCSet[PositionTestComp]()
	src := &PositionTestComp{X: 1, Y: 2, Z: 3}
	rid := ReuseID{index: 5, reuse: 0}
	e := rid.ToEntity()
	set.AddRaw(e, unsafe.Pointer(src))
	got := set.Get(e)
	if got == nil {
		t.Fatal("missing")
	}
	c := got.(*PositionTestComp)
	if c.X != 1 || c.Y != 2 || c.Z != 3 {
		t.Fatalf("got %+v", c)
	}
	// 源数据后续修改不影响集内副本（拷贝语义）
	src.X = 99
	if got.(*PositionTestComp).X != 1 {
		t.Fatal("AddRaw must copy, not alias")
	}
}
