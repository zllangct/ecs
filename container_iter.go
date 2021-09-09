package ecs

type Iterator[T any] struct {
	c 	  	IContainer[T]
	size  	int
	index	int
}

func NewIterator[T any](container IContainer[T]) IIterator[T] {
	return &Iterator[T]{
		c : container,
		size: container.Len(),
		index:      0,
	}
}

func (p *Iterator[T]) Next() *T {
	if p.index >= p.size-1 || p.size == 0 {
		return nil
	}
	item := p.c.Get(p.index)
	p.index ++
	return item
}
