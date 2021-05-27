package ecs

type IContainer[T any] interface {
	End() T
	Next() T
}