package ecs

import (
	"unsafe"
)

const (
	ChunkSize = 1024 * 16
)

const (
	ChunkAddCodeSuccess int = iota
	ChunkAddCodeFull
	ChunkAddCodeInvalidElement
)

type Chunk struct {
	data    [ChunkSize]byte
	seq     int64
	len     int64
	eleSize uintptr
	pend    uintptr
	pre     *Chunk
	next    *Chunk
}

func NewChunk(eleSize uintptr) *Chunk {
	c := &Chunk{
		eleSize: eleSize,
	}
	return c
}

func (c *Chunk) Add(element unsafe.Pointer) (unsafe.Pointer, int) {
	b := unsafe.Slice((*byte)(element), c.eleSize)
	if ChunkSize >= int64(c.pend+c.eleSize) {
		copy(c.data[c.pend+1:], b)
	} else {
		return nil, ChunkAddCodeFull
	}
	idx := c.len
	real := unsafe.Pointer(&(c.data[uintptr(idx)*c.eleSize]))
	c.len++
	c.pend += c.eleSize
	return real, ChunkAddCodeSuccess
}

func (c *Chunk) AddDiscrete(element []unsafe.Pointer, size []uintptr) (unsafe.Pointer, int) {
	if ChunkSize >= int64(c.pend+c.eleSize) {
		off := uintptr(0)
		for i, e := range element {
			b := unsafe.Slice((*byte)(e), size[i])
			copy(c.data[c.pend+off+1:], b)
			off += size[i]
		}
	} else {
		return nil, ChunkAddCodeFull
	}

	idx := c.len
	real := unsafe.Pointer(&(c.data[uintptr(idx)*c.eleSize]))
	c.len++
	c.pend += c.eleSize
	return real, ChunkAddCodeSuccess
}

func (c *Chunk[T, TP]) Remove(idx int64) {
	offset := idx * int64(c.eleSize)
	lastOffset := c.len * int64(c.eleSize)
	copy(c.data[offset:offset+int64(c.eleSize)], c.data[lastOffset:lastOffset+int64(c.eleSize)])
	c.len--
}

func (c *Chunk[T, TP]) Get(idx int64) unsafe.Pointer {
	if idx < 0 {
		return nil
	}
	return unsafe.Pointer(&(c.data[idx*int64(c.eleSize)]))
}

func (c *Chunk[T, TP]) Len() int {
	return int(c.len)
}
