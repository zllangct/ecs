package ecs

type IShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

type shapeCache[T ShapeObject] struct {
	cache  []T
	cached []bool
}

type ShapeIter[T ShapeObject, TP ShapeObjectPointer[T]] struct {
	c      ICollection
	cache  *shapeCache[T]
	len    int
	offset int
	req    []IRequirement
	cur    *T
	end    bool
}

func EmptyShapeIter[T ShapeObject, TP ShapeObjectPointer[T]]() IShapeIterator[T, TP] {
	return &ShapeIter[T, TP]{}
}

func NewShapeIterator[T ShapeObject, TP ShapeObjectPointer[T]](collection ICollection, req []IRequirement, cache any) IShapeIterator[T, TP] {
	iter := &ShapeIter[T, TP]{
		c:      collection,
		len:    collection.Len(),
		cache:  nil,
		req:    req,
		offset: 0,
	}

	if sc, ok := cache.(*shapeCache[T]); ok {
		iter.cache = sc
	}

	if iter.len != 0 {
		if iter.cache != nil && iter.cache.cached[0] {
			iter.cur = &iter.cache.cache[0]
		} else {
			iter.cur = new(T)
			com := collection.getByIndex(0)
			TP(iter.cur).parse(com.(IComponent).Owner(), iter.req)
			iter.cache.cache[0] = *iter.cur
			iter.cache.cached[0] = true
		}
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
		if i.cache != nil && i.cache.cached[i.offset] {
			i.cur = &i.cache.cache[i.offset]
		} else {
			if i.cur == nil {
				i.cur = new(T)
			}
			com := i.c.getByIndex(int64(i.offset))
			TP(i.cur).parse(com.(IComponent).Owner(), i.req)
			i.cache.cache[0] = *i.cur
			i.cache.cached[i.offset] = true
		}
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Val() *T {
	return i.cur
}

func (i *ShapeIter[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		if i.cache != nil && i.cache.cached[i.offset] {
			i.cur = &i.cache.cache[i.offset]
		} else {
			com := i.c.getByIndex(int64(i.offset))
			TP(i.cur).parse(com.(IComponent).Owner(), i.req)
			i.cache.cache[i.offset] = *i.cur
			i.cache.cached[i.offset] = true
		}
	} else {
		i.cur = nil
	}
	return i.cur
}

func (i *ShapeIter[T, TP]) Len() int {
	return i.len
}
