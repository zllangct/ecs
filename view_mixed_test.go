package ecs

import (
	"iter"
	"testing"
)

type mixedFilteredPool struct {
	*CSet[p2CompB, *p2CompB]
	hidden EntityIndex
}

func (p *mixedFilteredPool) Exist(key EntityIndex) bool {
	return key != p.hidden && p.CSet.Exist(key)
}

func TestMixedView2ROHonorsCustomPoolMembership(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity(WithComponents(&p2CompA{}, &p2CompB{}, &p2CompC{}))
	w.NewEntity(WithComponents(&p2CompB{})) // A is the smaller driving column.
	err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		id := GetIntType[p2CompB]()
		pool := w.components[id].(*CSet[p2CompB, *p2CompB])
		w.components[id] = &mixedFilteredPool{pool, e.Index()}
		for range View2RO[p2CompA, p2CompBView](ctx) {
			t.Fatal("custom component pool excluded this key")
		}
		return nil
	}, WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly), Dep[p2CompC]()), WithGroup(Dep[p2CompA](), Dep[p2CompC]()))
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Step(1); err != nil {
		t.Fatal(err)
	}
}

func TestMixedView2ROMatchesSuspendedQuery(t *testing.T) {
	for _, driver := range []string{"group", "pool", "tie"} {
		t.Run(driver, func(t *testing.T) {
			w := NewWorld()
			var initial []Entity
			for i := range 12 {
				components := []Component{&p2CompA{int32(i)}, &p2CompB{int32(100 + i)}}
				if i%2 == 0 {
					components = append(components, &p2CompC{})
				}
				initial = append(initial, w.NewEntity(WithComponents(components...)))
			}
			switch driver {
			case "group":
				w.NewEntity(WithComponents(&p2CompB{}))
			case "pool":
				w.NewEntity(WithComponents(&p2CompA{}))
			}
			type nextRow func() (*p2CompA, p2CompBView, bool)
			var got, want nextRow
			var stop, stopRef func()
			check := func() bool {
				ga, gb, gok := got()
				wa, wb, wok := want()
				if ga != wa || gb != wb || gok != wok {
					t.Fatalf("mismatched current pointers got %v/%v/%t want %v/%v/%t", ga, gb, gok, wa, wb, wok)
				}
				return gok
			}
			opts := []SystemOption{WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly), Dep[p2CompC]()), WithGroup(Dep[p2CompA](), Dep[p2CompC]())}
			err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
				if got == nil {
					q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompB]())
					ref := func(yield func(*p2CompA, p2CompBView) bool) {
						for k := range q.Iter() {
							a, aOK := QueryGet[p2CompA](&q, k)
							b, bOK := ctx.GetBuddyReadOnly[p2CompBView](k)
							if aOK && bOK && !yield(a, b) {
								return
							}
						}
					}
					want, stopRef = iter.Pull2(ref)
					got, stop = iter.Pull2(View2RO[p2CompA, p2CompBView](ctx))
					if !check() {
						t.Fatal("empty fixture")
					}
				} else {
					n := 0
					for check() {
						n++
					}
					if n == 0 {
						t.Fatal("fixture lost every remaining row")
					}
				}
				return nil
			}, opts...)
			if err != nil {
				t.Fatal(err)
			}
			if err := w.Update(); err != nil {
				t.Fatal(err)
			}
			defer stop()
			defer stopRef()
			info, _ := w.GetEntityInfo(initial[1])
			info.Add(&p2CompC{})
			info, _ = w.GetEntityInfo(initial[6])
			info.Remove(&p2CompC{})
			for i := range 128 {
				components := []Component{&p2CompA{int32(1000 + i)}, &p2CompB{int32(2000 + i)}}
				if i%2 == 0 {
					components = append(components, &p2CompC{})
				}
				w.NewEntity(WithComponents(components...))
			}
			if err := w.Update(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
