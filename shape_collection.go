package ecs

import (
	"reflect"
	"unsafe"
)

type ChunkIndex struct {
	index int64
	chunk *Chunk
}

type ShapeCollection[T ShapeObject, TP ShapePointer[T]] struct {
	data    *Chunk
	ids     map[int64]ChunkIndex
	pend    *Chunk
	eleSize []uintptr
	seq     int64
}

func NewShapeCollection[T ShapeObject, TP ShapePointer[T]](eleSize []uintptr) *ShapeCollection[T, TP] {
	chunkEleSize := uintptr(0)
	for _, size := range eleSize {
		chunkEleSize += size
	}
	c := &ShapeCollection[T, TP]{
		ids:  map[int64]ChunkIndex{},
		data: NewChunk(chunkEleSize),
	}
	return c
}

func (c *ShapeCollection[T, TP]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *ShapeCollection[T, TP]) Add(shape TP, entity Entity) (int64, *T) {
	elements := shape.getElements()
	c.pend.AddDiscrete(elements, c.eleSize)
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

// remove 非常复杂，当一个shape被移除某个组件之后，剩下的新组合需要重新添加到对应的shape容器中
func (c *ShapeCollection[T, TP]) Remove(id int64) *T {
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

func (c *ShapeCollection[T, TP]) shrink() {
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

func (c *ShapeCollection[T, TP]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return (*T)(unsafe.Pointer(&(c.data[idx*int64(c.eleSize)])))
}

func (c *ShapeCollection[T, TP]) GetByIndex(idx int64) *T {
	//if idx < 0 || idx >= c.len {
	//	return nil
	//}
	return (*T)(unsafe.Pointer(&(c.data[idx*int64(c.eleSize)])))
}

func (c *ShapeCollection[T, TP]) Len() int {
	return int(c.len)
}

func (c *ShapeCollection[T, TP]) ElementType() reflect.Type {
	return TypeOf[T]()
}
