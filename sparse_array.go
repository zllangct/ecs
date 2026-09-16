package ecs

import (
	"iter"
	"sort"
)

type SparseArray[K Integer, V any] struct {
	USet[V]
	indices         []int32
	idx2Key         []int32
	maxKey          K
	shrinkThreshold int32
	initSize        int
	isKOrder        bool
	// lastKey 最近一次 Add 的 key：递增添加不破坏 isKOrder（保序优化，
	// 使 merge-join 在"顺序创建实体"的常规负载下无需 Sort 即可生效）
	lastKey K
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

	// 保序：key 递增则有序性保持，否则标记失序
	if key < s.lastKey {
		s.isKOrder = false
	} else {
		s.lastKey = key
	}

	return &s.data[idx]
}

func (s *SparseArray[K, V]) Remove(key K) *V {
	if key < 0 || key >= K(len(s.indices)) {
		return nil
	}
	idx := s.indices[key] - 1
	if idx < 0 {
		return nil
	}

	// swap-remove：用末尾元素覆盖被删位置，同步维护 indices/idx2Key
	lastIdx := int32(s.len) - 1
	lastKey := s.idx2Key[lastIdx]
	removed := s.data[idx]
	s.data[idx] = s.data[lastIdx]
	s.len--

	s.idx2Key[idx] = lastKey
	s.indices[lastKey] = idx + 1
	s.indices[key] = 0
	s.idx2Key = s.idx2Key[:lastIdx]

	s.shrink(key)

	s.isKOrder = false

	return &removed
}

func (s *SparseArray[K, V]) Exist(key K) bool {
	if key < 0 || key >= K(len(s.indices)) {
		return false
	}
	return !(s.indices[key] == 0)
}

func (s *SparseArray[K, V]) Get(key K) *V {
	if key < 0 || key >= K(len(s.indices)) {
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
	s.lastKey = 0
}

func (s *SparseArray[K, V]) Less(i, j int) bool {
	return s.idx2Key[i] < s.idx2Key[j]
}

func (s *SparseArray[K, V]) Swap(i, j int) {
	s.USet.Swap(int64(i), int64(j))
	s.idx2Key[i], s.idx2Key[j] = s.idx2Key[j], s.idx2Key[i]
	// indices 以 key 为下标，swap 后用 idx2Key 回写两个位置的映射
	s.indices[s.idx2Key[i]] = int32(i + 1)
	s.indices[s.idx2Key[j]] = int32(j + 1)
}

// Sort 将密集数组按 key 升序重排，使多个 SparseArray 联合遍历时内存访问连续。
// 实现：先求"按 key 升序的位置排列"，再用置换环 O(n) 应用到 data/idx2Key，
// 最后由 idx2Key 全量重建 indices 映射（存活 key 全部覆盖，删除位保持 0）。
func (s *SparseArray[K, V]) Sort() {
	if s.isKOrder {
		return
	}
	n := int(s.len)
	// sortedPos[r] = 第 r 小的 key 当前所在位置
	sortedPos := make([]int, n)
	for i := range sortedPos {
		sortedPos[i] = i
	}
	sort.Slice(sortedPos, func(a, b int) bool { return s.idx2Key[sortedPos[a]] < s.idx2Key[sortedPos[b]] })
	// dst[p] = 位置 p 的元素应去的目标位置（其 key 的升序名次）
	dst := make([]int, n)
	for r, p := range sortedPos {
		dst[p] = r
	}
	// 置换环应用
	for i := 0; i < n; i++ {
		for dst[i] != i {
			j := dst[i]
			s.data[i], s.data[j] = s.data[j], s.data[i]
			s.idx2Key[i], s.idx2Key[j] = s.idx2Key[j], s.idx2Key[i]
			dst[i], dst[j] = dst[j], dst[i]
		}
	}
	// 重建稀疏索引
	for pos := 0; pos < n; pos++ {
		s.indices[s.idx2Key[pos]] = int32(pos + 1)
	}
	s.isKOrder = true
	if n > 0 {
		s.lastKey = K(s.idx2Key[n-1])
	}
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

	if len(s.indices) > 1024 && int(s.maxKey) < len(s.indices)/2 {
		m := (s.maxKey + 1) * 5 / 4
		newIndices := make([]int32, m)
		copy(newIndices, s.indices[:m])
		s.indices = newIndices
	}
}

func (s *SparseArray[K, V]) Iter() iter.Seq2[K, *V] {
	return func(yield func(K, *V) bool) {
		for i := 0; i < int(s.len); i++ {
			if !yield(K(s.idx2Key[i]), &s.data[i]) {
				return
			}
		}
	}
}
