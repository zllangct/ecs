package ecs

import (
	"strconv"
	"testing"
)

func TestNewCollection(t *testing.T) {
	type Item struct {
		Count int
		Name  string
		o1    int
	}
	caseCount := 10
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
		})
	}
	c := NewCollection[Item]()

	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp[id] = i
	}

	for id, idx := range cmp {
		item := c.Get(id)
		if *item != srcList[idx] {
			t.Errorf("src item %v != container item %v", srcList[idx], item)
		}
	}
}

func TestCollectionIterator(t *testing.T) {
	type Item struct {
		Count int
		Name  string
	}
	caseCount := 100
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
		})
	}

	c := NewCollection[Item]()

	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(&srcList[i])
		cmp[id] = i
	}

	//for iter := NewIterator(c) ; !iter.End(); iter.Next(){
	//	v := iter.Val()
	//	fmt.Printf("%+v", v)
	//}
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

		c.Add(item)
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
