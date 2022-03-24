package ecs

import (
	"unsafe"
)

const (
	//ChunkSize int64 = 1024 * 16
	ChunkSize  int64   = 512
	EntitySize uintptr = 8
)

const (
	ChunkAddCodeSuccess int = iota
	ChunkAddCodeFull
	ChunkAddCodeInvalidElement
)

type Chunk struct {
	data    [ChunkSize]byte
	ids     map[int64]int64
	len     int64
	eleSize uintptr
	pend    uintptr
	pre     *Chunk
	next    *Chunk
}

func NewChunk(eleSize uintptr) *Chunk {
	c := &Chunk{
		eleSize: eleSize,
		ids:     make(map[int64]int64),
	}
	return c
}

func (c *Chunk) MaxLength() int64 {
	return ChunkSize / int64(c.eleSize)
}

func (c *Chunk) Add(element unsafe.Pointer, entity Entity) (unsafe.Pointer, int) {
	bElement := unsafe.Slice((*byte)(element), c.eleSize)
	if ChunkSize >= int64(c.pend+c.eleSize) {
		copy(c.data[c.pend:], bElement)
	} else {
		return nil, ChunkAddCodeFull
	}
	idx := c.len
	c.ids[int64(entity)] = idx
	c.ids[-idx] = int64(-entity)
	real := unsafe.Pointer(&(c.data[c.pend]))
	c.len++
	c.pend += c.eleSize
	return real, ChunkAddCodeSuccess
}

func (c *Chunk) Remove(entity Entity) {
	idx, ok := c.ids[int64(entity)]
	if !ok {
		return
	}
	lastIdx := c.len - 1
	lastId := -c.ids[-lastIdx]
	c.ids[lastId] = idx
	c.ids[-idx] = -lastId
	delete(c.ids, -lastIdx)
	delete(c.ids, int64(entity))

	offset := uintptr(idx) * c.eleSize
	lastOffset := uintptr(lastIdx) * c.eleSize
	copy(c.data[offset:offset+c.eleSize], c.data[lastOffset:lastOffset+c.eleSize])
	c.len--
}

func (c *Chunk) RemoveAndReturn(entity Entity) unsafe.Pointer {
	idx, ok := c.ids[int64(entity)]
	if !ok {
		return nil
	}
	lastIdx := c.len - 1
	lastId := -c.ids[-lastIdx]
	c.ids[lastId] = idx
	c.ids[-idx] = -lastId
	delete(c.ids, -lastIdx)
	delete(c.ids, int64(entity))
	offset := uintptr(idx) * c.eleSize
	lastOffset := uintptr(lastIdx) * c.eleSize
	var r = make([]byte, c.eleSize, c.eleSize)
	copy(r, c.data[offset:offset+c.eleSize])
	copy(c.data[offset:offset+c.eleSize], c.data[lastOffset:lastOffset+c.eleSize])
	c.len--
	c.pend -= c.eleSize
	return unsafe.Pointer(&r[0])
}

func (c *Chunk) MoveTo(target *Chunk) []Entity {
	available := target.MaxLength() - target.len
	moveSize := uintptr(0)
	if c.len < available {
		moveSize = uintptr(c.len)
	} else {
		moveSize = uintptr(available)
	}
	copy(target.data[target.pend:target.pend+moveSize*c.eleSize], c.data[c.pend-moveSize*c.eleSize:c.pend])

	var entities []Entity
	for i := int64(0); i < int64(moveSize); i++ {
		idx := c.len - int64(moveSize) + i
		entity := Entity(c.ids[-idx])
		entities = append(entities, entity)
	}

	target.pend += moveSize * c.eleSize
	c.pend -= moveSize * c.eleSize
	target.len += int64(moveSize)
	c.len -= int64(moveSize)

	return entities
}

func (c *Chunk) Get(entity Entity) unsafe.Pointer {
	idx, ok := c.ids[int64(entity)]
	if !ok {
		return nil
	}
	return unsafe.Pointer(&(c.data[idx*int64(c.eleSize)]))
}

func (c *Chunk) GetByIndex(idx int64) unsafe.Pointer {
	if idx < 0 || idx >= c.len {
		return nil
	}
	return unsafe.Pointer(&(c.data[idx*int64(c.eleSize)]))
}

func (c *Chunk) Len() int64 {
	return c.len
}
