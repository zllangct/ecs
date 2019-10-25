package main

type ComponentCollectionIter struct {
	ls []*componentData
	index int
	indexInner int
}

func newComponentCollectionIter(ls []*componentData) *ComponentCollectionIter {
	return &ComponentCollectionIter{ls: ls}
}

func (p *ComponentCollectionIter) First() IComponent {
	return p.Next()
}

func (p *ComponentCollectionIter) Next() IComponent {
	if p.index == len(p.ls) {
		return nil
	}
	c:= p.ls[p.index].data[p.indexInner].(IComponent)
	if p.indexInner == len(p.ls[p.index].data)-1 {
		p.index += 1
		p.indexInner = 0
	}else{
		p.indexInner += 1
	}
	return c
}