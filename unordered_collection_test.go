package ecs

import (
	"testing"
	"unsafe"
)

type __unorderedCollection_Test_item struct {
	Component[__unorderedCollection_Test_item]
	ItemID int64
	Arr    [3]int
}

func TestUnorderedCollection_RW(t *testing.T) {
	c := NewUnorderedCollection[__unorderedCollection_Test_item](0)
	for i := 0; i < 5; i++ {
		_, _ = c.Add(&__unorderedCollection_Test_item{
			ItemID: int64(i),
			Arr:    [3]int{1, 2, 3},
		})
	}

	if get, want := c.Get(2).ItemID, int64(2); get != want {
		t.Errorf("want: %d, get: %d", want, get)
	}

	for i := 0; i < 5; i++ {
		c.Remove(0)
	}
}

func TestUnorderedCollection_Iterator(t *testing.T) {
	//准备数据
	caseCount := 50
	var srcList []__unorderedCollection_Test_item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, __unorderedCollection_Test_item{
			ItemID: int64(i),
			Arr:    [3]int{1, 2, 3},
		})
	}

	//创建容器(无序数据集)
	c := NewUnorderedCollection[__unorderedCollection_Test_item]()

	//添加数据
	for i := 0; i < caseCount; i++ {
		_, _ = c.Add(&srcList[i])
	}

	//遍历风格 1：
	for iter := NewUnorderedCollectionIterator(c); !iter.End(); iter.Next() {
		v := iter.Val()
		_ = v
		//t.Logf("style 1: %+v\n", v)
	}

	//遍历风格 2:
	iter := NewUnorderedCollectionIterator(c)
	for c := iter.Begin(); !iter.End(); c = iter.Next() {
		_ = c
		//t.Logf("style 2: %+v\n", c)
	}
}

func BenchmarkUnorderedCollection_SliceIter(b *testing.B) {
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
		for i := 0; i < 10000; i++ {
			_ = slice[i]
		}
	}
}

func BenchmarkUnorderedCollection_Write(b *testing.B) {
	c := NewUnorderedCollection[__unorderedCollection_Test_item]()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		ret, _ := c.Add(item)
		_ = ret
	}
}

func BenchmarkUnorderedCollection_SliceRead(b *testing.B) {
	total := 100000
	c := NewUnorderedCollection[__unorderedCollection_Test_item]()
	arr := make([]__unorderedCollection_Test_item, total)
	for n := 0; n < total; n++ {
		item := __unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		_, _ = c.Add(&item)
		arr = append(arr, item)
	}

	fn := func(idx int64) *__unorderedCollection_Test_item {
		return nil
	}

	b.ResetTimer()

	b.Run("slice", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = arr[n%total]
		}
	})
	b.Run("unordered collection", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = c.Get(int64(n % total))
		}
	})
	b.Run("unordered collection 2", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = (*__unorderedCollection_Test_item)(unsafe.Add(unsafe.Pointer(&c.data[0]), uintptr(n%total)*c.eleSize))
		}
	})
	b.Run("unordered collection empty func", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = fn(int64(n % total))
		}
	})
}

func BenchmarkUnorderedCollection_Read(b *testing.B) {
	c := NewUnorderedCollection[__unorderedCollection_Test_item]()
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
		_ = c.Get(int64((n + 1) % total))
	}
}

func BenchmarkUnorderedCollection_ReadUnsafe(b *testing.B) {
	var ids []int64
	total := 100000
	data := make([]__unorderedCollection_Test_item, total)
	for n := 0; n < total; n++ {
		item := __unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		data = append(data, item)
		ids = append(ids, int64(n+1))
	}
	eleSize := unsafe.Sizeof(__unorderedCollection_Test_item{})

	b.ResetTimer()
	b.Run("direct", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = data[int64((n+1)%total)]
		}
	})

	b.Run("unsafe", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = (*__unorderedCollection_Test_item)(unsafe.Add(unsafe.Pointer(&data[0]), uintptr(int64((n+1)%total))*eleSize))
		}
	})
}
