package ecs

import (
	"unsafe"
)

type ComponentSet interface {
	Len() int
	Reset()
	Add(entity Entity, comp Component) Component
	Get(entity Entity) Component
	Remove(entity Entity)
	RemoveAndReturn(entity Entity) Component
	EntityIndexes() []EntityIndex
	get(index EntityIndex) unsafe.Pointer
	Sort()
}

type CSet[T ComponentObject] struct {
	SparseArray[EntityIndex, T]
}

func NewCSet[T ComponentObject](initSize ...int) *CSet[T] {
	c := &CSet[T]{
		SparseArray: *NewSparseArray[EntityIndex, T](initSize...),
	}
	return c
}

func (c *CSet[T]) EntityIndexes() []EntityIndex {
	return *(*[]EntityIndex)(unsafe.Pointer(&c.idx2Key))
}

func (c *CSet[T]) Add(entity Entity, comp Component) Component {
	index := entity.Index()
	data := c.SparseArray.Add(index, (*T)((*iface)(unsafe.Pointer(&comp)).data))
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T]) remove(entity Entity) *T {
	index := entity.Index()
	return c.SparseArray.Remove(index)
}

func (c *CSet[T]) Remove(entity Entity) {
	data := c.remove(entity)
	if data == nil {
		return
	}
}

func (c *CSet[T]) RemoveAndReturn(entity Entity) Component {
	cpy := *c.remove(entity)
	return any(&cpy).(Component)
}

func (c *CSet[T]) Get(entity Entity) Component {
	data := c.SparseArray.Get(entity.Index())
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T]) get(entityIndex EntityIndex) unsafe.Pointer {
	data := c.SparseArray.Get(entityIndex)
	if data == nil {
		return nil
	}
	return unsafe.Pointer(data)
}

func (c *CSet[T]) getByIndex(index EntityIndex) *T {
	return c.SparseArray.Get(index)
}
