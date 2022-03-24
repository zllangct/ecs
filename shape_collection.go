package ecs

import (
	"unsafe"
)

type ChunkIndex struct {
	chunk *Chunk
}

type ShapeCollection[T ComponentObject, TP ComponentPointer[T]] struct {
	data       *Chunk
	ids        map[Entity]ChunkIndex
	pend       *Chunk
	holeList   map[*Chunk]struct{}
	eleSize    uintptr
	seq        int64
	len        int64
	chunkCount int64
}

func NewShapeCollection[T ComponentObject, TP ComponentPointer[T]]() *ShapeCollection[T, TP] {
	var e T
	size := unsafe.Sizeof(e)
	c := &ShapeCollection[T, TP]{
		ids:        map[Entity]ChunkIndex{},
		data:       NewChunk(size),
		eleSize:    size,
		holeList:   map[*Chunk]struct{}{},
		chunkCount: 1,
	}
	c.pend = c.data
	return c
}

func (c *ShapeCollection[T, TP]) Add(shape *T) *T {
	entity := shape.ID()
	if entity == 0 {
		return nil
	}
	var pElement = unsafe.Pointer(shape)
	var p unsafe.Pointer
	var code int
	if len(c.holeList) > 0 {
		// 优先放置到有空洞的chunk中
		for chunk := range c.holeList {
			p, code = chunk.Add(pElement, entity)
			if code != 0 {
				if code == 1 {
					delete(c.holeList, chunk)
					continue
				}
			}
		}
		if p == nil {
			return c.Add(shape)
		}
	} else {
		p, code = c.pend.Add(pElement, entity)
		if code != 0 {
			// 扩容
			if code == 1 {
				nt := NewChunk(c.eleSize)
				c.pend.next = nt
				nt.pre = c.pend
				c.pend = nt
				c.chunkCount++
				return c.Add(shape)
			}
		}
	}
	c.ids[entity] = ChunkIndex{
		chunk: c.pend,
	}
	c.len++
	return (*T)(p)
}

func (c *ShapeCollection[T, TP]) RemoveAndReturn(entity Entity) *T {
	index, ok := c.ids[entity]
	if !ok {
		return nil
	}
	p := index.chunk.RemoveAndReturn(entity)
	c.len--
	delete(c.ids, entity)
	if index.chunk != c.pend {
		if _, ok := c.holeList[index.chunk]; !ok {
			c.holeList[index.chunk] = struct{}{}
		}
	}
	c.shrink()
	return &*(*T)(p)
}

func (c *ShapeCollection[T, TP]) shrink() {
	chunkMaxSize := ChunkSize / int64(c.chunkEleSize)
	if c.len-chunkMaxSize*c.chunkCount < 2*chunkMaxSize {
		return
	}
	for chunk, _ := range c.holeList {
		length := chunk.Len()
		if length*2 < chunkMaxSize {
			c.tidy(chunk)
		}
	}
}

func (c *ShapeCollection[T, TP]) shrinkImmediately() {
	for chunk, _ := range c.holeList {
		c.tidy(chunk)
	}
}

func (c *ShapeCollection[T, TP]) tidy(chunk *Chunk) {
	length := chunk.Len()
	entities := chunk.MoveTo(c.pend)
	for _, entity := range entities {
		c.ids[entity] = ChunkIndex{
			chunk: c.pend,
		}
	}
	// 尾结点已满，该节点还有剩余，该节点搬迁到末尾
	if int64(len(entities)) < length {
		chunk.pre.next = chunk.next
		chunk.next.pre = chunk.pre
		c.pend.next = chunk
		chunk.pre = c.pend
		delete(c.holeList, chunk)
	}
}

// 非帧执行阶段优化，不影响占用主逻辑时间
func (c *ShapeCollection[T, TP]) optimize() {
	c.shrinkImmediately()
}

func (c *ShapeCollection[T, TP]) Get(entity Entity) *T {
	if entity < 0 {
		return nil
	}
	index, ok := c.ids[entity]
	if !ok {
		return nil
	}
	return (*T)(index.chunk.Get(entity))
}

func (c *ShapeCollection[T, TP]) Len() int {
	return int(c.len)
}
