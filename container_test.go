package ecs

import (
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestNewContainer(t *testing.T) {
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
		typ := reflect.TypeOf(Item{})
		c := NewContainer(typ.Size())

		cmp := map[int]int{}
		for i := 0; i < caseCount; i++ {
			idx, _ := c.Add(unsafe.Pointer(&srcList[i]))
			cmp[idx] = i
		}

		for idx, isrc := range cmp {
			p := c.Get(idx)
			item := *((*Item)(p))
			if item != srcList[isrc] {
				t.Errorf("src item %v != container item %v", srcList[isrc], item)
			}
		}
	})
}

func BenchmarkContainerNormalWrite(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}

	typ := reflect.TypeOf(Item{})
	c := NewContainer(typ.Size())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}

		c.Add(unsafe.Pointer(item))
	}
}

func BenchmarkContainerNormalRead(b *testing.B) {
	type Item struct {
		Count int
		Name  string
	}

	typ := reflect.TypeOf(Item{})
	c := NewContainer(typ.Size())

	for n := 0; n < b.N; n++ {
		item := &Item{
			Count: n,
			Name:  "foo" + strconv.Itoa(n),
		}

		c.Add(unsafe.Pointer(item))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		p := c.Get(n)
		item := *((*Item)(p))
		_ = item
	}
}

// func BenchmarkContainerGenericWrite(b *testing.B) {
// 	type Item struct {
// 		Count int
// 		Name  string
// 	}
// 	c := NewTContainer[Item]()
// 	b.ResetTimer()
// 	for n := 0; n < b.N; n++ {
// 		item := Item{
// 			Count: n,
// 			Name:  "foo" + strconv.Itoa(n),
// 		}

// 		c.Add(item)
// 	}
// }
