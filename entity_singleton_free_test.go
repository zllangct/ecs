package ecs

import (
	"math/rand"
	"slices"
	"testing"
)

// Allocate the complete ID space. Only every stride-th entity will be freed;
// the intervening IDs remain live, leaving holes between occupied free pages.
func singletonFreeInputs(n, stride int) (*EntityIDGenerator, []Entity) {
	g := NewEntityIDGenerator(n*stride+2, 1)
	ids := make([]Entity, 0, n)
	for i := 0; i < n*stride; i++ {
		id := g.NewID()
		if i%stride == 0 {
			ids = append(ids, id)
		}
	}
	return g, ids
}

func orderSingletonFrees(ids []Entity, order string) {
	switch order {
	case "Descending":
		slices.Reverse(ids)
	case "Shuffled", "ShuffledSpaced":
		rng := rand.New(rand.NewSource(915))
		rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
}

func TestEntitySingletonFreePages(t *testing.T) {
	for _, order := range []string{"Ascending", "Descending", "ShuffledSpaced"} {
		const n, stride = 512, 256
		g, ids := singletonFreeInputs(n, stride)
		orderSingletonFrees(ids, order)
		for _, id := range ids {
			g.FreeID(id)
		}
		if (len(g.freePages) > 0) != (order == "ShuffledSpaced") {
			t.Fatalf("%s built unexpected page index of size %d", order, len(g.freePages))
		}
		for i := range n {
			id := g.NewID().toReuseID()
			if id.index != EntityIndex(i*stride+1) || id.reuse != 1 {
				t.Fatalf("%s reused %+v at position %d", order, id, i)
			}
		}
	}
}

func BenchmarkEntitySingletonFree(b *testing.B) {
	for _, order := range []string{"Ascending", "Descending", "Shuffled", "ShuffledSpaced"} {
		b.Run(order, func(b *testing.B) {
			n, stride := 30_000, 1
			if order == "ShuffledSpaced" {
				n, stride = 4096, 256
			}
			b.StopTimer()
			b.ReportAllocs()
			for sample := 0; sample < b.N; sample++ {
				g, ids := singletonFreeInputs(n, stride)
				orderSingletonFrees(ids, order)
				b.StartTimer()
				for _, id := range ids {
					g.FreeID(id)
				}
				b.StopTimer()
				for i := range n {
					id := g.NewID().toReuseID()
					if id.index != EntityIndex(i*stride+1) || id.reuse != 1 {
						b.Fatalf("reuse at position %d: %+v", i, id)
					}
				}
			}
		})
	}
}
