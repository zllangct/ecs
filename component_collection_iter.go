package ecs

type ComponentCollectionIter struct {
	ls         []*componentData
	index      int
	indexInner int
}

func NewComponentCollectionIter(ls []*componentData) *ComponentCollectionIter {
	return &ComponentCollectionIter{ls: ls}
}

func (p *ComponentCollectionIter) first() IComponent {
	return p.next()
}

func (p *ComponentCollectionIter) next() IComponent {
	if p.index == len(p.ls) {
		return nil
	}
	c := p.ls[p.index].data[p.indexInner].(IComponent)
	if p.indexInner == len(p.ls[p.index].data)-1 {
		p.index += 1
		p.indexInner = 0
	} else {
		p.indexInner += 1
	}
	return c
}
