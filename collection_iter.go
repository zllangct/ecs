package ecs

type IIterator[T any] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
}

type Iterator[T any] struct {
	c     *Collection[T]
	size  int
	index int
	cur   *T
}

func NewIterator[T any](collection *Collection[T]) IIterator[T] {
	iter := &Iterator[T]{
		c:     collection,
		size:  collection.Len(),
		index: 0,
		cur:   &(collection.data[0]),
	}
	return iter
}

func (i *Iterator[T]) End() bool {
	if i.index > i.size-1 || i.size == 0 {
		return true
	}
	return false
}

func (i *Iterator[T]) Begin() *T {
	return &(i.c.data[i.index])
}

func (i *Iterator[T]) Val() *T {
	return i.cur
}

func (i *Iterator[T]) Next() *T {
	i.index++
	if !i.End() {
		i.cur = &(i.c.data[i.index])
	}
	return i.cur
}
