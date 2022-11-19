package ecs

import (
	"math/rand"
	"testing"
)

func TestComponentSet_Sort(t *testing.T) {
	//准备数据
	caseCount := 50
	var srcList []__unorderedCollection_Test_item
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, __unorderedCollection_Test_item{
			Component: Component[__unorderedCollection_Test_item]{
				seq:   uint32(caseCount - i),
				owner: Entity(i),
			},
			ItemID: int64(i),
			Arr:    [3]int{1, 2, 3},
		})
	}

	//创建容器(无序数据集)
	c := NewComponentSet[__unorderedCollection_Test_item](&ComponentMetaInfo{})

	//添加数据
	for i := 0; i < caseCount; i++ {
		_ = c.Add(&srcList[i], srcList[i].Owner())
	}

	i := 0
	c.Range(func(item IComponent) bool {
		if item.(*__unorderedCollection_Test_item).seq != uint32(caseCount-i) {
			t.Errorf("sort error, want %d, got %d", caseCount, item.(*__unorderedCollection_Test_item).seq)
			return false
		}
		i++
		return true
	})

	//排序
	c.Sort()

	//验证
	i = 1
	c.Range(func(item IComponent) bool {
		if item.(*__unorderedCollection_Test_item).seq != uint32(i) {
			t.Errorf("sort error, want %d, got %d", i, item.(*__unorderedCollection_Test_item).seq)
			return false
		}
		i++
		return true
	})
}

func TestNewComponentSet(t *testing.T) {
	cs := NewComponentSet[__unorderedCollection_Test_item](&ComponentMetaInfo{})
	if cs.GetElementMeta().it != 0 {
		t.Error("element meta error")
	}
}

func BenchmarkComponentSet_Read(b *testing.B) {
	c := NewComponentSet[__unorderedCollection_Test_item](&ComponentMetaInfo{})
	var ids []int64
	total := 1000000
	for n := 0; n < total; n++ {
		item := &__unorderedCollection_Test_item{
			ItemID: int64(n),
		}
		_ = c.Add(item, Entity(n))
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
			_ = c.Get(Entity(seq[n%total]))
		}
	})
	b.Run("random", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = c.Get(Entity(r[n%total]))
		}
	})
}
