package ecs

type Iterator[T any] struct {
	c 	  	*Collection
	size  	int
	index	int
}

func NewIterator[T any](collection *Collection) IIterator[T] {
	iter:= &Iterator[T]{
		c :    collection,
		size:  collection.Len(),
		index: 0,
	}
	return iter
}

func (i *Iterator[T]) End() bool  {
	if i.index >= i.size-1 || i.size == 0 {
		return true
	}
	return false
}

func (i *Iterator[T]) Val() *T {
	item := i.c.get(i.index)
	return  (*T)(item)
}

func (i *Iterator[T]) Next() {
	i.index ++
}
