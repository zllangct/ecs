package ecs

import "unsafe"

type ShapeIterator[T ShapeObject, TP ShapePointer[T]] interface {
	Begin() T
	Val() T
	Next() T
	End() bool
	Empty() bool
}

type ShapeIter[T ShapeObject, TP ShapePointer[T]] struct {
	c              *ShapeCollection[T, TP]
	len            int
	offset         int
	cur            *T
	curChunk       *Chunk
	curChunkOffset int64
}

func EmptyShapeIter[T ShapeObject, TP ShapePointer[T]]() ShapeIterator[T, TP] {
	return &ShapeIter[T, TP]{}
}

func NewShapeIterator[T ShapeObject, TP ShapePointer[T]](collection *ShapeCollection[T, TP]) ShapeIterator[T, TP] {
	iter := &ShapeIter[T, TP]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
	}
	if iter.len != 0 {
		iter.cur = new(T)
		TP(iter.cur).parse(unsafe.Pointer(&collection.data.data[0]), collection.eleSize)
		iter.curChunk = collection.data
	}
	return iter
}

func (i *ShapeIter[T, TP]) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *ShapeIter[T, TP]) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *ShapeIter[T, TP]) Begin() T {
	if i.len != 0 {
		i.offset = 0
		i.curChunkOffset = 0
		i.curChunk = i.c.data
		TP(i.cur).parse(unsafe.Pointer(&i.curChunk.data[0]), i.c.eleSize)
	}
	return *(i.cur)
}

func (i *ShapeIter[T, TP]) Val() T {
	return *(i.cur)
}

func (i *ShapeIter[T, TP]) Next() T {
	i.offset++
	i.curChunkOffset++
	if !i.End() {
		if i.curChunkOffset >= i.curChunk.len {
			if i.curChunk.next != nil {
				i.curChunk = i.curChunk.next
				i.curChunkOffset = 0
			} else {
				var empty T
				return empty
			}
		}
		TP(i.cur).parse(i.curChunk.GetByIndex(i.curChunkOffset), i.c.eleSize)
	}
	return *(i.cur)
}
