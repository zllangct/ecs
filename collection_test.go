package ecs

import (
	"fmt"
	"testing"
)

//待存储的数据定义
type Item struct {
	Component[Item, *Item]
	Count int
	Arr   [3]int
}

func TestCollectionIterator(t *testing.T) {
	//准备数据
	caseCount := 50
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Arr:   [3]int{1, 2, 3},
		})
	}

	//创建容器(无序数据集)
	c := NewCollection[Item]()

	//添加数据
	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		_, _ = c.Add(&srcList[i])
		cmp[int64(i)] = i
	}

	//遍历风格 1：
	for iter := NewIterator(c); !iter.End(); iter.Next() {
		v := iter.Val()
		fmt.Printf("style 1: %+v\n", v)
	}

	//遍历风格 2:
	iter := NewIterator(c)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		fmt.Printf("style 2: %+v\n", c)
	}
}

func BenchmarkSliceWrite(b *testing.B) {
	var slice []Item
	var id2index = map[int]int{}

	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
		}
		slice = append(slice, *item)
		id2index[n] = n
	}
}

func BenchmarkSliceRead(b *testing.B) {
	var slice []Item
	// collection 有ID生成，此处用通常方式模拟
	var id2index = map[int]int{}

	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
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

func BenchmarkSliceIter(b *testing.B) {
	var slice []Item
	// collection 有ID生成，此处用通常方式模拟
	var id2index = map[int]int{}

	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
		}
		slice = append(slice, *item)
		id2index[n] = n
		ids = append(ids, int64(n))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < 10000; i++ {
			_ = slice[i]
		}
	}
}

func BenchmarkCollectionWrite(b *testing.B) {
	c := NewCollection[Item]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
		}
		ret, _ := c.Add(item)
		_ = ret
	}
}

func BenchmarkCollectionRead(b *testing.B) {
	c := NewCollection[Item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
		}
		_, _ = c.Add(item)
		ids = append(ids, int64(n+1))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = c.Get(ids[(n+1)%total])
	}
}

func BenchmarkCollectionIter(b *testing.B) {
	c := NewCollection[Item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
		}
		_, _ = c.Add(item)
		ids = append(ids, int64(n))
	}

	iter := NewIterator(c)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for item := iter.Begin(); !iter.End(); item = iter.Next() {
			_ = item
		}
	}
}
