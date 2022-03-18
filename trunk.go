package ecs

type Trunk[T any] struct {
	id   string
	next *Trunk[T]
}
