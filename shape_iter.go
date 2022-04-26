package ecs

type IShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

type ShapeIter[T ShapeObject, TP ShapeObjectPointer[T]] struct {
	shapes []T
	len    int
	offset int
	req    []IRequirement
	cur    *T
	end    bool
}

func EmptyShapeIter[T ShapeObject, TP ShapeObjectPointer[T]]() IShapeIterator[T, TP] {
	return &ShapeIter[T, TP]{}
}

func NewShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]](shapes []T) IShapeIterator[T, TP] {
	iter := &ShapeIter[T, TP]{
		len:    len(shapes),
		shapes: shapes,
		offset: 0,
	}

	if iter.len != 0 {
		iter.cur = &iter.shapes[0]
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

func (i *ShapeIter[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		i.cur = &i.shapes[i.offset]
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Val() *T {
	return i.cur
}

func (i *ShapeIter[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		i.cur = &i.shapes[i.offset]
	} else {
		i.cur = nil
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Len() int {
	return i.len
}
