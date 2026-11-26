package ecs

import "testing"

// Measure allocation independently of FreeID. The second case releases two
// overlapping windows first so a generator with a page index will build it.
func BenchmarkEntityNewIDBatch(b *testing.B) {
	const n = 100_000
	for _, mode := range []string{"Fresh", "ReuseAfterScatteredFrees"} {
		b.Run(mode, func(b *testing.B) {
			b.StopTimer()
			b.ReportAllocs()
			ids := make([]Entity, n)
			for sample := 0; sample < b.N; sample++ {
				g := NewEntityIDGenerator(n+2, 10)
				var generation int32
				if mode == "ReuseAfterScatteredFrees" {
					generation = 1
					for i := range ids {
						ids[i] = g.NewID()
					}
					ids[0], ids[511] = ids[511], ids[0]
					ids[10], ids[767] = ids[767], ids[10]
					for _, id := range ids {
						g.FreeID(id)
					}
				}
				b.StartTimer()
				for i := range ids {
					ids[i] = g.NewID()
				}
				b.StopTimer()
				for i, id := range ids {
					if value := id.toReuseID(); value.index != EntityIndex(i+1) || value.reuse != generation {
						b.Fatalf("unexpected entity at position %d: %+v", i, value)
					}
				}
			}
		})
	}
}
