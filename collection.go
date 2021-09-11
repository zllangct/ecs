package ecs

import (
	"reflect"
	"unsafe"
)

type Collection struct {
	buf  []byte
	len  int
	unit int
	ids map[int64]int64
	seq int64
}

func NewCollection(elementSize int) *Collection {
	c := &Collection{
		buf:  make([]byte, 0, elementSize),
		len:  0,
		unit: elementSize,
		ids: map[int64]int64{},
	}
	return c
}

func (c *Collection) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist :=c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *Collection) add(element unsafe.Pointer) (int, unsafe.Pointer) {
	data := reflect.SliceHeader{
		Data: uintptr(element),
		Len:  c.unit,
		Cap:  c.unit,
	}
	c.buf = append(c.buf, *(*[]byte)(unsafe.Pointer(&data))...)
	c.len += 1
	return c.len - 1, unsafe.Pointer(&c.buf)
}

func (c *Collection) Add(element unsafe.Pointer) (int64, unsafe.Pointer)  {
	idx, ptr := c.add(element)
	id := c.getID()
	c.ids[id] = int64(idx)
	c.ids[int64(-idx)] = -id
	return id, ptr
}

func (c *Collection) remove(idx int) {
	if idx < 0 || idx >= c.len {
		return
	}
	offsetDelete := idx * c.unit
	offsetEnd := c.len * c.unit
	copy(c.buf[offsetDelete:offsetDelete + c.unit], c.buf[offsetEnd:])
	c.buf = c.buf[:offsetEnd]
	c.len -= 1
}

func (c *Collection) Remove(id int64) {
	if id < 0 {
		return
	}
	idx, ok := c.ids[id]
	if !ok {
		return
	}
	c.ids[c.ids[int64(c.len)]] = idx
	delete(c.ids, c.ids[-idx])
	c.ids[-idx] = c.ids[int64(c.len)]
	delete(c.ids, int64(c.len))

	c.remove(int(idx))
}

func (c *Collection) get(idx int) unsafe.Pointer {
	if idx < 0 || idx >= c.len {
		return nil
	}
	return unsafe.Add(unsafe.Pointer(&c.buf), idx * c.unit)
}

func (c *Collection) Get(id int64) unsafe.Pointer {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return c.get(int(idx))
}

func (c *Collection) Len() int {
	return c.len
}
