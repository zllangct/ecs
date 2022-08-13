package ecs

import (
	"sort"
	"unsafe"
)

type IComponentSet interface {
	ICollectionWithID
	GetElementMeta() ComponentMetaInfo
	GetComponent(entity Entity) IComponent
	Sort()
}

type ComponentSet[T ComponentObject] struct {
	UnorderedCollectionWithID[T]
	meta ComponentMetaInfo
}

func NewComponentSet[T ComponentObject]() *ComponentSet[T] {
	typ := TypeOf[T]()
	return &ComponentSet[T]{
		UnorderedCollectionWithID: *NewUnorderedCollectionWithID[T](),
		meta:                      ComponentMeta.GetComponentMetaInfo(typ),
	}
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
		seq2id[cp.seq] = c.idx2id[i]
	}
	sort.Slice(c.data, func(i, j int) bool {
		return (*Component[T])(unsafe.Pointer(&(c.data[i]))).seq < (*Component[T])(unsafe.Pointer(&(c.data[j]))).seq
	})
	var id int64
	var seq uint32
	for i := int64(0); i < int64(c.Len()); i++ {
		seq = (*Component[T])(unsafe.Pointer(&(c.data[i]))).seq
		id = seq2id[seq]
		c.ids[id] = i
		c.idx2id[i] = id
	}
	c.ChangeReset()
}

func (c *ComponentSet[T]) GetComponent(entity Entity) IComponent {
	return c.UnorderedCollectionWithID.GetByID(int64(entity)).(IComponent)
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
