package ecs

import "unsafe"

type Iter2[T ComponentObject, TP ComponentPointer[T]] struct {
	c              *Collection2[T, TP]
	len            int
	offset         int
	cur            *T
	curChunk       *Chunk[T, TP]
	curChunkOffset int64
}

func EmptyIter2[T ComponentObject, TP ComponentPointer[T]]() Iterator[T, TP] {
	return &Iter2[T, TP]{}
}

func NewIterator2[T ComponentObject, TP ComponentPointer[T]](collection *Collection2[T, TP]) Iterator[T, TP] {
	iter := &Iter2[T, TP]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
	}
	if iter.len != 0 {
		iter.cur = (*T)(unsafe.Pointer(&collection.data.data[0]))
		iter.curChunk = collection.data
	}
	return iter
}

func (i *Iter2[T, TP]) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *Iter2[T, TP]) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *Iter2[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		i.curChunkOffset = 0
		i.curChunk = i.c.data
		i.cur = (*T)(unsafe.Pointer(&i.curChunk.data[0]))
	}
	return i.cur
}

func (i *Iter2[T, TP]) Val() *T {
	return i.cur
}

func (i *Iter2[T, TP]) Next() *T {
	i.offset++
	i.curChunkOffset++
	if !i.End() {
		if i.curChunkOffset >= i.curChunk.len {
			if i.curChunk.next != nil {
				i.curChunk = i.curChunk.next
				i.curChunkOffset = 0
			} else {
				return nil
			}
		}
		i.cur = &i.curChunk.data[i.curChunkOffset]
	}
	return i.cur
}
