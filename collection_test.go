package ecs

import (
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestNewCollection(t *testing.T) {
	type Item struct {
		Count int
		Name  string
	}
	caseCount := 10
	var srcList []Item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, Item{
			Count: i,
			Name:  "foo" + strconv.Itoa(i),
		})
	}

	c := NewCollection(int(reflect.TypeOf(Item{}).Size()))

	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(unsafe.Pointer(&srcList[i]))
		cmp[id] = i
	}

	for id, idx := range cmp {
		p := c.Get(id)
		item := *(*Item)(p)
		if item != srcList[idx] {
			t.Errorf("src item %v != container item %v", srcList[idx], item)
		}
	}
}

func TestCollectionIterator(t *testing.T){
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

	c := NewCollection(int(reflect.TypeOf(Item{}).Size()))

	cmp := map[int64]int{}
	for i := 0; i < caseCount; i++ {
		id, _ := c.Add(unsafe.Pointer(&srcList[i]))
		cmp[id] = i
	}

	for iter := NewIterator[Item](c) ; !iter.End(); iter.Next(){
		_ = iter.Val
	}
}

func BenchmarkCollectionWrite(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}
	c := NewCollection(int(reflect.TypeOf(Item{}).Size()))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		item := Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}

		c.Add(unsafe.Pointer(&item))
	}
}

func BenchmarkCollectionRead(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}
	c := NewCollection(int(reflect.TypeOf(Item{}).Size()))
	for n := 0; n < b.N; n++ {
		item := Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}
		c.Add(unsafe.Pointer(&item))
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = c.get(n)
	}
}
