package ecs

import "unsafe"

type IItem interface {
	Ref() *Ref
	ID() int64
}

type Ref struct {
	idx int64
}

type UnorderedCollectionWithIDRef[T IItem] struct {
	UnorderedCollection[T]
	id2Ref map[int64]*Ref
	seq    int64
}

func NewUnorderedCollectionWithIDRef[T IItem](initSize ...int) *UnorderedCollectionWithIDRef[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	if len(initSize) > 0 {
		eleSize = uintptr(initSize[0]) / eleSize
	}
	c := &UnorderedCollectionWithIDRef[T]{
		id2Ref: map[int64]*Ref{},
		UnorderedCollection: UnorderedCollection[T]{
			data:    make([]T, 0, size),
			eleSize: eleSize,
		},
	}
	return c
}

func (c *UnorderedCollectionWithIDRef[T]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.id2Ref[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *UnorderedCollectionWithIDRef[T]) Add(element *T, elementID ...int64) (*T, int64) {
	_, idx := c.UnorderedCollection.Add(element)
	var id int64
	if len(elementID) > 0 {
		id = elementID[0]
	} else {
		id = c.getID()
	}
	ref := c.data[idx].Ref()
	c.id2Ref[id] = ref

	return &c.data[idx], id
}

func (c *UnorderedCollectionWithIDRef[T]) remove(id int64) *T {
	ref, ok := c.id2Ref[id]
	if !ok {
		return nil
	}

	removed, oldIndex, newIndex := c.UnorderedCollection.Remove(ref.idx)
	c.data[oldIndex].Ref().idx = newIndex
	c.id2Ref[newIndex] = c.data[oldIndex].Ref()
	delete(c.id2Ref, id)

	return removed
}

func (c *UnorderedCollectionWithIDRef[T]) Remove(id int64) {
	c.remove(id)
}

func (c *UnorderedCollectionWithIDRef[T]) RemoveAndReturn(id int64) *T {
	cpy := *c.remove(id)
	return &cpy
}

func (c *UnorderedCollectionWithIDRef[T]) getByID(id int64) *T {
	ref, ok := c.id2Ref[id]
	if !ok {
		return nil
	}
	return c.UnorderedCollection.Get(ref.idx)
}

func (c *UnorderedCollectionWithIDRef[T]) GetByID(id int64) any {
	return c.getByID(id)
}

func (c *UnorderedCollectionWithIDRef[T]) GetByIndex(idx int64) any {
	return c.UnorderedCollection.Get(idx)
}

func NewUnorderedCollectionWithIDRefIterator[T IItem](collection *UnorderedCollectionWithIDRef[T], readOnly ...bool) Iterator[T] {
	iter := &Iter[T]{
		data:    collection.data,
		len:     collection.Len(),
		eleSize: collection.eleSize,
		offset:  0,
	}
	if len(readOnly) > 0 {
		iter.readOnly = readOnly[0]
	}
	if iter.len != 0 {
		iter.head = unsafe.Pointer(&collection.data[0])
		if iter.readOnly {
			iter.curTemp = collection.data[0]
			iter.cur = &iter.curTemp
		} else {
			iter.cur = &(collection.data[0])
		}
	}

	return iter
}
