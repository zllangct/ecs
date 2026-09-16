package ecs

import "testing"

// Pool[T] 是 sync.Pool 的泛型封装：Get/Put 免断言，可复用于任意类型。
func TestPool_GetPut(t *testing.T) {
	p := NewPool(func() *opTask { return new(opTask) })

	a := p.Get()
	if a == nil {
		t.Fatal("Get on empty pool should return a new instance")
	}
	a.op = ComponentOperateAdd

	p.Put(a)

	b := p.Get()
	if b == nil {
		t.Fatal("Get after Put should return an instance")
	}
}

func TestPool_GenericOverTypes(t *testing.T) {
	// 同一套 Pool 实现服务两种不同类型，证明泛型复用
	p1 := NewPool(func() *opTask { return new(opTask) })
	p2 := NewPool(func() *dummyComponent { return new(dummyComponent) })

	if p1.Get() == nil || p2.Get() == nil {
		t.Fatal("both pools should yield non-nil instances")
	}
}
