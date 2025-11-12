//go:build goexperiment.arenas

package ecs

import "arena"

type allocator[T any] struct {
	a *arena.Arena
}

func (a *allocator[T]) isArena() bool {
	return true
}

func (a *allocator[T]) alloc(size int, cap int) []T {
	if a.a == nil {
		a.a = arena.NewArena()
	}
	return arena.MakeSlice[T](a.a, size, cap)
}

func (a *allocator[T]) free() {
	if a.a != nil {
		a.a.Free()
	}
}
