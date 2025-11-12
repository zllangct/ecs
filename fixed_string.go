package ecs

import (
	"unsafe"
)

const (
	__FixedMax = 1024
)

//go:generate go run ./cmd/ecs_internal_gen/main.go FixedString -p "ecs"
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
	f.len = len(s)
	if f.len != 0 {
		copy((*(*[__FixedMax]byte)(unsafe.Pointer(&(f.data))))[:unsafe.Sizeof(f.data)], s)
	}
}
