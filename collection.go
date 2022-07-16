package ecs

import (
	"reflect"
	"sort"
	"unsafe"
)

const (
	InitMaxSize        = 1024 * 16
	SeqMax      uint32 = 0xFFFFFFFF
)

type ICollection interface {
	Len() int
	Range(func(v any) bool)
	ChangeCount() int64
	ChangeReset()
	ElementType() reflect.Type
	ElementMeta() ComponentMetaInfo
	Sort()

	getByIndex(idx int64) any
}

type Collection[T ComponentObject] struct {
	data    []T
	ids     map[int64]int64
	idx2id  map[int64]int64
	eleSize uintptr
	seq     int64
	len     int64
	meta    ComponentMetaInfo
	change  int64
}

func NewCollection[T ComponentObject]() *Collection[T] {
	typ := TypeOf[T]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	c := &Collection[T]{
		ids:     map[int64]int64{},
		idx2id:  map[int64]int64{},
		data:    make([]T, 0, size),
		eleSize: eleSize,
		meta:    ComponentMeta.GetComponentMetaInfo(typ),
	}
	return c
}

func (c *Collection[T]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *Collection[T]) Add(element *T, elementID ...int64) (*T, int64) {
	if int64(len(c.data)) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	idx := c.len
	var id int64
	if len(elementID) > 0 {
		id = elementID[0]
	} else {
		id = c.getID()
	}
	c.ids[id] = idx
	c.idx2id[idx] = id
	c.len++
	c.change++
	return &c.data[idx], id
}

func (c *Collection[T]) Remove(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	lastIdx := c.len - 1
	lastId := c.idx2id[lastIdx]

	c.ids[lastId] = idx
	c.idx2id[idx] = lastId
	delete(c.idx2id, lastIdx)
	delete(c.ids, id)

	c.data[idx], c.data[lastIdx] = c.data[lastIdx], c.data[idx]
	c.shrink()
	c.len--
	c.change++
	return &c.data[lastIdx]
}

func (c *Collection[T]) RemoveAndReturn(id int64) *T {
	cpy := *c.Remove(id)
	return &cpy
}

func (c *Collection[T]) shrink() {
	var threshold int64
	if len(c.data) < 1024 {
		threshold = c.len * 2
	} else {
		threshold = int64(float64(c.len) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		c.data = c.data[:threshold]
	}
}

func (c *Collection[T]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	base := uintptr(unsafe.Pointer(&c.data[0]))
	return (*T)(unsafe.Pointer(base + uintptr(idx)*c.eleSize))
}

func (c *Collection[T]) getByIndex(idx int64) any {
	base := uintptr(unsafe.Pointer(&c.data[0]))
	return (*T)(unsafe.Pointer(base + uintptr(idx)*c.eleSize))
}

func (c *Collection[T]) Sort() {
	if c.change == 0 {
		return
	}
	var zeroSeq = SeqMax
	seq2id := map[uint32]int64{}
	var cp *Component[T]
	for i := int64(0); i < c.len; i++ {
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
	for i := int64(0); i < c.len; i++ {
		seq = (*Component[T])(unsafe.Pointer(&(c.data[i]))).seq
		id = seq2id[seq]
		c.ids[id] = i
		c.idx2id[i] = id
	}
	c.change = 0
}

func (c *Collection[T]) Len() int {
	return int(c.len)
}

func (c *Collection[T]) ChangeCount() int64 {
	return c.change
}

func (c *Collection[T]) ChangeReset() {
	c.change = 0
}

func (c *Collection[T]) ElementType() reflect.Type {
	return TypeOf[T]()
}

func (c *Collection[T]) getData() []T {
	return c.data
}

func (c *Collection[T]) ElementMeta() ComponentMetaInfo {
	return c.meta
}

func (c *Collection[T]) Range(f func(element any) bool) {
	for i := int64(0); i < c.len; i++ {
		if !f(&c.data[i]) {
			break
		}
	}
}
