package ecs

import (
	"reflect"
	"unsafe"
)

type Collection struct {
	buf  []byte
	len  int
	unit uintptr
}

func NewCollection(elementSize uintptr) *Collection {
	c := &Collection{
		buf:  make([]byte, 0, elementSize),
		len:  0,
		unit: elementSize,
	}
	return c
}

func (p *Collection) Head() uintptr {
	return uintptr(unsafe.Pointer(&p.buf))
}

func (p *Collection) Add(pointer unsafe.Pointer) (int, unsafe.Pointer) {
	data := reflect.SliceHeader{
		Data: uintptr(pointer),
		Len:  int(p.unit),
		Cap:  int(p.unit),
	}
	p.buf = append(p.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	p.len += 1
	return p.len - 1, unsafe.Pointer(&p.buf)
}

func (p *Collection) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	offsetDelete := p.Head() + uintptr(idx)*p.unit
	offsetEnd := p.Head() + uintptr(p.len)*p.unit
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.unit)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *Collection) Get(idx int) unsafe.Pointer {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return unsafe.Pointer(p.Head() + uintptr(idx)*p.unit)
}

func (p *Collection) Len() int {
	return p.len
}
