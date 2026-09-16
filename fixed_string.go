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
	max := int(unsafe.Sizeof(f.data))
	n := len(s)
	if n > max {
		n = max // 超长截断，防止 String() 越界
	}
	f.len = n
	if n != 0 {
		copy((*(*[__FixedMax]byte)(unsafe.Pointer(&(f.data))))[:max], s[:n])
	}
}
