package ecs

import (
	"reflect"
	"unsafe"
)

type Collection2[T ComponentObject, TP ComponentPointer[T]] struct {
	data       *Chunk[T, TP]
	ids        map[int64]*Chunk[T, TP]
	pend       *Chunk[T, TP]
	holeList   map[*Chunk[T, TP]]struct{}
	eleSize    uintptr
	seq        int64
	len        int64
	chunkCount int64
}

func NewCollection2[T ComponentObject, TP ComponentPointer[T]]() *Collection2[T, TP] {
	var e T
	size := unsafe.Sizeof(e)
	c := &Collection2[T, TP]{
		ids:        map[int64]*Chunk[T, TP]{},
		data:       NewChunk[T, TP](),
		eleSize:    size,
		holeList:   map[*Chunk[T, TP]]struct{}{},
		chunkCount: 1,
	}
	c.pend = c.data
	return c
}

func (c *Collection2[T, TP]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *Collection2[T, TP]) Add(element *T) *T {
	id := TP(element).ID()
	if id == 0 {
		id = c.getID()
		TP(element).setID(id)
	}
	var p *T
	var code int
	if len(c.holeList) > 0 {
		// 优先放置到有空洞的chunk中
		for chunk := range c.holeList {
			p, code = chunk.Add(element, id)
			if code != 0 {
				if code == 1 {
					delete(c.holeList, chunk)
					continue
				}
			}
		}
		if p == nil {
			return c.Add(element)
		}
	} else {
		p, code = c.pend.Add(element, id)
		if code != 0 {
			// 扩容
			if code == 1 {
				nt := NewChunk[T, TP]()
				c.pend.next = nt
				nt.pre = c.pend
				c.pend = nt
				c.chunkCount++
				return c.Add(element)
			}
		}
	}
	c.ids[id] = c.pend
	c.len++
	return (*T)(p)
}

func (c *Collection2[T, TP]) Remove(id int64) *T {
	chunk, ok := c.ids[id]
	if !ok {
		return nil
	}
	p := chunk.RemoveAndReturn(id)
	c.len--
	delete(c.ids, id)
	if chunk != c.pend {
		if _, ok := c.holeList[chunk]; !ok {
			c.holeList[chunk] = struct{}{}
		}
	}
	c.shrink()
	return (*T)(p)
}

func (c *Collection2[T, TP]) shrink() {
	chunkMaxSize := int64(ChunkSize / c.eleSize)
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

func (c *Collection2[T, TP]) shrinkImmediately() {
	for chunk, _ := range c.holeList {
		c.tidy(chunk)
	}
}

func (c *Collection2[T, TP]) tidy(chunk *Chunk[T, TP]) {
	length := chunk.Len()
	entities := chunk.MoveTo(c.pend)
	for _, id := range entities {
		c.ids[id] = c.pend
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
func (c *Collection2[T, TP]) optimize() {
	c.shrinkImmediately()
}

func (c *Collection2[T, TP]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	chunk, ok := c.ids[id]
	if !ok {
		return nil
	}
	return (*T)(chunk.Get(id))
}

func (c *Collection2[T, TP]) Len() int {
	return int(c.len)
}

func (c *Collection2[T, TP]) ElementType() reflect.Type {
	return TypeOf[T]()
}
