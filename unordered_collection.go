package ecs

import (
	"reflect"
	"unsafe"
)

const (
	//InitMaxSize        = 1024 * 16
	InitMaxSize        = 0
	SeqMax      uint32 = 0xFFFFFFFF
)

type ICollection interface {
	Len() int
	Range(func(v any) bool)
	ChangeCount() int64
	ChangeReset()
	ElementType() reflect.Type
}

type UnorderedCollection[T any] struct {
	eleSize uintptr
	len     int64
	change  int64
	data    []T
}

func NewUnorderedCollection[T any](initSize ...int) *UnorderedCollection[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	if len(initSize) > 0 {
		eleSize = uintptr(initSize[0]) / eleSize
	}
	c := &UnorderedCollection[T]{
		data:    make([]T, 0, size),
		eleSize: eleSize,
	}
	return c
}

func (c *UnorderedCollection[T]) Get(idx int64) *T {
	return (*T)(unsafe.Pointer(uintptr(unsafe.Pointer(&c.data[0])) + uintptr(idx)*c.eleSize))
}

func (c *UnorderedCollection[T]) Add(element *T) (*T, int64) {
	if int64(len(c.data)) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	idx := c.len
	c.len++
	c.change++
	return &c.data[idx], idx
}

func (c *UnorderedCollection[T]) Remove(idx int64) (*T, int64, int64) {
	if idx < 0 {
		return nil, 0, 0
	}
	lastIdx := c.len - 1

	c.data[idx], c.data[lastIdx] = c.data[lastIdx], c.data[idx]
	c.shrink()
	c.len--
	c.change++
	return &c.data[lastIdx], lastIdx, idx
}

func (c *UnorderedCollection[T]) Len() int {
	return int(c.len)
}

func (c *UnorderedCollection[T]) RangeRaw(f func(element *T) bool) {
	for i := int64(0); i < c.len; i++ {
		if !f(&c.data[i]) {
			break
		}
	}
}

func (c *UnorderedCollection[T]) Range(f func(element any) bool) {
	for i := int64(0); i < c.len; i++ {
		if !f(&c.data[i]) {
			break
		}
	}
}

func (c *UnorderedCollection[T]) shrink() {
	var threshold int64
	if len(c.data) < InitMaxSize {
		return
	} else {
		threshold = int64(float64(c.len) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		c.data = c.data[:threshold]
	}
}

func (c *UnorderedCollection[T]) ChangeCount() int64 {
	return c.change
}

func (c *UnorderedCollection[T]) ChangeReset() {
	c.change = 0
}

func (c *UnorderedCollection[T]) ElementType() reflect.Type {
	return TypeOf[T]()
}

func (c *UnorderedCollection[T]) getData() []T {
	return c.data
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
