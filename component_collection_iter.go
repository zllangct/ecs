package ecs

import "unsafe"

type ComponentCollectionIter = *componentCollectionIter

type componentCollectionIter struct {
	ls         []*ContainerWithId
	index      int
	indexInner int
	temp       IComponent
}

func NewComponentCollectionIter(ls []*ContainerWithId) ComponentCollectionIter {
	return &componentCollectionIter{ls: ls}
}

func (p *componentCollectionIter) End() *IComponent {
	return nil
}

func (p *componentCollectionIter) Next() *IComponent {
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
	efaceStruct := (*eface)(unsafe.Pointer(&p.temp))
	efaceStruct.data = c
	return &p.temp
}
