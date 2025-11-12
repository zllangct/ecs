package ecs

import (
	"fmt"
	"math/rand/v2"
	"sort"
	"testing"
)

func TestSparseArray_Iter(t *testing.T) {
	type V struct {
		Seq int
	}

	sa := NewSparseArray[int, V]()

	for i := range 10 {
		sa.Add(i, &V{i})
	}

	removed := sa.Remove(5)
	if removed.Seq != 5 {
		t.Errorf("want %v, got %v", 5, removed.Seq)
	}

	got := sa.Get(5)
	if got != nil {
		t.Errorf("want nil, got %v", got)
	}

	got2 := sa.Get(9)
	if got2.Seq != 9 {
		t.Errorf("want %v, got %v", 9, got2.Seq)
	}

	if sa.idx2Key[5] != 9 {
		t.Errorf("want %v, got %v", 9, sa.idx2Key[4])
	}

	sa.Add(15, &V{15})

	if sa.Get(15).Seq != 15 {
		t.Errorf("want %v, got %v", 15, sa.Get(15).Seq)
	}

	if sa.idx2Key[9] != 15 {
		t.Errorf("want %v, got %v", 9, sa.idx2Key[15])
	}

	w := []int{0, 1, 2, 3, 4, 9, 6, 7, 8, 15}

	index := 0
	for i, v := range sa.Iter() {
		if i != w[index] {
			t.Errorf("want %v, got %v", w[index], i)
		}
		if v.Seq != w[index] {
			t.Errorf("want %v, got %v", w[index], v.Seq)
		}
		index++
	}
}

func TestSparseArray_Sort(t *testing.T) {
	type V struct {
		Seq int
	}

	sa := NewSparseArray[int, V]()
	ids := []int{7, 4, 6, 5, 1}
	for _, i := range ids {
		sa.Add(i, &V{i})
	}

	for i, v := range sa.Iter() {
		fmt.Printf("(%v %v) ", i, v)
	}
	fmt.Printf("\n")

	sa.Sort()

	for i, v := range sa.Iter() {
		fmt.Printf("(%v %v) ", i, v)
	}
	fmt.Printf("\n")
}

func TestSparseArray_Sort2(t *testing.T) {
	type V struct {
		Seq int
	}

	sa := NewSparseArray[int, V]()
	ids := []int{7, 4, 6, 5, 1}
	for _, i := range ids {
		sa.Add(i, &V{i})
	}

	for i, v := range sa.Iter() {
		fmt.Printf("(%v %v) ", i, v)
	}
	fmt.Printf("\n")

	sort.Sort(sa)

	for i, v := range sa.Iter() {
		fmt.Printf("(%v %v) ", i, v)
	}
	fmt.Printf("\n")
}

func BenchmarkSparseArray_Sort(b *testing.B) {
	maxNum := 10000
	var ids []int = make([]int, maxNum)
	for i := range maxNum {
		ids[i] = i
	}
	rand.Shuffle(10000, func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})

	type V struct {
		Seq int
	}

	b.Run("sort", func(b *testing.B) {
		sa := NewSparseArray[int, V]()
		for _, i := range ids {
			sa.Add(i, &V{i})
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sa.Sort()
		}
	})
	b.Run("sort2", func(b *testing.B) {
		sa := NewSparseArray[int, V]()
		for _, i := range ids {
			sa.Add(i, &V{i})
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sort.Sort(sa)
		}
	})
}
