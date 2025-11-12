package ecs

import (
	"math/rand"
	"testing"
)

func TestComponentSet(t *testing.T) {
	//prepare test data
	caseCount := 50
	var srcList []dummyComponent
	for i := 0; i < caseCount; i++ {
		srcList = append(srcList, dummyComponent{
			Seq: int32(i),
		})
	}

	//create component container
	c := NewCSet[dummyComponent]()

	//add test data
	for i := 0; i < caseCount; i++ {
		_ = c.Add(Entity(i), &srcList[i])
	}

	i := int32(0)
	for entity, comp := range c.Iter() {
		if comp.Seq != i {
			t.Errorf("error 1, want %d, got %d", i, comp.Seq)
		}
		if int32(entity) != i {
			t.Errorf("error 2, want %d, got %d", Entity(i), entity)
		}
		i++
	}
}

func BenchmarkComponentSet_Read(b *testing.B) {
	c := NewCSet[dummyComponent]()
	var ids []int64
	total := 1000000
	for n := 0; n < total; n++ {
		item := &dummyComponent{
			Seq: int32(n),
		}
		_ = c.Add(Entity(n), item)
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
