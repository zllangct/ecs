package ecs

type IterByte[T ComponentObject, TP ComponentPointer[T]] struct {
	c      *CollectionByte[T, TP]
	len    int
	offset int
	cur    *T
}

func EmptyIterByte[T ComponentObject, TP ComponentPointer[T]]() Iterator[T, TP] {
	return &IterByte[T, TP]{}
}

func NewIteratorByte[T ComponentObject, TP ComponentPointer[T]](collection *CollectionByte[T, TP]) Iterator[T, TP] {
	iter := &IterByte[T, TP]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
	}
	if iter.len != 0 {
		iter.cur = iter.c.GetByIndex(0)
	}
	return iter
}

func (i *IterByte[T, TP]) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *IterByte[T, TP]) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *IterByte[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		i.cur = i.c.GetByIndex(0)
	}
	return i.cur
}

func (i *IterByte[T, TP]) Val() *T {
	return i.cur
}

func (i *IterByte[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		i.cur = i.c.GetByIndex(int64(i.offset))
	}
	return i.cur
}
