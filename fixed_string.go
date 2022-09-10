package ecs

import (
	"fmt"
	"unsafe"
)

const (
	__FixedMax = 128
)

type FixedString[T any] struct {
	data T
	len  int
}

func (f *FixedString[T]) Clear() {
	f.len = 0
}

func (f *FixedString[T]) Empty() bool {
	return f.len == 0
}

func (f *FixedString[T]) Len() int {
	return f.len
}

func (f *FixedString[T]) String() string {
	return string((*(*[__FixedMax]byte)(unsafe.Pointer(&(f.data))))[:f.len])
}

func (f *FixedString[T]) Set(s string) {
	if len(s) > int(unsafe.Sizeof(f.data)) {
		panic(fmt.Sprintf("fixed string max size: %d, received size: %d", unsafe.Sizeof(f.data), len(s)))
	}
	f.len = len(s)
	if f.len != 0 {
		copy((*(*[__FixedMax]byte)(unsafe.Pointer(&(f.data))))[:unsafe.Sizeof(f.data)], s)
	}
}
