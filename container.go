package ecs

type IContainer[T any] interface {
	Add(*T) (int *T)
	Remove(int)
	Get(int) *T
	Len() int
}

type IIterator[T any] interface {
	Next() *T
	End() *T
}