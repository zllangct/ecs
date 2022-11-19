package ecs

import "unsafe"

type Iterator[T any] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
}

type Iter[T any] struct {
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

func EmptyIter[T any]() Iterator[T] {
	return &Iter[T]{}
}

func (i *Iter[T]) End() bool {
	if i.offset >= i.len || i.len == 0 {
		return true
	}
	return false
}

func (i *Iter[T]) Begin() *T {
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

func (i *Iter[T]) Val() *T {
	return i.cur
}

func (i *Iter[T]) Next() *T {
	i.offset++
	i.pend += i.eleSize
	if !i.End() {
		if i.readOnly {
			//i.curTemp = i.data[i.offset]
			i.curTemp = *(*T)(unsafe.Add(i.head, i.pend))
			i.cur = &i.curTemp
		} else {
			//i.cur = &(i.data[i.offset])
			i.cur = (*T)(unsafe.Add(i.head, i.pend))
		}
	} else {
		i.cur = nil
	}
	return i.cur
}
