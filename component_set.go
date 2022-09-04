package ecs

import (
	"sort"
	"unsafe"
)

type IComponentSet interface {
	ICollection
	GetByEntity(entity Entity) any
	GetElementMeta() ComponentMetaInfo
	GetComponent(entity Entity) IComponent

	Sort()
	Shrink()
	pointer() unsafe.Pointer
	getPointerByEntity(entity Entity) unsafe.Pointer
}

type ComponentSet[T ComponentObject] struct {
	UnorderedCollection[T]
	indexMax int32
	indices  []int32
	meta     ComponentMetaInfo
}

func NewComponentSet[T ComponentObject](meta ComponentMetaInfo, initSize ...int) *ComponentSet[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	if len(initSize) > 0 {
		size = uintptr(initSize[0]) / eleSize
	}
	c := &ComponentSet[T]{
		UnorderedCollection: UnorderedCollection[T]{
			data:    make([]T, 0, size),
			eleSize: eleSize,
		},
		meta: meta,
	}
	if size == 0 {
		c.indices = make([]int32, 1)
	}
	return c
}

func (c *ComponentSet[T]) Add(element *T, entity Entity) *T {
	realID := entity.ToRealID()
	_, idx := c.UnorderedCollection.Add(element)
	length := len(c.indices)
	if realID.index >= int32(length) {
		m := int32(0)
		if length == 0 {
			m = realID.index + 1
		} else if length < 1024 {
			m = realID.index * 2
		} else {
			m = realID.index * 5 / 4
		}
		newIndices := make([]int32, m)
		count := copy(newIndices, c.indices)
		if count != length {
			panic("copy failed")
		}
		c.indices = newIndices
	}

	c.indices[realID.index] = int32(idx)

	return &c.data[idx]
}

func (c *ComponentSet[T]) remove(entity Entity) *T {
	realID := entity.ToRealID()
	if realID.index >= int32(len(c.indices)) {
		return nil
	}
	idx := c.indices[realID.index]
	removed, oldIndex, newIndex := c.UnorderedCollection.Remove(int64(idx))
	oldRealID := c.data[oldIndex].OwnerEntity().ToRealID()
	c.indices[oldRealID.index] = int32(newIndex)
	c.indices[realID.index] = -1

	return removed
}

func (c *ComponentSet[T]) Remove(entity Entity) {
	c.remove(entity)
}

func (c *ComponentSet[T]) RemoveAndReturn(entity Entity) *T {
	cpy := *c.remove(entity)
	return &cpy
}

func (c *ComponentSet[T]) getByEntity(entity Entity) *T {
	idx := c.indices[entity.ToRealID().index]
	if idx < 0 {
		return nil
	}
	return c.UnorderedCollection.Get(int64(idx))
}

func (c *ComponentSet[T]) getPointerByEntity(entity Entity) unsafe.Pointer {
	idx := c.indices[entity.ToRealID().index]
	if idx < 0 {
		return nil
	}
	return c.UnorderedCollection.getPointer(int64(idx))
}

func (c *ComponentSet[T]) GetByEntity(entity Entity) any {
	return c.getByEntity(entity)
}

func (c *ComponentSet[T]) pointer() unsafe.Pointer {
	return unsafe.Pointer(c)
}

func (c *ComponentSet[T]) Shrink() {
	// TODO indices缩容
}

func (c *ComponentSet[T]) Sort() {
	if c.ChangeCount() == 0 {
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
		c.indices[cp.owner.ToRealID().index] = i
	}
	c.ChangeReset()
}

func (c *ComponentSet[T]) GetComponent(entity Entity) IComponent {
	return c.GetByEntity(entity).(IComponent)
}

func (c *ComponentSet[T]) GetElementMeta() ComponentMetaInfo {
	return c.meta
}

func NewComponentSetIterator[T ComponentObject](collection *ComponentSet[T], readOnly ...bool) Iterator[T] {
	iter := &Iter[T]{
		data:    collection.getData(),
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
