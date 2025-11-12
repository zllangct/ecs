package ecs

import (
	"iter"
)

type SparseArray[K Integer, V any] struct {
	USet[V]
	indices         []int32
	idx2Key         []int32
	maxKey          K
	shrinkThreshold int32
	initSize        int
	isKOrder        bool
}

func NewSparseArray[K Integer, V any](initSize ...int) *SparseArray[K, V] {
	size := 0
	if len(initSize) > 0 {
		size = initSize[0]
	}
	c := &SparseArray[K, V]{
		idx2Key:  []int32{},
		initSize: int(size),
		USet:     *NewUSet[V](size),
		isKOrder: true,
	}

	if size > 0 {
		c.indices = make([]int32, 0, size)
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

func (s *SparseArray[K, V]) Add(key K, value *V) *V {
	length := len(s.indices)
	// already existed
	if key < K(length) && s.indices[key] != 0 {
		return nil
	}
	_, idx := s.USet.Add(value)
	if key >= K(length) {
		m := K(0)
		if length == 0 {
			m = key + 1
		} else if length < int(s.shrinkThreshold) {
			m = key * 2
		} else {
			m = key * 5 / 4
		}
		newIndices := make([]int32, m)
		count := copy(newIndices, s.indices)
		if count != length {
			panic("copy failed")
		}
		s.indices = newIndices
	}
	if int(idx) < len(s.idx2Key) {
		s.idx2Key[int32(idx)] = int32(key)
	} else {
		s.idx2Key = append(s.idx2Key, int32(key))
	}
	s.indices[key] = int32(idx + 1)
	if key > s.maxKey {
		s.maxKey = key
	}

	s.isKOrder = false

	return &s.data[idx]
}

func (s *SparseArray[K, V]) Remove(key K) *V {
	if key > s.maxKey {
		return nil
	}
	idx := s.indices[key] - 1
	removed, oldIndex, newIndex := s.USet.Remove(int64(idx))

	lastKey := s.idx2Key[int32(oldIndex)]
	s.indices[lastKey] = int32(newIndex + 1)
	s.indices[key] = 0
	s.idx2Key[idx] = lastKey

	// swap
	s.idx2Key[newIndex], s.idx2Key[oldIndex] = s.idx2Key[oldIndex], s.idx2Key[newIndex]
	// remove last
	s.idx2Key = s.idx2Key[:len(s.idx2Key)]

	s.shrink(key)

	s.isKOrder = false

	return removed
}

func (s *SparseArray[K, V]) Exist(key K) bool {
	if key > s.maxKey {
		return false
	}
	return !(s.indices[key] == 0)
}

func (s *SparseArray[K, V]) Get(key K) *V {
	if key > s.maxKey {
		return nil
	}
	idx := s.indices[key] - 1
	if idx < 0 {
		return nil
	}
	return s.USet.Get(int64(idx))
}

func (s *SparseArray[K, V]) Reset() {
	if s.Len() == 0 {
		return
	}
	s.USet.Reset()
	if int(s.maxKey) < 1024 {
		for i := 0; i < len(s.indices); i++ {
			s.indices[i] = 0
		}
	} else {
		s.indices = make([]int32, 0, s.initSize)
	}
	s.maxKey = 0
	s.idx2Key = []int32{}
	s.isKOrder = true
}

func (s *SparseArray[K, V]) Less(i, j int) bool {
	return s.idx2Key[i] < s.idx2Key[j]
}

func (s *SparseArray[K, V]) Swap(i, j int) {
	Key := s.idx2Key[j]
	s.USet.Swap(int64(i), int64(j))
	// swap
	s.idx2Key[i], s.idx2Key[j] = s.idx2Key[j], s.idx2Key[i]
	s.indices[Key], s.indices[i] = s.indices[i], s.indices[Key]
}

func (s *SparseArray[K, V]) Sort() {
	if s.isKOrder {
		return
	}
	seq := int64(0)
	for i, index := range s.indices {
		if index == 0 {
			continue
		}
		idx := index - 1
		if idx != s.idx2Key[seq] {
			Key := s.idx2Key[seq]
			s.USet.Swap(int64(idx), seq)
			// swap
			s.idx2Key[idx], s.idx2Key[seq] = s.idx2Key[seq], s.idx2Key[idx]
			s.indices[Key], s.indices[i] = s.indices[i], s.indices[Key]
		}
		seq++
	}
	s.isKOrder = true
}

func (s *SparseArray[K, V]) shrink(key K) {
	if key < s.maxKey {
		return
	}

	s.maxKey = 0
	for i := key; i > 0; i-- {
		if s.indices[i] != 0 {
			s.maxKey = i
			break
		}
	}

	if int32(s.maxKey) < s.shrinkThreshold {
		s.maxKey = K(s.shrinkThreshold)
	}

	if len(s.indices) > 1024 && int(s.maxKey) < len(s.indices)/2 {
		m := (s.maxKey + 1) * 5 / 4
		newIndices := make([]int32, m)
		copy(newIndices, s.indices[:m])
	}
}

func (s *SparseArray[K, V]) Iter() iter.Seq2[K, *V] {
	return func(yield func(K, *V) bool) {
		for i := 0; i < int(s.len); i++ {
			yield(K(s.idx2Key[i]), &s.data[i])
		}
	}
}

func (s *SparseArray[K, V]) IterReadOnly() iter.Seq2[K, *V] {
	return func(yield func(K, *V) bool) {
		for i := 0; i < int(s.len); i++ {
			cpy := s.data[i]
			yield(K(s.idx2Key[i]), &cpy)
		}
	}
}
