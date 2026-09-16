package ecs

import "testing"

func TestUSetRemove_OutOfRange(t *testing.T) {
	u := NewUSet[int](4)
	for i := 0; i < 4; i++ {
		v := i
		u.Add(&v)
	}
	// idx >= len 应安全返回 nil，不得触碰未使用区域
	removed, _, _ := u.Remove(10)
	if removed != nil {
		t.Errorf("want nil, got %v", *removed)
	}
	if u.Len() != 4 {
		t.Errorf("len changed: want 4, got %d", u.Len())
	}
}

func TestUSetIter_EarlyBreak(t *testing.T) {
	u := NewUSet[int](4)
	for i := 0; i < 5; i++ {
		v := i
		u.Add(&v)
	}
	count := 0
	for range u.Iter() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("want 1, got %d", count)
	}
}

func TestUSet_EleSizeWithoutInitSize(t *testing.T) {
	// 不传 initSize 时 eleSize 也必须正确，否则 Get 指针运算全错
	u := NewUSet[int]()
	for i := 0; i < 3; i++ {
		v := i * 10
		u.Add(&v)
	}
	for i := 0; i < 3; i++ {
		if got := *u.Get(int64(i)); got != i*10 {
			t.Errorf("idx %d: want %d, got %d", i, i*10, got)
		}
	}
}
