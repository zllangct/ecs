package ecs

type SparseArray[K Integer, V any] struct {
	UnorderedCollection[V]
	indices         []int32
	idx2Key         map[int32]int32
	maxKey          K
	shrinkThreshold int32
}

func NewSparseArray[K Integer, V any](initSize ...int) *SparseArray[K, V] {
	typ := TypeOf[V]()
	eleSize := typ.Size()
	size := InitMaxSize / eleSize
	if len(initSize) > 0 {
		size = uintptr(initSize[0]) / eleSize
	}
	c := &SparseArray[K, V]{
		UnorderedCollection: UnorderedCollection[V]{
			data:    make([]V, 0, size),
			eleSize: eleSize,
		},
		idx2Key: map[int32]int32{},
	}
	switch any(*new(K)).(type) {
	case int8, uint8:
		c.shrinkThreshold = 127
	case uint16:
		c.shrinkThreshold = 255
	default:
		c.shrinkThreshold = 1024
	}

	if size == 0 {
		c.indices = make([]int32, 1)
	}
	return c
}

func (g *SparseArray[K, V]) Add(key K, value *V) *V {
	length := len(g.indices)
	// already existed
	if key < K(length) && g.indices[key] != 0 {
		return nil
	}
	_, idx := g.UnorderedCollection.Add(value)
	if key >= K(length) {
		m := K(0)
		if length == 0 {
			m = key + 1
		} else if length < int(g.shrinkThreshold) {
			m = key * 2
		} else {
			m = key * 5 / 4
		}
		newIndices := make([]int32, m)
		count := copy(newIndices, g.indices)
		if count != length {
			panic("copy failed")
		}
		g.indices = newIndices
	}

	g.idx2Key[int32(idx)] = int32(key)
	g.indices[key] = int32(idx + 1)
	if key > g.maxKey {
		g.maxKey = key
	}

	return &g.data[idx]
}

func (g *SparseArray[K, V]) Remove(key K) *V {
	if key > g.maxKey {
		return nil
	}
	idx := g.indices[key] - 1
	removed, oldIndex, newIndex := g.UnorderedCollection.Remove(int64(idx))

	delete(g.idx2Key, idx)
	g.indices[g.idx2Key[int32(oldIndex)]] = int32(newIndex + 1)
	g.indices[key] = 0

	g.shrink(key)

	return removed
}

func (g *SparseArray[K, V]) Exist(key K) bool {
	if key > g.maxKey {
		return false
	}
	return !(g.indices[key] == 0)
}

func (g *SparseArray[K, V]) Get(key K) *V {
	if key > g.maxKey {
		return nil
	}
	idx := g.indices[key] - 1
	if idx < 0 {
		return nil
	}
	return g.UnorderedCollection.Get(int64(idx))
}

func (g *SparseArray[K, V]) shrink(key K) {
	if key < g.maxKey {
		return
	}

	g.maxKey = 0
	for i := key; i > 0; i-- {
		if g.indices[i] != 0 {
			g.maxKey = i
			break
		}
	}

	if int32(g.maxKey) < g.shrinkThreshold {
		g.maxKey = K(g.shrinkThreshold)
	}

	if len(g.indices) > 1024 && int(g.maxKey) < len(g.indices)/2 {
		m := (g.maxKey + 1) * 5 / 4
		newIndices := make([]int32, m)
		count := copy(newIndices, g.indices[:m])
		if count != int(m) {
			panic("copy failed")
		}
	}
}
