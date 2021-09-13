package ecs

type IContainer[T any] interface {
	Add(*T) (int, *T)
	Remove(int)
	Get(int) *T
	Len() int
}

type IIterator[T any] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
}
