package ecs

import (
	"arena"
	"unsafe"
)

const (
	SeqMax        uint32 = 0xFFFFFFFF
	ShrinkMinSize int    = 1024
)

type UnorderedCollection[T any] struct {
	isArena  bool
	a        *arena.Arena
	eleSize  uintptr
	len      int
	initSize int
	data     []T
}

func NewUnorderedCollection[T any](initSize ...int) *UnorderedCollection[T] {
	c := &UnorderedCollection[T]{}
	size := 0
	for _, v := range initSize {
		size = v
	}
	c.init(true, size)
	return c
}

func NewUnorderedCollectionNoArena[T any](initSize ...int) *UnorderedCollection[T] {
	c := &UnorderedCollection[T]{}
	size := 0
	for _, v := range initSize {
		size = v
	}
	c.init(false, size)
	return c
}

func (c *UnorderedCollection[T]) init(isArena bool, initSize int) {
	c.isArena = isArena
	typ := TypeOf[T]()
	c.initSize = initSize
	c.eleSize = typ.Size()
	c.alloc(0, initSize)
	c.len = 0
}

func (c *UnorderedCollection[T]) alloc(size int, cap int) {
	if c.isArena {
		if c.a != nil {
			c.a.Free()
		}
		c.a = arena.NewArena()
		c.data = arena.MakeSlice[T](c.a, size, cap)
	} else {
		c.a = nil
		c.data = make([]T, size, cap)
	}
}

func (c *UnorderedCollection[T]) Free() {
	if c.a != nil {
		c.a.Free()
	}
	c.len = 0
	c.data = nil
}

func (c *UnorderedCollection[T]) Get(idx int64) *T {
	//return &c.data[idx]
	return (*T)(unsafe.Add(unsafe.Pointer(&c.data[0]), uintptr(idx)*c.eleSize))
}

func (c *UnorderedCollection[T]) Add(element *T) (*T, int64) {
	if len(c.data) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	idx := c.len
	c.len++
	return &c.data[idx], int64(idx)
}

func (c *UnorderedCollection[T]) Remove(idx int64) (*T, int64, int64) {
	if idx < 0 {
		return nil, 0, 0
	}
	lastIdx := c.len - 1

	c.data[idx], c.data[lastIdx] = c.data[lastIdx], c.data[idx]
	c.len--
	removed := c.data[lastIdx]
	c.shrink()
	return &removed, int64(lastIdx), idx
}

func (c *UnorderedCollection[T]) Len() int {
	return c.len
}

func (c *UnorderedCollection[T]) Range(f func(element *T) bool) {
	for i := 0; i < c.len; i++ {
		if !f(&c.data[i]) {
			break
		}
	}
}

func (c *UnorderedCollection[T]) Reset() {
	c.data = c.data[0:0]
	c.len = 0
	c.shrink()
}

func (c *UnorderedCollection[T]) shrink() {
	capSize := cap(c.data)
	if capSize < ShrinkMinSize || capSize <= c.initSize {
		return
	}
	threshold := int(float64(c.len) * 2)
	if cap(c.data) > threshold {
		resize := int(float64(c.len) * 1.25)
		var temp []T
		if c.isArena {
			temp = arena.Clone(c.data)
		} else {
			temp = c.data
		}
		c.alloc(c.len, resize)
		copy(c.data, temp[:c.len])
	}
}

func (c *UnorderedCollection[T]) getIndexByElePointer(element *T) int64 {
	if c.len == 0 {
		return -1
	}
	offset := uintptr(unsafe.Pointer(element)) - uintptr(unsafe.Pointer(&c.data[0]))
	if offset%c.eleSize != 0 {
		return -1
	}
	idx := int(offset / c.eleSize)
	if idx < 0 || idx > c.len-1 {
		return -1
	}
	return int64(idx)
}

func NewUnorderedCollectionIterator[T ComponentObject](collection *UnorderedCollection[T], readOnly ...bool) Iterator[T] {
	iter := &Iter[T]{
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
