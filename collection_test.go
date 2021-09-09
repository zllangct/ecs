package ecs

import (
	"strconv"
	"testing"
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

	t.Run("test1", func(t *testing.T) {
		c := NewCollection[Item]()

		cmp := map[int]int{}
		for i := 0; i < caseCount; i++ {
			idx, _ := c.Add(srcList[i])
			cmp[idx] = i
		}

		for idx, isrc := range cmp {
			p := c.Get(idx)
			item := *p
			if item != srcList[isrc] {
				t.Errorf("src item %v != container item %v", srcList[isrc], item)
			}
		}
	})
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

	c := NewCollection[Item]()

	cmp := map[int]int{}
	for i := 0; i < caseCount; i++ {
		idx, _ := c.Add(srcList[i])
		cmp[idx] = i
	}

	iter := NewIterator(c)
	for v := iter.Next(); v != nil; v = iter.Next(){
		_ = v
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
		item := Item{
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
	for n := 0; n < b.N; n++ {
		item := Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}
		c.Add(item)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = c.Get(n)
	}
}
