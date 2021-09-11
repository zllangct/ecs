package ecs

type IContainer[T any] interface {
	Add(*T) (int, *T)
	Remove(int)
	Get(int) *T
	Len() int
}

type IIterator[T any] interface {
	Val() *T
	Next()
	End() bool
}
