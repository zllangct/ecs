package ecs

import (
	"reflect"
	"unsafe"
)

type UnorderedContainerByte[T any] struct {
	buf  []byte
	len  int
	unit uintptr
	head uintptr
}

func NewUnorderedContainerByte[T any]() *UnorderedContainerByte[T] {
	var ins T
	size := reflect.TypeOf(ins).Size()
	c := &UnorderedContainerByte{
		buf:  make([]byte, 0, size),
		len:  0,
		unit: size,
	}
	return c
}

func (p *UnorderedContainerByte[T]) Add(item *T) (int, *T) {
	pointer	:= unsafe.Pointer(item)
	data := reflect.SliceHeader{
		Data: uintptr(pointer),
		Len:  int(p.unit),
		Cap:  int(p.unit),
	}
	p.buf = append(p.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	p.head = (*reflect.SliceHeader)(unsafe.Pointer(&p.buf)).Data
	p.len += 1
	return p.len - 1, (*T)(unsafe.Pointer(p.head + uintptr(p.len-1)*p.unit))
}

func (p *UnorderedContainerByte[T]) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	offsetDelete := p.head + uintptr(idx)*p.unit
	offsetEnd := p.head + uintptr(p.len)*p.unit
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.unit)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *UnorderedContainerByte[T]) Get(idx int) *T {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return (*T)(unsafe.Pointer(p.head + uintptr(idx)*p.unit))
}
