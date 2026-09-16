package ecs

import "sync"

// Pool 是 sync.Pool 的泛型封装，Get/Put 免类型断言，可复用于任意指针类型。
type Pool[T any] struct {
	pool  sync.Pool
	newFn func() *T
}

// NewPool 创建对象池，newFn 用于在池空时构造新实例。
func NewPool[T any](newFn func() *T) *Pool[T] {
	p := &Pool[T]{newFn: newFn}
	p.pool.New = func() any {
		return newFn()
	}
	return p
}

func (p *Pool[T]) Get() *T {
	v := p.pool.Get()
	if v == nil {
		return p.newFn()
	}
	return v.(*T)
}

func (p *Pool[T]) Put(t *T) {
	p.pool.Put(t)
}
