package ecs

type Iterator[T any] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
}

type Iter[T any] struct {
	c      *Collection[T]
	len    int
	offset int
	cur    *T
}

func NewIterator[T any](collection *Collection[T]) Iterator[T] {
	iter := &Iter[T]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
		cur:    &(collection.data[0]),
	}
	return iter
}

func (i *Iter[T]) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T]) Begin() *T {
	return i.cur
}

func (i *Iter[T]) Val() *T {
	return i.cur
}

func (i *Iter[T]) Next() *T {
	i.offset++
	if !i.End() {
		i.cur = &(i.c.data[i.offset])
	}
	return i.cur
}
