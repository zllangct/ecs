package ecs

import (
	"testing"
)

func BenchmarkUnorderedCollectionWithID_Read(b *testing.B) {
	c := NewUnorderedCollectionWithID[__unorderedCollection_Test_item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		_, _ = c.Add(item)
		ids = append(ids, int64(n+1))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = c.GetByID(ids[(n+1)%total])
	}
}

func BenchmarkUnorderedCollectionWithID_Iter(b *testing.B) {
	c := NewUnorderedCollectionWithID[__unorderedCollection_Test_item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		_, _ = c.Add(item)
		ids = append(ids, int64(n))
	}

	iter := NewUnorderedCollectionWithIDIterator(c)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for item := iter.Begin(); !iter.End(); item = iter.Next() {
			_ = item
		}
	}
}

func BenchmarkUnorderedCollectionWithID_SliceWrite(b *testing.B) {
	var slice []__unorderedCollection_Test_item
	var id2index = map[int]int{}

	for n := 0; n < b.N; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		slice = append(slice, *item)
		id2index[n] = n
	}
}

func BenchmarkUnorderedCollectionWithID_SliceRead(b *testing.B) {
	var slice []__unorderedCollection_Test_item
	// collection 有ID生成，此处用通常方式模拟
	var id2index = map[int]int{}

	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		slice = append(slice, *item)
		id2index[n] = n
		ids = append(ids, int64(n))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = slice[id2index[n%total]]
	}
}
