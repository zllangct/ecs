package ecs

import (
	"unsafe"
)

const (
	InitMaxSize = 1024 * 16
	//InitMaxSize        = 0
	SeqMax uint32 = 0xFFFFFFFF
)

type UnorderedCollection[T any] struct {
	eleSize  uintptr
	len      int64
	initSize int64
	data     []T
}

func NewUnorderedCollection[T any](initSize ...int) *UnorderedCollection[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	c := &UnorderedCollection[T]{
		data: make([]T, 0, size),
	}
	if len(initSize) > 0 {
		c.initSize = int64(initSize[0])
		c.eleSize = uintptr(initSize[0]) / eleSize
	}
	return c
}

func (c *UnorderedCollection[T]) Get(idx int64) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(&c.data[0]), uintptr(idx)*c.eleSize))
}

func (c *UnorderedCollection[T]) Add(element *T) (*T, int64) {
	if int64(len(c.data)) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	idx := c.len
	c.len++
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
	removed := c.data[lastIdx]
	return &removed, lastIdx, idx
}

func (c *UnorderedCollection[T]) Len() int {
	return int(c.len)
}

func (c *UnorderedCollection[T]) Range(f func(element *T) bool) {
	for i := int64(0); i < c.len; i++ {
		if !f(&c.data[i]) {
			break
		}
	}
}

func (c *UnorderedCollection[T]) Clear() {
	c.data = make([]T, 0, c.initSize)
	c.len = 0
}

func (c *UnorderedCollection[T]) shrink() {
	var threshold int64
	if len(c.data) < InitMaxSize {
		return
	} else {
		threshold = int64(float64(c.len) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		//c.data = c.data[:threshold]
		newData := make([]T, threshold)
		copy(newData, c.data)
		c.data = newData
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
	idx := int64(offset / c.eleSize)
	if idx < 0 || idx > c.len-1 {
		return -1
	}
	return idx
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
