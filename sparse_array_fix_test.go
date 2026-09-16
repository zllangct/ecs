package ecs

import (
	"math/rand/v2"
	"testing"
)

type saTestVal struct {
	Seq int
}

// D03: Remove 后 idx2Key 必须截断，不得残留已删除 key
func TestSparseArrayRemove_Idx2KeyTruncated(t *testing.T) {
	sa := NewSparseArray[int, saTestVal]()
	for i := 0; i < 10; i++ {
		sa.Add(i, &saTestVal{i})
	}
	sa.Remove(5)
	if len(sa.idx2Key) != sa.Len() {
		t.Errorf("idx2Key not truncated: len(idx2Key)=%d, Len()=%d", len(sa.idx2Key), sa.Len())
	}
	for _, k := range sa.idx2Key {
		if int(k) == 5 {
			t.Errorf("stale key 5 remains in idx2Key: %v", sa.idx2Key)
		}
	}
}

// D03: swap-remove 后所有存活元素的 Get 必须正确
func TestSparseArrayRemove_SwapRemoveConsistency(t *testing.T) {
	sa := NewSparseArray[int, saTestVal]()
	for i := 0; i < 10; i++ {
		sa.Add(i, &saTestVal{i * 100})
	}
	sa.Remove(3)
	sa.Remove(7)
	sa.Remove(9) // 连续删除，含多次 swap

	for i := 0; i < 10; i++ {
		got := sa.Get(i)
		if i == 3 || i == 7 || i == 9 {
			if got != nil {
				t.Errorf("key %d: want nil, got %v", i, got)
			}
			continue
		}
		if got == nil || got.Seq != i*100 {
			t.Errorf("key %d: want %d, got %v", i, i*100, got)
		}
	}
	// 迭代只能产出一个元素一次，且无已删除 key
	seen := map[int]int{}
	for k, v := range sa.Iter() {
		seen[k]++
		if v.Seq != k*100 {
			t.Errorf("key %d: value mismatch %d", k, v.Seq)
		}
	}
	if len(seen) != 7 {
		t.Errorf("want 7 elements iterated, got %d", len(seen))
	}
}

// D05: shrink 抬高 maxKey 后 Get 越界 panic
func TestSparseArrayGet_BoundaryAfterShrink(t *testing.T) {
	sa := NewSparseArray[int, saTestVal]()
	for i := 0; i <= 100; i++ {
		sa.Add(i, &saTestVal{i})
	}
	sa.Remove(100) // 触发 shrink，maxKey 被抬到 1024，超过 indices 实际长度
	// 不应 panic，应返回 nil
	if got := sa.Get(500); got != nil {
		t.Errorf("want nil, got %v", got)
	}
	if sa.Exist(500) {
		t.Errorf("key 500 should not exist")
	}
}

// D04: Sort 后所有 key 的 Get 必须仍返回原值，且迭代按 key 升序
func TestSparseArraySort_GetConsistencyAfterSort(t *testing.T) {
	sa := NewSparseArray[int, saTestVal]()
	ids := []int{7, 4, 6, 5, 1, 100, 42, 3}
	for _, i := range ids {
		sa.Add(i, &saTestVal{i * 10})
	}
	sa.Sort()
	// 迭代必须升序
	prev := -1
	count := 0
	for k, v := range sa.Iter() {
		if k <= prev {
			t.Errorf("not sorted: %d after %d", k, prev)
		}
		if v.Seq != k*10 {
			t.Errorf("iter key %d: value mismatch %d", k, v.Seq)
		}
		prev = k
		count++
	}
	if count != len(ids) {
		t.Errorf("want %d, got %d", len(ids), count)
	}
	// Sort 后随机访问必须正确
	for _, i := range ids {
		got := sa.Get(i)
		if got == nil || got.Seq != i*10 {
			t.Errorf("Get(%d) after sort: want %d, got %v", i, i*10, got)
		}
	}
}

// D04: 随机排列下 Sort 的正确性（排序+Get 一致性）
func TestSparseArraySort_RandomPermutation(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 0))
	for round := 0; round < 200; round++ {
		n := 1 + rng.IntN(50)
		ids := rng.Perm(n * 3)[:n] // 稀疏 key，含空洞
		sa := NewSparseArray[int, saTestVal]()
		for _, i := range ids {
			sa.Add(i, &saTestVal{i})
		}
		sa.Sort()
		prev := -1
		for k, v := range sa.Iter() {
			if k <= prev {
				t.Fatalf("round %d: not sorted: %d after %d (ids=%v)", round, k, prev, ids)
			}
			if v.Seq != k {
				t.Fatalf("round %d: iter key %d value %d", round, k, v.Seq)
			}
			prev = k
		}
		for _, i := range ids {
			got := sa.Get(i)
			if got == nil || got.Seq != i {
				t.Fatalf("round %d: Get(%d) wrong: %v (ids=%v)", round, i, got, ids)
			}
		}
	}
}

// D10: 迭代器必须响应 break
func TestSparseArrayIter_EarlyBreak(t *testing.T) {
	sa := NewSparseArray[int, saTestVal]()
	for i := 0; i < 5; i++ {
		sa.Add(i, &saTestVal{i})
	}
	count := 0
	for range sa.Iter() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("want 1, got %d", count)
	}
}
