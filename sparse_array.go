package ecs

type SparseArray[K comparable, V any] struct {
	indices []K
	values  []V
}

func NewSparseArray[K comparable, V any](initCap ...int) *SparseArray[K, V] {
	cap := 0
	if len(initCap) > 0 {
		cap = initCap[0]
	}
	return &SparseArray[K, V]{
		indices: make([]K, 0, cap),
		values:  make([]V, 0, cap),
	}
}

func (g *SparseArray[K, V]) Add(key K, value V) {
	g.indices = append(g.indices, key)
	g.values = append(g.values, value)
}

func (g *SparseArray[K, V]) Remove(key K) {
	idx := -1
	for i, t := range g.indices {
		if t == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	// remove from indices
	g.indices = append(g.indices[:idx], g.indices[idx+1:]...)
	// remove from values
	g.values = append(g.values[:idx], g.values[idx+1:]...)
}

func (g *SparseArray[K, V]) Get(key K) (V, bool) {
	for i, t := range g.indices {
		if t == key {
			return g.values[i], true
		}
	}
	v := new(V)
	return *v, false
}
