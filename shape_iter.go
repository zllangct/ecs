package ecs

type IShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

type ShapeIter[T ShapeObject, TP ShapeObjectPointer[T]] struct {
	c      ICollection
	len    int
	offset int
	req    []IRequirement
	cur    *T
	end    bool
}

func EmptyShapeIter[T ShapeObject, TP ShapeObjectPointer[T]]() IShapeIterator[T, TP] {
	return &ShapeIter[T, TP]{}
}

func NewShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]](collection ICollection, req []IRequirement) IShapeIterator[T, TP] {
	iter := &ShapeIter[T, TP]{
		c:      collection,
		len:    collection.Len(),
		cur:    new(T),
		req:    req,
		offset: 0,
	}

	if iter.len != 0 {
		com := collection.getByIndex(0)
		TP(iter.cur).parse(com.Owner(), iter.req)
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
		com := i.c.getByIndex(int64(i.offset))
		TP(i.cur).parse(com.Owner(), i.req)
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Val() *T {
	return i.cur
}

func (i *ShapeIter[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		com := i.c.getByIndex(int64(i.offset))
		TP(i.cur).parse(com.Owner(), i.req)
	}
	return i.cur
}
