package ecs

import "unsafe"

type ComponentCollectionIter struct {
	ls         []*ContainerWithId
	index      int
	indexInner int
}

func NewComponentCollectionIter(ls []*ContainerWithId) *ComponentCollectionIter {
	return &ComponentCollectionIter{ls: ls}
}

func (p *ComponentCollectionIter) End() unsafe.Pointer {
	return nil
}

func (p *ComponentCollectionIter) Next() unsafe.Pointer {
	if p.index == len(p.ls) {
		return nil
	}
	c := p.ls[p.index].Get(p.indexInner)
	if p.indexInner == p.ls[p.index].Len()-1 {
		p.index += 1
		p.indexInner = 0
	} else {
		p.indexInner += 1
	}
	return c
}
