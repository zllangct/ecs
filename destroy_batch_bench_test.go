package ecs

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"
	"time"
)

func orderDestroyBenchmark(ids []Entity, order string) {
	switch order {
	case "Descending":
		slices.Reverse(ids)
	case "Shuffled":
		rng := rand.New(rand.NewSource(915))
		rng.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	}
}

func BenchmarkEntityFreeBatch(b *testing.B) {
	for _, n := range []int{1_000, 100_000} {
		for _, order := range []string{"Ascending", "Descending", "Shuffled"} {
			b.Run(fmt.Sprintf("%d/%s", n, order), func(b *testing.B) {
				b.StopTimer()
				b.ReportAllocs()
				for sample := 0; sample < b.N; sample++ {
					g := NewEntityIDGenerator(n+2, 10)
					ids := make([]Entity, n)
					for i := range ids {
						ids[i] = g.NewID()
					}
					orderDestroyBenchmark(ids, order)
					b.StartTimer()
					for _, id := range ids {
						g.FreeID(id)
					}
					b.StopTimer()
					for i := 1; i <= n; i++ {
						id := g.NewID().toReuseID()
						if id.index != EntityIndex(i) || id.reuse != 1 {
							b.Fatalf("ID reuse differs at %d: %+v", i, id)
						}
					}
				}
			})
		}
	}
}

// One operation enqueues and flushes the destruction of N existing entities.
// Setup, the initial component flush, and all postconditions are untimed.
func BenchmarkWorldDestroyBatch(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		for _, order := range []string{"Ascending", "Descending", "Shuffled"} {
			b.Run(fmt.Sprintf("%d/%s", n, order), func(b *testing.B) {
				benchmarkDestroyWorld(b, n, order, false)
			})
		}
	}
}

// A world with eight component pools and one component per destroyed entity
// exercises the cost of searching storage that cannot own the target.
func BenchmarkWorldDestroySparse(b *testing.B) {
	benchmarkDestroyWorld(b, 1_000, "Ascending", true)
}

func benchmarkDestroyWorld(b *testing.B, n int, order string, sparse bool) {
	b.StopTimer()
	b.ReportAllocs()
	for sample := 0; sample < b.N; sample++ {
		w := NewWorld()
		if !sparse {
			if err := w.RegisterLight(func(*SystemContext, Event) error { return nil },
				WithGroup(Dep[groupCompA](), Dep[groupCompB]())); err != nil {
				b.Fatal(err)
			}
		}
		targets := make([]Entity, n)
		for i := range targets {
			switch {
			case sparse || i%4 == 1:
				targets[i] = w.NewEntity(WithComponents(newGroupCompA(int32(i))))
			case i%4 == 0:
				targets[i] = w.NewEntity()
			case i%4 == 2:
				targets[i] = w.NewEntity(WithComponents(newGroupCompA(int32(i)), newGroupCompB(2)))
			default:
				targets[i] = w.NewEntity(WithComponents(newGroupCompA(int32(i)), newGroupCompB(2), &gCompC{V: 3}))
			}
		}
		components := []Component{newGroupCompA(41), newGroupCompB(42), &gCompC{V: 43}}
		if sparse {
			components = append(components, &gCompD{}, &gCompE{}, &gCompF{}, &gCompG{}, &gCompH{})
		}
		survivor := w.NewEntity(WithComponents(components...))
		if err := w.Step(time.Millisecond); err != nil {
			b.Fatal(err)
		}
		orderDestroyBenchmark(targets, order)

		b.StartTimer()
		for _, target := range targets {
			w.DestroyEntity(target)
		}
		err := w.Step(time.Millisecond)
		b.StopTimer()

		if err != nil {
			b.Fatal(err)
		}
		if w.entities.Len() != 1 {
			b.Fatalf("surviving entities %d, want 1", w.entities.Len())
		}
		for _, set := range w.components {
			for _, idx := range set.EntityIndexes() {
				if idx != survivor.Index() {
					b.Fatalf("unexpected residual entity %d", idx)
				}
			}
		}
		if a, ok := w.GetComponent[groupCompA](survivor); !ok || a.V != 41 {
			b.Fatal("lost survivor")
		}
		for idx := 1; idx <= n; idx++ {
			id := w.NewEntity().toReuseID()
			if id.index != EntityIndex(idx) || id.reuse != 1 {
				b.Fatalf("ID reuse differs at %d: %+v", idx, id)
			}
		}
	}
}
