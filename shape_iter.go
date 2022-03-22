package ecs

type ShapeIterator[T ShapeObject] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

/*
type ShapeIter[T ShapeObject, TP ShapePointer[T]] struct {
	c      *ShapeCollection[T, TP]
	len    int
	offset int
	cur    *T

	eleSize []uintptr
}

func EmptyShapeIter[T ShapeObject, TP ShapePointer[T]]() ShapeIter[T, TP] {
	return &ShapeIter[T, TP]{}
}

func NewShapeIterator[T ShapeObject, TP ShapePointer[T]](collection *ShapeCollection[T, TP]) ShapeIter[T, TP] {
	iter := &ShapeIter[T, TP]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
	}
	if iter.len != 0 {
		iter.cur = &(collection.data[0])
	}
	iter.eleSize = collection.getEleSize()
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

func (i *ShapeIter[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		i.cur = &(i.c.data[0])
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Val() *T {
	return i.cur
}

func (i *ShapeIter[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		i.cur = &(i.c.data[i.offset])
	}
	return i.cur
}
*/
