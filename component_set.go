package ecs

import (
	"sort"
	"unsafe"
)

type IComponentSet interface {
	Len() int
	Range(fn func(com IComponent) bool)
	Clear()
	GetByEntity(entity Entity) any
	GetElementMeta() *ComponentMetaInfo
	GetComponent(entity Entity) IComponent
	GetComponentRaw(entity Entity) unsafe.Pointer
	Remove(entity Entity)
	Sort()

	getPointerByIndex(index int64) unsafe.Pointer
	changeCount() int64
	changeReset()
	pointer() unsafe.Pointer
	getPointerByEntity(entity Entity) unsafe.Pointer
}

type ComponentSet[T ComponentObject] struct {
	SparseArray[int32, T]
	change int64
	meta   *ComponentMetaInfo
}

func NewComponentSet[T ComponentObject](meta *ComponentMetaInfo, initSize ...int) *ComponentSet[T] {
	c := &ComponentSet[T]{
		SparseArray: *NewSparseArray[int32, T](initSize...),
		meta:        meta,
	}
	return c
}

func (c *ComponentSet[T]) Add(element *T, entity Entity) *T {
	index := entity.ToRealID().index
	data := c.SparseArray.Add(index, element)
	if data == nil {
		return nil
	}
	c.change++
	return data
}

func (c *ComponentSet[T]) remove(entity Entity) *T {
	index := entity.ToRealID().index
	return c.SparseArray.Remove(index)
}

func (c *ComponentSet[T]) Remove(entity Entity) {
	data := c.remove(entity)
	if data == nil {
		return
	}
	c.change++
}

func (c *ComponentSet[T]) RemoveAndReturn(entity Entity) *T {
	cpy := *c.remove(entity)
	return &cpy
}

func (c *ComponentSet[T]) getByEntity(entity Entity) *T {
	return c.SparseArray.Get(entity.ToRealID().index)
}

func (c *ComponentSet[T]) getPointerByEntity(entity Entity) unsafe.Pointer {
	return unsafe.Pointer(c.getByEntity(entity))
}

func (c *ComponentSet[T]) GetByEntity(entity Entity) any {
	return c.getByEntity(entity)
}

func (c *ComponentSet[T]) Get(entity Entity) *T {
	return c.getByEntity(entity)
}

func (c *ComponentSet[T]) pointer() unsafe.Pointer {
	return unsafe.Pointer(c)
}

func (c *ComponentSet[T]) changeCount() int64 {
	return c.change
}

func (c *ComponentSet[T]) changeReset() {
	c.change = 0
}

func (c *ComponentSet[T]) Sort() {
	if c.changeCount() == 0 {
		return
	}
	var zeroSeq = SeqMax
	seq2id := map[uint32]int64{}
	var cp *Component[T]
	for i := int64(0); i < int64(c.Len()); i++ {
		cp = (*Component[T])(unsafe.Pointer(&(c.data[i])))
		if cp.seq == 0 {
			zeroSeq--
			cp.seq = zeroSeq
		}
		seq2id[cp.seq] = cp.owner.ToInt64()
	}
	sort.Slice(c.data, func(i, j int) bool {
		return (*Component[T])(unsafe.Pointer(&(c.data[i]))).seq < (*Component[T])(unsafe.Pointer(&(c.data[j]))).seq
	})
	for i := int32(0); i < int32(c.Len()); i++ {
		cp = (*Component[T])(unsafe.Pointer(&(c.data[i])))
		c.indices[cp.owner.ToRealID().index] = i + 1
	}
	c.changeReset()
}

func (c *ComponentSet[T]) GetComponent(entity Entity) IComponent {
	return c.GetByEntity(entity).(IComponent)
}

func (c *ComponentSet[T]) GetComponentRaw(entity Entity) unsafe.Pointer {
	return unsafe.Pointer(c.getByEntity(entity))
}

func (c *ComponentSet[T]) getPointerByIndex(index int64) unsafe.Pointer {
	return unsafe.Pointer(c.SparseArray.UnorderedCollection.Get(index))
}

func (c *ComponentSet[T]) GetElementMeta() *ComponentMetaInfo {
	return c.meta
}

func (c *ComponentSet[T]) Range(fn func(com IComponent) bool) {
	c.SparseArray.Range(func(com *T) bool {
		return fn(any(com).(IComponent))
	})
}

func NewComponentSetIterator[T ComponentObject](collection *ComponentSet[T], readOnly ...bool) Iterator[T] {
	iter := &Iter[T]{
		data:    collection.data,
		len:     collection.Len(),
		eleSize: collection.eleSize,
		offset:  0,
	}
	if len(readOnly) > 0 {
		iter.readOnly = readOnly[0]
	}
	if iter.len != 0 {
		iter.head = unsafe.Pointer(&collection.data[0])
		if iter.readOnly {
			iter.curTemp = collection.data[0]
			iter.cur = &iter.curTemp
		} else {
			iter.cur = &(collection.data[0])
		}
	}

	return iter
}
