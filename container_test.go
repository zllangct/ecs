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
