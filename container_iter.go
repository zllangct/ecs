package ecs

import "unsafe"

type Iterator *iterator

type iterator struct {
	memberSize uintptr
	size       int
	index      int
	head       uintptr
}

func NewIterator(container *UnorderedContainer) Iterator {
	return &iterator{
		memberSize: container.unit,
		size:       container.len,
		index:      -1,
		head:       container.head,
	}
}

func EmptyIterator() Iterator {
	return &iterator{
		size:  0,
		index: -1,
	}
}

func (p *iterator) End() unsafe.Pointer {
	return nil
}

func (p *iterator) Next() unsafe.Pointer {
	if p.index >= p.size-1 || p.size == 0 {
		return nil
	}
	p.index++
	return unsafe.Pointer(p.head + uintptr(p.index)*p.memberSize)
}
