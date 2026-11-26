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
	// Last peak that actually triggered a shrink. A repeated peak can reuse a bounded backing array.
	lastShrunkKey K
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
	if key < 0 {
		return nil
	}
	// already existed
	if uint64(key) < uint64(len(s.indices)) && s.indices[key] != 0 {
		return nil
	}
	if uint64(key) >= uint64(len(s.indices)) {
		s.growIndices(int(key))
	}
	_, idx := s.USet.Add(value)
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

func (s *SparseArray[K, V]) growIndices(key int) {
	n := key + 1
	if len(s.indices) > 0 {
		if len(s.indices) < int(s.shrinkThreshold) {
			n = max(n, key*2)
		} else {
			n = max(n, key+key/4)
		}
	}
	if n <= cap(s.indices) {
		// Every slot above maxKey was zero when the slice was shortened.
		s.indices = s.indices[:n]
		return
	}
	indices := make([]int32, n)
	copy(indices, s.indices)
	s.indices = indices
}

func (s *SparseArray[K, V]) Remove(key K) *V {
	idx := s.denseIndex(key)
	if idx < 0 {
		return nil
	}
	removed := s.data[idx]
	s.removeAt(key, idx)
	return &removed
}

func (s *SparseArray[K, V]) denseIndex(key K) int32 {
	if uint64(key) >= uint64(len(s.indices)) {
		return -1
	}
	return s.indices[key] - 1
}

func (s *SparseArray[K, V]) removeAt(key K, idx int32) {
	// swap-remove：用末尾元素覆盖被删位置，同步维护 indices/idx2Key
	lastIdx := int32(s.len) - 1
	if idx != lastIdx {
		movedKey := s.idx2Key[lastIdx]
		s.data[idx] = s.data[lastIdx]
		s.idx2Key[idx] = movedKey
		s.indices[movedKey] = idx + 1
		s.isKOrder = false
	} else if lastIdx > 0 {
		s.lastKey = K(s.idx2Key[lastIdx-1])
	} else {
		s.isKOrder = true
		s.lastKey = 0
	}
	// 清零空出的尾槽，避免滞留指针/GC 引用
	var zero V
	s.data[lastIdx] = zero
	s.len--

	s.indices[key] = 0
	s.idx2Key = s.idx2Key[:lastIdx]
	s.shrink(key)
}

func (s *SparseArray[K, V]) Exist(key K) bool {
	if uint64(key) >= uint64(len(s.indices)) {
		return false
	}
	return !(s.indices[key] == 0)
}

func (s *SparseArray[K, V]) Get(key K) *V {
	if uint64(key) >= uint64(len(s.indices)) {
		return nil
	}
	idx := s.indices[key] - 1
	if idx < 0 {
		return nil
	}
	return s.USet.Get(int64(idx))
}

func (s *SparseArray[K, V]) Reset() {
	retained := s.lastShrunkKey != 0 && cap(s.indices) > 1024 && cap(s.indices) > len(s.indices)
	s.lastShrunkKey = 0
	if s.Len() == 0 {
		if retained {
			s.indices = make([]int32, len(s.indices))
		}
		return
	}
	s.USet.Reset()
	if int(s.maxKey) < 1024 {
		if retained {
			s.indices = make([]int32, len(s.indices))
		} else {
			clear(s.indices)
		}
	} else {
		s.indices = make([]int32, 0, s.initSize)
	}
	s.maxKey = 0
	s.idx2Key = []int32{}
	s.isKOrder = true
	s.lastKey = 0
}

// Free 释放底层数组与稀疏索引内存
func (s *SparseArray[K, V]) Free() {
	s.USet.Free()
	s.indices = nil
	s.idx2Key = nil
	s.maxKey = 0
	s.lastKey = 0
	s.lastShrunkKey = 0
	s.isKOrder = true
}

func (s *SparseArray[K, V]) Less(i, j int) bool {
	return s.idx2Key[i] < s.idx2Key[j]
}

func (s *SparseArray[K, V]) Swap(i, j int) {
	s.USet.Swap(int64(i), int64(j))
	s.idx2Key[i], s.idx2Key[j] = s.idx2Key[j], s.idx2Key[i]
	if i != j {
		s.isKOrder = false
	}
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

	if s.isKOrder {
		s.maxKey = s.lastKey
	} else {
		s.maxKey = s.maxAfterRemove(int(key))
	}

	if len(s.indices) > 1024 && int(s.maxKey) < len(s.indices)/2 {
		n := (int(s.maxKey) + 1) * 5 / 4
		// Release a one-off peak immediately. Repeated toggles may retain up to
		// 1 MiB per pool; shorten the logical index in either case for snapshots.
		const maxRetainedSlots = (1 << 20) / 4
		if key == s.lastShrunkKey && cap(s.indices) <= maxRetainedSlots {
			s.indices = s.indices[:n]
		} else {
			indices := make([]int32, n)
			copy(indices, s.indices[:n])
			s.indices = indices
		}
		s.lastShrunkKey = key
	}
}

func (s *SparseArray[K, V]) maxAfterRemove(key int) K {
	// Nearby dense keys need only a few probes. Bound probes across a sparse
	// gap by the cost of inspecting the remaining dense keys.
	stop := max(0, key-len(s.idx2Key))
	for i := key - 1; i >= stop; i-- {
		if s.indices[i] != 0 {
			return K(i)
		}
	}
	var highest int32
	for _, k := range s.idx2Key {
		highest = max(highest, k)
	}
	return K(highest)
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
