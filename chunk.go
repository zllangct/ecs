package ecs

const (
	ChunkSize uintptr = 1024 * 16
	//ChunkSize  int64   = 512
	EntitySize uintptr = 8
)

const (
	ChunkAddCodeSuccess int = iota
	ChunkAddCodeFull
	ChunkAddCodeInvalidElement
)

type Chunk[T ComponentObject, TP ComponentPointer[T]] struct {
	data    []T
	ids     map[int64]int64
	idx2id  map[int64]int64
	len     int64
	max     int64
	eleSize uintptr
	pend    uintptr
	pre     *Chunk[T, TP]
	next    *Chunk[T, TP]
}

func NewChunk[T ComponentObject, TP ComponentPointer[T]]() *Chunk[T, TP] {
	size := TypeOf[T]().Size()
	max := ChunkSize / size
	c := &Chunk[T, TP]{
		data:    make([]T, max, max),
		eleSize: size,
		max:     int64(max),
		ids:     make(map[int64]int64),
		idx2id:  make(map[int64]int64),
	}
	return c
}

func (c *Chunk[T, TP]) Add(element *T, id int64) (*T, int) {
	if uintptr(len(c.data)) >= c.pend+c.eleSize {
		c.data[c.pend+1] = *element
	} else {
		return nil, ChunkAddCodeFull
	}
	idx := c.len
	c.ids[id] = idx
	c.idx2id[idx] = id
	real := &(c.data[c.pend])
	c.len++
	c.pend += c.eleSize
	return real, ChunkAddCodeSuccess
}

func (c *Chunk[T, TP]) Remove(id int64) {
	if id < 0 {
		return
	}
	idx, ok := c.ids[id]
	if !ok {
		return
	}
	lastIdx := c.len - 1
	lastId := c.idx2id[lastIdx]

	c.ids[lastId] = idx
	c.idx2id[idx] = lastId
	delete(c.idx2id, lastIdx)
	delete(c.ids, id)

	c.data[idx], c.data[lastIdx] = c.data[lastIdx], c.data[idx]
	c.len--
}

func (c *Chunk[T, TP]) RemoveAndReturn(id int64) *T {
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
	r := &(c.data[lastIdx])
	c.len--
	c.pend -= c.eleSize
	return r
}

func (c *Chunk[T, TP]) MoveTo(target *Chunk[T, TP]) []int64 {
	moveSize := uintptr(0)
	if c.len < c.max {
		moveSize = uintptr(c.len)
	} else {
		moveSize = uintptr(c.max)
	}
	copy(target.data[target.pend:target.pend+moveSize], c.data[c.pend-moveSize:c.pend])

	var moved []int64
	for i := int64(0); i < int64(moveSize); i++ {
		idx := c.len - int64(moveSize) + i
		id := c.idx2id[idx]
		moved = append(moved, id)
	}

	target.pend += moveSize * c.eleSize
	c.pend -= moveSize * c.eleSize
	target.len += int64(moveSize)
	c.len -= int64(moveSize)

	return moved
}

func (c *Chunk[T, TP]) Get(id int64) *T {
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return &(c.data[idx])
}

func (c *Chunk[T, TP]) GetByIndex(idx int64) *T {
	if idx < 0 || idx >= c.len {
		return nil
	}
	return &(c.data[idx])
}

func (c *Chunk[T, TP]) Len() int64 {
	return c.len
}
