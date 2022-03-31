package ecs

type Iterator[T ComponentObject, TP ComponentPointer[T]] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

type Iter[T ComponentObject, TP ComponentPointer[T]] struct {
	c        *Collection[T]
	len      int
	offset   int
	cur      *T
	curTemp  T
	readOnly bool
}

func EmptyIter[T ComponentObject, TP ComponentPointer[T]]() Iterator[T, TP] {
	return &Iter[T, TP]{}
}

func NewIterator[T ComponentObject, TP ComponentPointer[T]](collection *Collection[T], readOnly ...bool) Iterator[T, TP] {
	iter := &Iter[T, TP]{
		c:      collection,
		len:    collection.Len(),
		offset: 0,
	}
	if len(readOnly) > 0 {
		iter.readOnly = readOnly[0]
	}
	if iter.len != 0 {
		if iter.readOnly {
			iter.curTemp = collection.data[0]
			iter.cur = &iter.curTemp
		} else {
			iter.cur = &(collection.data[0])
		}
	}

	return iter
}

func (i *Iter[T, TP]) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T, TP]) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		if i.readOnly {
			i.curTemp = i.c.data[0]
			i.cur = &i.curTemp
		} else {
			i.cur = &(i.c.data[0])
		}
	}
	return i.cur
}

func (i *Iter[T, TP]) Val() *T {
	return i.cur
}

func (i *Iter[T, TP]) Next() *T {
	i.offset++
	if !i.End() {
		if i.readOnly {
			i.curTemp = i.c.data[i.offset]
			i.cur = &i.curTemp
		} else {
			i.cur = &(i.c.data[i.offset])
		}
		i.cur = &(i.c.data[i.offset])
	}
	return i.cur
}
