package ecs

import "testing"

// The flag-only probe control keeps both pools and their key locality identical.
// Every sample goes through a registered system's real Query.Iter implementation.
func BenchmarkQueryPair(b *testing.B) {
	const size = 65536
	for _, shape := range []struct {
		name string
		hasA func(int) bool
		hasB func(int) bool
	}{
		{"Dense", func(i int) bool { return true }, func(i int) bool { return true }},
		{"Near", func(i int) bool { return i%8 != 0 }, func(i int) bool { return i%9 != 0 }},
		{"Spread4", func(i int) bool { return i%4 == 1 }, func(i int) bool { return i%9 != 0 }},
		{"Spread16", func(i int) bool { return i%16 == 1 }, func(i int) bool { return i%9 != 0 }},
		{"Spread64", func(i int) bool { return i%64 == 1 }, func(i int) bool { return i%9 != 0 }},
		{"Prefix64", func(i int) bool { return i < 64 }, func(i int) bool { return i%9 != 0 }},
		{"Tail64", func(i int) bool { return i >= size-64 }, func(i int) bool { return i%9 != 0 }},
		{"Tail1", func(i int) bool { return i == size-1 }, func(i int) bool { return i%9 != 0 }},
	} {
		for _, mode := range []string{"Ordered", "Probe"} {
			b.Run(shape.name+"/"+mode, func(b *testing.B) {
				world := NewWorld(WithoutGroups())
				var wantCount int
				var wantSum int64
				for i := range size {
					var components []Component
					a, bb := shape.hasA(i), shape.hasB(i)
					if a {
						components = append(components, newGroupCompA(1))
					}
					if bb {
						components = append(components, newGroupCompB(1))
					}
					e := world.NewEntity(WithComponents(components...))
					if a && bb {
						wantCount++
						wantSum += int64(e.Index())
					}
				}
				err := world.RegisterLight(func(ctx *SystemContext, _ Event) error {
					set, _ := world.getComponentSet(GetIntType[groupCompA]())
					if mode == "Probe" {
						set.(*CSet[groupCompA, *groupCompA]).isKOrder = false
					}
					query := ctx.NewQuery(WithComp[groupCompA](), WithComp[groupCompB]())
					b.ReportAllocs()
					for b.Loop() {
						count := 0
						var sum int64
						for key, info := range query.Iter() {
							if info.Entity().Index() != key {
								b.Fatal("incorrect entity key")
							}
							count++
							sum += int64(key)
						}
						if count != wantCount || sum != wantSum {
							b.Fatalf("intersection count/sum %d/%d, expected %d/%d", count, sum, wantCount, wantSum)
						}
					}
					return nil
				}, WithDeps(Dep[groupCompA](), Dep[groupCompB]()))
				if err != nil {
					b.Fatal(err)
				}
				if err := world.Update(); err != nil {
					b.Fatal(err)
				}
			})
		}
	}
}
