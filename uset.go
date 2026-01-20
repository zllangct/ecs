package ecs

import (
	"iter"
	"unsafe"
)

var EnableShrink = true

const (
	InitMaxSize = 64
)

type USetSerializable[T any] struct {
	eleSize  uint64
	len      int64
	initSize int64
	data     []T
}

type USet[T any] struct {
	USetSerializable[T]
	a allocator[T]
}

func NewUSet[T any](initSize ...int) *USet[T] {
	typ := TypeOf[T]()
	size := int64(0)
	if len(initSize) > 0 {
		size = int64(initSize[0])
	}
	c := &USet[T]{}
	c.a = allocator[T]{}
	c.data = c.a.alloc(0, int(size))

	if len(initSize) > 0 {
		c.initSize = size
		c.eleSize = uint64(typ.Size())
	}
	return c
}

func (u *USet[T]) Free() {
	u.a.free()
	u.len = 0
	u.data = nil
}

func (u *USet[T]) Get(idx int64) *T {
	return (*T)(unsafe.Add(unsafe.Pointer(&u.data[0]), uintptr(idx)*uintptr(u.eleSize)))
}

func (u *USet[T]) Add(element *T) (*T, int64) {
	if int64(len(u.data)) > u.len {
		u.data[u.len] = *element
	} else {
		u.data = append(u.data, *element)
	}
	idx := u.len
	u.len++
	return &u.data[idx], idx
}

func (u *USet[T]) Remove(idx int64) (*T, int64, int64) {
	if idx < 0 {
		return nil, 0, 0
	}
	lastIdx := u.len - 1

	u.data[idx], u.data[lastIdx] = u.data[lastIdx], u.data[idx]
	u.shrink()
	u.len--
	removed := u.data[lastIdx]
	return &removed, lastIdx, idx
}

func (u *USet[T]) Len() int {
	return int(u.len)
}

func (u *USet[T]) Reset() {
	u.a.free()
	u.data = u.a.alloc(0, int(u.initSize))
	u.len = 0
}

func (u *USet[T]) Swap(i, j int64) {
	u.data[i], u.data[j] = u.data[j], u.data[i]
}

func (u *USet[T]) shrink() {
	var threshold int64
	if !EnableShrink || len(u.data) < 1024 || len(u.data) < InitMaxSize {
		return
	}
	if u.a.isArena() {
		threshold = int64(float64(u.len) * 5)
		if int64(len(u.data)) > threshold {
			targetSize := u.len * 2
			temp := make([]T, u.len)
			copy(temp, u.data[:u.len])
			u.a.free()
			u.data = u.a.alloc(int(targetSize), int(targetSize))
			copy(u.data, temp)
		}
	} else {
		threshold = int64(float64(u.len) * 3)
		if int64(len(u.data)) > threshold {
			targetSize := int64(float64(u.len) * 1.5)
			newData := make([]T, targetSize)
			copy(newData, u.data[:u.len])
			u.data = newData
		}
	}
}

func (u *USet[T]) getIndexByElePointer(element *T) int64 {
	if u.len == 0 {
		return -1
	}
	offset := uintptr(unsafe.Pointer(element)) - uintptr(unsafe.Pointer(&u.data[0]))
	if offset%uintptr(u.eleSize) != 0 {
		return -1
	}
	idx := int64(offset / uintptr(u.eleSize))
	if idx < 0 || idx > u.len-1 {
		return -1
	}
	return idx
}

func (u *USet[T]) Iter() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		for i := 0; i < int(u.len); i++ {
			yield(i, &u.data[i])
		}
	}
}
