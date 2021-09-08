package ecs

import "unsafe"

type ComponentCollectionIter = *componentCollectionIter

type componentCollectionIter[T any] struct {
	ls         []*ContainerWithId[T]
	index      int
	indexInner int
	len        int
	temp       IComponent
}

func NewComponentCollectionIter[T any](ls []*ContainerWithId[T]) ComponentCollectionIter {
	return &componentCollectionIter{
		ls:         ls,
		index:      0,
		indexInner: -1,
		len:        len(ls),
		temp:       &ComponentBase{},
	}
}

func (p *componentCollectionIter) End() IComponent {
	return nil
}

func (p *componentCollectionIter) Next() IComponent {
	if p.len == 0 {
		return nil
	}
	if p.indexInner == p.ls[p.index].Len()-1 {
		p.index += 1
		p.indexInner = 0
	} else {
		p.indexInner += 1
	}
	if p.index == p.len {
		p.temp = nil
		return nil
	}
	c := p.ls[p.index].Get(p.indexInner)
	efaceStruct := (*iface)(unsafe.Pointer(&p.temp))
	efaceStruct.data = c
	return p.temp
}
