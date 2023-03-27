package ecs

type SparseArray[K Integer, V any] struct {
	UnorderedCollection[V]
	indices         []int32
	idx2Key         map[int32]int32
	maxKey          K
	sum             int64
	shrinkThreshold int32
}

func NewSparseArray[K Integer, V any](initSize ...int) *SparseArray[K, V] {
	initCap := 0
	if len(initSize) > 0 {
		initCap = initSize[0]
	}
	c := &SparseArray[K, V]{
		idx2Key: map[int32]int32{},
	}

	c.UnorderedCollection.init(true, initCap)

	if initCap > 0 {
		c.indices = make([]int32, 0, initCap)
	}

	switch any(*new(K)).(type) {
	case int8, uint8:
		c.shrinkThreshold = 127
	case uint16:
		c.shrinkThreshold = 255
	default:
		c.shrinkThreshold = 1024
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
		g.indices = append(g.indices, make([]int32, m, m)...)
	}

	g.idx2Key[int32(idx)] = int32(key)
	g.indices[key] = int32(idx + 1)
	if key > g.maxKey {
		g.maxKey = key
	}
	g.sum += int64(key)

	return &g.data[idx]
}

func (g *SparseArray[K, V]) Remove(key K) *V {
	if key > g.maxKey {
		return nil
	}
	idx := g.indices[key] - 1
	removed, oldIndex, newIndex := g.UnorderedCollection.Remove(int64(idx))

	lastKey := g.idx2Key[int32(oldIndex)]
	g.indices[lastKey] = int32(newIndex + 1)
	g.indices[key] = 0
	g.idx2Key[idx] = lastKey
	delete(g.idx2Key, int32(oldIndex))

	g.sum -= int64(key)

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

func (g *SparseArray[K, V]) Clear() {
	if g.Len() == 0 {
		return
	}
	g.UnorderedCollection.Reset()
	if int(g.maxKey) < 1024 {
		for i := 0; i < len(g.indices); i++ {
			g.indices[i] = 0
		}
	} else {
		g.indices = make([]int32, 0, g.initSize)
	}
	g.maxKey = 0
	g.sum = 0
	g.idx2Key = map[int32]int32{}
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

	if len(g.indices) > 1024 && int(g.maxKey) < len(g.indices)/2 && int64(g.maxKey) > g.sum/int64(g.len)*2 {
		m := (g.maxKey + 1) * 5 / 4
		temp := g.indices
		g.indices = make([]int32, g.maxKey+1, m)
		count := copy(g.indices, temp[:g.maxKey+1])
		if count != int(m) {
			panic("copy failed")
		}
	}
}
