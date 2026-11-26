package ecs

import (
	"slices"
	"testing"
)

func TestQueryPairResultAndReuse(t *testing.T) {
	for _, shape := range []struct {
		name string
		a    func(int) bool
		b    func(int) bool
	}{
		{"covered-by-dense-range", func(i int) bool { return i >= 10 && i < 60 && i%3 == 0 }, func(i int) bool { return i >= 10 && i < 60 }},
		{"partly-inside-dense-range", func(i int) bool { return i%17 == 0 }, func(i int) bool { return i >= 10 && i < 60 }},
		{"merge-with-misses", func(i int) bool { return i%4 != 0 }, func(i int) bool { return i%3 != 0 }},
		{"probe-with-misses", func(i int) bool { return i == 1 || i == 32 || i == 255 }, func(i int) bool { return i%2 != 0 }},
		{"disjoint", func(i int) bool { return i > 128 }, func(i int) bool { return i < 128 }},
		{"missing-pool", func(i int) bool { return true }, func(i int) bool { return false }},
	} {
		t.Run(shape.name, func(t *testing.T) {
			world := NewWorld(WithoutGroups())
			var want []EntityIndex
			for i := range 256 {
				var components []Component
				if shape.a(i) {
					components = append(components, newGroupCompA(1))
				}
				if shape.b(i) {
					components = append(components, newGroupCompB(1))
				}
				e := world.NewEntity(WithComponents(components...))
				if shape.a(i) && shape.b(i) {
					want = append(want, e.Index())
				}
			}
			called := false
			err := world.RegisterLight(func(ctx *SystemContext, _ Event) error {
				called = true
				q := ctx.NewQuery(WithComp[groupCompA](), WithComp[groupCompB](), WithComp[groupCompA]())
				rows := q.Iter()
				for range rows {
					break
				}
				for range 2 {
					var got []EntityIndex
					for key, info := range rows {
						if info == nil || info.Entity().Index() != key {
							t.Fatalf("incorrect entity for key %v", key)
						}
						if _, writable := QueryGet[groupCompB](&q, key); writable {
							t.Fatal("read-only component has a writable slot")
						}
						got = append(got, key)
					}
					if !slices.Equal(got, want) {
						t.Fatalf("got %v, want %v", got, want)
					}
				}
				return nil
			}, WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly)))
			if err != nil {
				t.Fatal(err)
			}
			if err := world.Update(); err != nil {
				t.Fatal(err)
			}
			if !called {
				t.Fatal("system did not run")
			}
		})
	}
}
