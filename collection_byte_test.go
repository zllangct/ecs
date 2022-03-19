package ecs

import (
	"fmt"
	"testing"
)

func TestCollectionByte(t *testing.T) {
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
	c := NewCollectionByte[Item]()

	//添加数据
	var cmp []int64
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp = append(cmp, id)
	}

	//查询数据
	item := c.Get(cmp[5])
	Log.Infof("item get : %+v", item)

	//迭代器
	iter := NewIteratorByte(c)
	for item := iter.Begin(); !iter.End(); item = iter.Next() {
		fmt.Printf("style 2: %+v\n", item)
	}
}

func BenchmarkCollectionByteWrite(b *testing.B) {
	c := NewCollectionByte[Item]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
		}
		id, ret := c.Add(item)
		_, _ = id, ret
	}
}

func BenchmarkCollectionByteRead(b *testing.B) {
	c := NewCollectionByte[Item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
		}
		id, _ := c.Add(item)
		ids = append(ids, id)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		_ = c.Get(ids[n%total])
	}
}

func BenchmarkCollectionByteIter(b *testing.B) {
	c := NewCollectionByte[Item]()
	var ids []int64
	total := 100000
	for n := 0; n < total; n++ {
		item := &Item{
			Count: n,
		}
		id, _ := c.Add(item)
		ids = append(ids, id)
	}

	iter := NewIteratorByte(c)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		v := iter.Val()
		_ = v
		iter.Next()
		if iter.End() {
			iter.Begin()
		}
	}
}
