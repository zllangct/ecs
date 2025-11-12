//go:build !goexperiment.arenas

package ecs

type allocator[T any] struct {
}

func (a *allocator[T]) isArena() bool {
	return false
}

func (a *allocator[T]) alloc(size int, cap int) []T {
	return make([]T, size, cap)
}

func (a *allocator[T]) free() {}
