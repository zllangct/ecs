package ecs

import "unsafe"

type Iterator struct {
	memberSize uintptr
	size       int
	index      int
	head       uintptr
}

func NewIterator(container *Container) *Iterator {
	return &Iterator{
		memberSize: container.unit,
		size:       container.len,
		index:      -1,
		head:       container.head,
	}
}

func EmptyIterator() *Iterator {
	return &Iterator{
		size:  0,
		index: -1,
	}
}

func (p *Iterator) End() unsafe.Pointer {
	return nil
}

func (p *Iterator) Next() unsafe.Pointer {
	if p.index >= p.size || p.size == 0 {
		return nil
	}
	p.index++
	return unsafe.Pointer(p.head + uintptr(p.index)*p.memberSize)
}

func (p *Iterator) NextIV() (int, unsafe.Pointer) {
	if p.index >= p.size || p.size == 0 {
		return -1, nil
	}
	p.index++
	return p.index, unsafe.Pointer(p.head + uintptr(p.index)*p.memberSize)
}
