package ecs

import "unsafe"

type ICollectionWithID interface {
	ICollection
	GetByIndex(idx int64) any
	GetByID(id int64) any
}

type UnorderedCollectionWithID[T any] struct {
	UnorderedCollection[T]
	ids    map[int64]int64
	idx2id map[int64]int64
	seq    int64
}

func NewUnorderedCollectionWithID[T any](initSize ...int) *UnorderedCollectionWithID[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	if len(initSize) > 0 {
		eleSize = uintptr(initSize[0]) / eleSize
	}
	c := &UnorderedCollectionWithID[T]{
		ids:    map[int64]int64{},
		idx2id: map[int64]int64{},
		UnorderedCollection: UnorderedCollection[T]{
			data:    make([]T, 0, size),
			eleSize: eleSize,
		},
	}
	return c
}

func (c *UnorderedCollectionWithID[T]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *UnorderedCollectionWithID[T]) Add(element *T, elementID ...int64) (*T, int64) {
	_, idx := c.UnorderedCollection.Add(element)
	var id int64
	if len(elementID) > 0 {
		id = elementID[0]
	} else {
		id = c.getID()
	}
	c.ids[id] = idx
	c.idx2id[idx] = id

	return &c.data[idx], id
}

func (c *UnorderedCollectionWithID[T]) remove(id int64) *T {
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}

	removed, oldIndex, newIndex := c.UnorderedCollection.Remove(idx)

	lastId := c.idx2id[oldIndex]
	c.ids[lastId] = newIndex
	c.idx2id[newIndex] = lastId
	delete(c.idx2id, oldIndex)
	delete(c.ids, id)

	return removed
}

func (c *UnorderedCollectionWithID[T]) Remove(id int64) {
	c.remove(id)
}

func (c *UnorderedCollectionWithID[T]) RemoveAndReturn(id int64) *T {
	cpy := *c.remove(id)
	return &cpy
}

func (c *UnorderedCollectionWithID[T]) getByID(id int64) *T {
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return c.UnorderedCollection.Get(idx)
}

func (c *UnorderedCollectionWithID[T]) GetByID(id int64) any {
	return c.getByID(id)
}

func (c *UnorderedCollectionWithID[T]) GetByIndex(idx int64) any {
	return c.UnorderedCollection.Get(idx)
}

func NewUnorderedCollectionWithIDIterator[T ComponentObject](collection *UnorderedCollectionWithID[T], readOnly ...bool) Iterator[T] {
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
