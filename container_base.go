package ecs

import (
	"reflect"
	"unsafe"
)

type Container struct {
	buf        []byte
	len        int
	memberSize uintptr
	begin      uintptr
}

func NewContainer(size uintptr) *Container {
	return &Container{
		buf:        make([]byte, 0, size),
		len:        -1,
		memberSize: size,
	}
}

func (p *Container) add(pointer unsafe.Pointer) (int, unsafe.Pointer) {
	data := reflect.SliceHeader{
		Data: uintptr(pointer),
		Len:  int(p.memberSize),
		Cap:  int(p.memberSize),
	}
	p.buf = append(p.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	p.begin = (*reflect.SliceHeader)(unsafe.Pointer(&p.buf)).Data
	p.len += 1
	return p.len, unsafe.Pointer(p.begin + uintptr(p.len)*p.memberSize)
}

func (p *Container) remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	offsetDelete := p.begin + uintptr(idx)*p.memberSize
	offsetEnd := p.begin + uintptr(p.len)*p.memberSize
	copy(p.buf[offsetDelete:offsetDelete+uintptr(p.memberSize)], p.buf[offsetEnd:])
	p.buf = p.buf[:offsetEnd]
	p.len -= 1
}

func (p *Container) get(idx int) unsafe.Pointer {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return unsafe.Pointer(p.begin + uintptr(idx)*p.memberSize)
}

type Iterator struct {
	memberSize uintptr
	size       int
	index      int
	begin      uintptr
}

func NewIterator(container *Container) Iterator {
	return Iterator{
		memberSize: container.memberSize,
		size:       container.len,
		index:      -1,
		begin:      container.begin,
	}
}

func (p *Iterator) End() unsafe.Pointer {
	return nil
}

func (p *Iterator) Next() unsafe.Pointer {
	if p.index == p.size {
		return nil
	}
	p.index++
	return unsafe.Pointer(p.begin + uintptr(p.index)*p.memberSize)
}
