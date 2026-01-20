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

type CSet[T any, TP ComponentPointer[T]] struct {
	SparseArray[EntityIndex, T]
}

func NewCSet[T any, TP ComponentPointer[T]](initSize ...int) *CSet[T, TP] {
	c := &CSet[T, TP]{
		SparseArray: *NewSparseArray[EntityIndex, T](initSize...),
	}
	return c
}

func (c *CSet[T, TP]) EntityIndexes() []EntityIndex {
	return *(*[]EntityIndex)(unsafe.Pointer(&c.idx2Key))
}

func (c *CSet[T, TP]) Add(entity Entity, comp Component) Component {
	index := entity.Index()
	data := c.SparseArray.Add(index, (*T)((*iface)(unsafe.Pointer(&comp)).data))
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T, TP]) remove(entity Entity) *T {
	index := entity.Index()
	return c.SparseArray.Remove(index)
}

func (c *CSet[T, TP]) Remove(entity Entity) {
	data := c.remove(entity)
	if data == nil {
		return
	}
}

func (c *CSet[T, TP]) RemoveAndReturn(entity Entity) Component {
	cpy := *c.remove(entity)
	return any(&cpy).(Component)
}

func (c *CSet[T, TP]) Get(entity Entity) Component {
	data := c.SparseArray.Get(entity.Index())
	if data == nil {
		return nil
	}
	return any(data).(Component)
}

func (c *CSet[T, TP]) get(entityIndex EntityIndex) unsafe.Pointer {
	data := c.SparseArray.Get(entityIndex)
	if data == nil {
		return nil
	}
	return unsafe.Pointer(data)
}

func (c *CSet[T, TP]) getByIndex(index EntityIndex) *T {
	return c.SparseArray.Get(index)
}
