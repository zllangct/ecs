package ecs

import (
	"reflect"
	"unsafe"
)

type UnorderedContainer[T any] struct {
	buf  []byte
	len  int
	unit uintptr
	head uintptr
}

func NewUnorderedContainer[T any]() *UnorderedContainer[T] {
	var ins T
	size := reflect.TypeOf(ins).Size()
	c := &UnorderedContainer{
		buf:  make([]byte, 0, size),
		len:  0,
		unit: size,
	}
	return c
}

func (p *UnorderedContainer[T]) Add(item T) (int, *T) {
	pointer := unsafe.Pointer(&item)
	idx, newPointer := p.AddOriginal(pointer)
	return idx, (*T)(newPointer)
}

func (p *UnorderedContainer[T]) AddOriginal(pointer unsafe.Pointer) (int, unsafe.Pointer) {
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

func (p *UnorderedContainer[T]) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	offsetDelete := p.head + uintptr(idx)*p.unit
	offsetEnd := p.head + uintptr(p.len)*p.unit
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.unit)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *UnorderedContainer[T]) Get(idx int) *T {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return (*T)(unsafe.Pointer(p.head + uintptr(idx)*p.unit))
}

func (p *UnorderedContainer[T]) Len() int {
	return p.len
}

func (p *UnorderedContainer[T]) Iterator() IIterator[T] {
	return NewIterator[T](p)
}
