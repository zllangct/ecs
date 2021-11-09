package ecs

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCollectionIterator(t *testing.T) {
	//待存储的数据定义
	type Item struct {
		Count int
		Name  string
		Arr   []int
	}

	//准备数据
	caseCount := 50
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
			Arr:   []int{1, 2, 3},
		})
	}

	//创建容器(无序数据集)
	c := NewCollection[Item]()

	//添加数据
	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp[id] = i
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

func BenchmarkCollectionWrite(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}
	c := NewCollection[Item]()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}

		id, ret := c.Add(item)
		_, _ = id, ret
	}
}

func BenchmarkCollectionRead(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}
	c := NewCollection[Item]()
	var ids []int64
	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}
		id, _ := c.Add(item)
		ids = append(ids, id)
	}
	b.ResetTimer()
	for n := 0; n < len(ids); n++ {
		_ = c.Get(ids[n])
	}
}
