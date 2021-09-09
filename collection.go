package ecs

import (
	"reflect"
	"unsafe"
)

type Collection[T any] struct {
	buf  []byte
	len  int
	unit uintptr
	head uintptr
}

func NewCollection[T any]() *Collection[T] {
	var ins T
	size := reflect.TypeOf(ins).Size()
	c := &Collection[T]{
		buf:  make([]byte, 0, size),
		len:  0,
		unit: size,
	}
	return c
}

func (p *Collection[T]) Add(item T) (int, *T) {
	pointer := unsafe.Pointer(&item)
	idx, newPointer := p.AddOriginal(pointer)
	return idx, (*T)(newPointer)
}

func (p *Collection[T]) AddOriginal(pointer unsafe.Pointer) (int, unsafe.Pointer) {
	data := reflect.SliceHeader{
		Data: uintptr(pointer),
		Len:  int(p.unit),
		Cap:  int(p.unit),
	}
	p.buf = append(p.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	p.head = (*reflect.SliceHeader)(unsafe.Pointer(&p.buf)).Data
	p.len += 1
	return p.len - 1, unsafe.Pointer(p.head + uintptr(p.len-1)*p.unit)
}

func (p *Collection[T]) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return nil
	}
	offsetDelete := p.head + uintptr(idx)*p.unit
	offsetEnd := p.head + uintptr(p.len)*p.unit
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.unit)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *Collection[T]) Get(idx int) *T {
	return (*T)(p.GetOriginal(idx))
}

func (p *Collection[T]) GetOriginal(idx int) unsafe.Pointer {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return unsafe.Pointer(p.head + uintptr(idx)*p.unit)
}

func (p *Collection[T]) Len() int {
	return p.len
}

//func (p *Collection[T]) Iterator() IIterator[T] {
//	return NewIterator[T](p)
//}
