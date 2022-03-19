package ecs

import (
	"reflect"
	"unsafe"
)

type CollectionByte[T ComponentObject, TP ComponentPointer[T]] struct {
	data    []byte
	ids     map[int64]int64
	seq     int64
	len     int64
	eleSize uintptr
	pend    uintptr
}

func NewCollectionByte[T ComponentObject, TP ComponentPointer[T]]() *CollectionByte[T, TP] {
	eleSize := unsafe.Sizeof(*new(T))
	c := &CollectionByte[T, TP]{
		ids:     map[int64]int64{},
		eleSize: eleSize,
		data:    make([]byte, 0, eleSize),
	}
	return c
}

func (c *CollectionByte[T, TP]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *CollectionByte[T, TP]) Add(element *T) (int64, *T) {
	p := (*byte)(unsafe.Pointer(element))
	b := unsafe.Slice(p, c.eleSize)
	if int64(len(c.data)) >= int64(c.pend+c.eleSize) {
		copy(c.data[c.pend+1:], b)
	} else {
		c.data = append(c.data, b...)
	}
	idx := c.len
	id := c.getID()
	c.ids[id] = idx
	c.ids[-idx] = -id
	real := (*T)(unsafe.Pointer(&(c.data[uintptr(idx)*c.eleSize])))
	c.len++
	c.pend += c.eleSize
	TP(real).setID(id)
	return id, real
}

func (c *CollectionByte[T, TP]) Remove(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	//Log.Info("collection Remove:", ObjectToString(c.data[idx]))
	lastIdx := c.len - 1
	lastId := -c.ids[-lastIdx]

	c.ids[lastId] = idx
	c.ids[-idx] = -lastId
	delete(c.ids, -lastIdx)
	delete(c.ids, id)

	offset := idx * int64(c.eleSize)
	lastOffset := c.len * int64(c.eleSize)
	r := *(*T)(unsafe.Pointer(&(c.data[lastOffset])))
	copy(c.data[offset:offset+int64(c.eleSize)], c.data[lastOffset:lastOffset+int64(c.eleSize)])
	c.shrink()
	c.len--
	return &r
}

func (c *CollectionByte[T, TP]) shrink() {
	var threshold int64
	if len(c.data) < 1024*int(c.eleSize) {
		threshold = int64(c.pend) * 2
	} else {
		threshold = int64(float64(c.pend) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		c.data = c.data[:threshold]
	}
}

func (c *CollectionByte[T, TP]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return (*T)(unsafe.Pointer(&(c.data[idx*int64(c.eleSize)])))
}

func (c *CollectionByte[T, TP]) GetByIndex(idx int64) *T {
	//if idx < 0 || idx >= c.len {
	//	return nil
	//}
	return (*T)(unsafe.Pointer(&(c.data[idx*int64(c.eleSize)])))
}

func (c *CollectionByte[T, TP]) Len() int {
	return int(c.len)
}

func (c *CollectionByte[T, TP]) ElementType() reflect.Type {
	return TypeOf[T]()
}
