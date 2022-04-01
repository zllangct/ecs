package ecs

import "unsafe"

type Iterator[T ComponentObject, TP ComponentPointer[T]] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
	Empty() bool
}

type Iter[T ComponentObject, TP ComponentPointer[T]] struct {
	head     unsafe.Pointer
	data     []T
	len      int
	offset   int
	pend     uintptr
	cur      *T
	curTemp  T
	eleSize  uintptr
	readOnly bool
}

func EmptyIter[T ComponentObject, TP ComponentPointer[T]]() Iterator[T, TP] {
	return &Iter[T, TP]{}
}

func NewIterator[T ComponentObject, TP ComponentPointer[T]](collection *Collection[T], readOnly ...bool) Iterator[T, TP] {
	iter := &Iter[T, TP]{
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

func (i *Iter[T, TP]) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T, TP]) End() bool {
	if i.offset >= i.len || i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T, TP]) Begin() *T {
	if i.len != 0 {
		i.offset = 0
		if i.readOnly {
			i.curTemp = i.data[0]
			i.cur = &i.curTemp
		} else {
			i.cur = &(i.data[0])
		}
	}
	return i.cur
}

func (i *Iter[T, TP]) Val() *T {
	return i.cur
}

func (i *Iter[T, TP]) Next() *T {
	i.offset++
	i.pend += i.eleSize
	if !i.End() {
		if i.readOnly {
			i.curTemp = *(*T)(unsafe.Pointer(uintptr(i.head) + i.pend))
			i.cur = &i.curTemp
		} else {
			i.cur = (*T)(unsafe.Pointer(uintptr(i.head) + i.pend))
		}
	} else {
		i.cur = nil
	}
	return i.cur
}
