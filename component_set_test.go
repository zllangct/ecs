package ecs

import (
	"math/rand"
	"testing"
)

func TestComponentSet_Sort(t *testing.T) {
	//准备数据
	caseCount := 50
	var srcList []item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, item{
			Component: Component[item]{
				seq: uint32(caseCount - i),
			},
			ID:  int64(i),
			Arr: [3]int{1, 2, 3},
		})
	}

	//创建容器(无序数据集)
	c := NewComponentSet[item]()

	//添加数据
	for i := 0; i < caseCount; i++ {
		_, _ = c.Add(&srcList[i], srcList[i].ID)
	}

	//排序
	c.Sort()

	//验证
	for i := 0; i < caseCount; i++ {
		if c.Get(int64(i)).ID != int64(caseCount-i-1) {
			t.Error("sort error")
			break
		}
	}
}

func TestNewComponentSet(t *testing.T) {
	cs := NewComponentSet[item]()
	if cs.GetElementMeta().it == 0 {
		t.Error("element meta error")
	}
}

func BenchmarkComponentSet_Read(b *testing.B) {
	c := NewComponentSet[item]()
	var ids []int64
	total := 1000000
	for n := 0; n < total; n++ {
		item := &item{
			ID: int64(n),
		}
		_, _ = c.Add(item)
		ids = append(ids, int64(n+1))
	}

	seq := make([]int, total)
	r := make([]int, total)

	for i := 0; i < total; i++ {
		seq[i] = i
		r[i] = i
	}
	rand.Shuffle(len(r), func(i, j int) {
		r[i], r[j] = r[j], r[i]
	})

	b.ResetTimer()

	b.Run("sequence", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = c.Get(int64(seq[n%total]))
		}
	})
	b.Run("random", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = c.Get(int64(r[n%total]))
		}
	})
}
