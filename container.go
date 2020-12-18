package ecs

import (
	"reflect"
	"unsafe"
)

type Container struct {
	buf  []byte
	len  int
	unit uintptr
	head uintptr
}

func NewContainer(size uintptr) *Container {
	return &Container{
		buf:  make([]byte, 0, size),
		len:  -1,
		unit: size,
	}
}

func (p *Container) Add(pointer unsafe.Pointer) (int, unsafe.Pointer) {
	data := reflect.SliceHeader{
		Data: uintptr(pointer),
		Len:  int(p.unit),
		Cap:  int(p.unit),
	}
	p.buf = append(p.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	p.head = (*reflect.SliceHeader)(unsafe.Pointer(&p.buf)).Data
	p.len += 1
	return p.len, unsafe.Pointer(p.head + uintptr(p.len)*p.unit)
}

func (p *Container) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	offsetDelete := p.head + uintptr(idx)*p.unit
	offsetEnd := p.head + uintptr(p.len)*p.unit
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.unit)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *Container) Get(idx int) unsafe.Pointer {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return unsafe.Pointer(p.head + uintptr(idx)*p.unit)
}

func (p Container) GetIterator() *iterator {
	return &iterator{
		memberSize: p.unit,
		size:       p.len,
		index:      -1,
		head:       p.head,
	}
}
