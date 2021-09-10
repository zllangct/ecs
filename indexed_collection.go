package ecs

import "unsafe"

type IndexedCollection[T any] struct {
	c   IContainer[T]
	ids map[int64]int64
}

func NewContainerWithId[T any]() *IndexedCollection[T] {
	size := unsafe.Sizeof(*new(T))
	return &IndexedCollection[T]{
		c:   NewCollection(size),
		ids: map[int64]int64{},
	}
}

func (p *IndexedCollection[T]) Add(item *T, id ...int64) (int, *T) {
	if len(id) > 0 {
		_, ok := p.ids[id[0]]
		if ok {
			return -1, nil
		}
	}
	idx, ptr := p.c.Add(item)
	if len(id) > 0 {
		p.ids[id[0]] = int64(idx)
		p.ids[-int64(idx)] = id[0]
	}
	return idx, ptr
}

func (p *IndexedCollection[T]) remove(idx int) {
	len := p.c.Len()
	if idx < 0 || idx >= len {
		return
	}
	p.ids[p.ids[int64(len)]] = int64(idx)
	delete(p.ids, p.ids[-int64(idx)])
	p.ids[-int64(idx)] = p.ids[int64(len)]
	delete(p.ids, int64(len))

	p.c.Remove(idx)
}

func (p *IndexedCollection[T]) Remove(id int64) {
	idx, ok := p.ids[id]
	if !ok {
		return
	}
	p.remove(int(idx))
}

func (p *IndexedCollection[T]) get(idx int) *T {
	return p.c.Get(idx)
}

func (p *IndexedCollection[T]) Get(id int64) *T {
	idx, ok := p.ids[id]
	if !ok {
		return nil
	}
	return p.get(int(idx))
}

func (p *IndexedCollection[T]) Len() int {
	return p.c.Len()
}
