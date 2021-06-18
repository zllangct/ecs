package ecs

import "unsafe"

type IComponentContainer interface {
	
}

type ContainerWithId[T any] struct {
	c	IContainer[T]
	ids	map[int64]int
}

func NewContainerWithIdByte[T any]() *ContainerWithId {
	return &ContainerWithId{
		c : NewUnorderedContainerByte[T](),
		ids: map[int64]int{},
	}
}

func NewContainerWithId[T any]() *ContainerWithId {
	return &ContainerWithId{
		c : NewUnorderedContainer[T](),
		ids: map[int64]int{},
	}
}

func (p *ContainerWithId[T]) Add(item T, id ...uint64) (int, *T) {
	if len(id) > 0 {
		_, ok := p.ids[id[0]]
		if ok {
			return -1, nil
		}
	}
	idx, ptr := p.c.Add(item)
	if len(id) > 0 {
		p.ids[id[0]] = idx	
		p.ids[-idx] = id[0]
	}
	return idx, ptr
}

func (p *ContainerWithId[T]) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	p.ids[p.ids[p.Len()]] = idx
	delete(p.ids, p.ids[-idx])
	p.ids[-idx] = p.ids[p.Len()]
	delete(p.ids, p.Len())

	p.c.Remove(idx)
}

func (p *ContainerWithId[T]) RemoveById(id uint64) {
	idx, ok := p.ids[id]
	if !ok {
		return
	}
	p.Remove(idx)
}

func (p *ContainerWithId[T]) Get(idx int) *T {
	return p.c.Get(idx)
}

func (p *ContainerWithId[T]) GetById(id uint64) *T {
	idx, ok := p.ids[id]
	if !ok {
		return nil
	}
	return p.c.Get(idx)
}

func (p *ContainerWithId[T]) GetId(idx int) uint64 {
	if id, ok := p.ids[-idx]; ok {
		return id
	}
	return 0
}

func (p *ContainerWithId[T]) Len() int {
	return p.c.Len()
}
