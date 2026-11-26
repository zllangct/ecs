package ecs

import (
	"iter"
	"testing"
)

func TestView2ROStorageAndAccess(t *testing.T) {
	for _, layout := range []string{"CSet", "SubsetWithResidual", "Mixed"} {
		t.Run(layout, func(t *testing.T) {
			w := NewWorld()
			opts := []SystemOption{WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly), Dep[p2CompC]())}
			switch layout {
			case "SubsetWithResidual":
				opts = append(opts, WithGroup(Dep[p2CompA](), Dep[p2CompB](), Dep[p2CompC]()))
			case "Mixed":
				opts = append(opts, WithGroup(Dep[p2CompA](), Dep[p2CompC]()))
			}
			full := w.NewEntity(WithComponents(&p2CompA{1}, &p2CompB{10}, &p2CompC{}))
			residual := w.NewEntity(WithComponents(&p2CompA{2}, &p2CompB{20}))
			unmatched := w.NewEntity(WithComponents(&p2CompA{3}))
			if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
				rows := View2RO[p2CompA, p2CompBView](ctx)
				first := 0
				for range rows {
					first++
					break
				}
				if first != 1 {
					t.Fatalf("early break visited %d rows", first)
				}
				for a, bv := range rows {
					a.V += bv.V()
				}
				return nil
			}, opts...); err != nil {
				t.Fatal(err)
			}
			for _, rejected := range [][]SystemOption{
				{WithDeps(Dep[p2CompA](ReadOnly), Dep[p2CompB](ReadOnly))},
				{WithDeps(Dep[p2CompA](), Dep[p2CompB]())},
				{WithDeps(Dep[p2CompA]())},
			} {
				if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
					for range View2RO[p2CompA, p2CompBView](ctx) {
						t.Fatal("view yielded without the required read/write permissions")
					}
					return nil
				}, rejected...); err != nil {
					t.Fatal(err)
				}
			}
			if err := w.Update(); err != nil {
				t.Fatal(err)
			}
			for entity, want := range map[Entity]int32{full: 11, residual: 22, unmatched: 3} {
				if a, ok := w.GetComponent[p2CompA](entity); !ok || a.V != want {
					t.Fatalf("entity %v: got %v, want %d", entity, a, want)
				}
			}
		})
	}
}

func TestView2ROCSetSuspendedIteratorRespectsContext(t *testing.T) {
	w := NewWorld(WithoutGroups())
	for range 2 {
		w.NewEntity(WithComponents(&p2CompA{}, &p2CompB{}))
	}
	var next func() (*p2CompA, p2CompBView, bool)
	var stop func()
	if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		next, stop = iter.Pull2(View2RO[p2CompA, p2CompBView](ctx))
		if _, _, ok := next(); !ok {
			t.Fatal("active system could not read a row")
		}
		return nil
	}, WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly))); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	defer stop()
	if _, _, ok := next(); ok {
		t.Fatal("suspended iterator resumed access outside the system callback")
	}
}
