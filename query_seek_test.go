package ecs

import (
	"iter"
	"slices"
	"testing"
)

func TestSkewedQueryPreservesCapturedPeerCursors(t *testing.T) {
	for _, mutation := range []string{"grow", "swap-remove"} {
		t.Run(mutation, func(t *testing.T) {
			w := NewWorld(WithoutGroups())
			entities := make([]Entity, 257)
			for i := range entities {
				comps := []Component{&p2CompC{}}
				if i%64 == 0 {
					comps = append(comps, &p2CompA{})
				}
				if i != 128 {
					comps = append(comps, &p2CompB{})
				}
				entities[i] = w.NewEntity(WithComponents(comps...))
			}
			var next, reference func() (EntityIndex, *EntityInfo, bool)
			var stop, stopReference func()
			err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
				if next != nil {
					for {
						key, info, ok := next()
						wantKey, wantInfo, wantOK := reference()
						if key != wantKey || info != wantInfo || ok != wantOK {
							t.Fatalf("resumed row (%d,%p,%v), want (%d,%p,%v)", key, info, ok, wantKey, wantInfo, wantOK)
						}
						if !ok {
							break
						}
					}
					return nil
				}
				q := ctx.NewQuery(WithComp[p2CompA](), WithComp[p2CompB](), WithComp[p2CompC]())
				rows := q.Iter()
				for range rows {
					break
				}
				want := []EntityIndex{entities[0].Index(), entities[64].Index(), entities[192].Index(), entities[256].Index()}
				for range 2 {
					var got []EntityIndex
					for key, info := range rows {
						if info.Entity().Index() != key {
							t.Fatal("wrong entity metadata")
						}
						got = append(got, key)
					}
					if !slices.Equal(got, want) {
						t.Fatalf("got %v, want %v", got, want)
					}
				}
				next, stop = iter.Pull2(rows)
				reference, stopReference = iter.Pull2(mergeJoinIter(w, pickMinNonEmptyPool(q.cached.pools), q.cached.pools))
				first, _, ok := next()
				refFirst, _, refOK := reference()
				if first != want[0] || first != refFirst || !ok || !refOK {
					t.Fatal("expected the same first row")
				}
				return nil
			}, WithDeps(Dep[p2CompA](), Dep[p2CompB](), Dep[p2CompC]()))
			if err != nil {
				t.Fatal(err)
			}
			if err := w.Step(1); err != nil {
				t.Fatal(err)
			}
			defer stop()
			defer stopReference()
			if mutation == "grow" {
				info, _ := w.GetEntityInfo(entities[128])
				info.Add(&p2CompB{})
				for range 512 {
					w.NewEntity(WithComponents(&p2CompB{}, &p2CompC{}))
				}
			} else {
				info, _ := w.GetEntityInfo(entities[2])
				info.Remove(&p2CompB{})
			}
			if err := w.Step(1); err != nil {
				t.Fatal(err)
			}
		})
	}
}

type customOrderedPeer struct {
	*CSet[p2CompB, *p2CompB]
	keys []EntityIndex
}

func (p *customOrderedPeer) EntityIndexes() []EntityIndex { return p.keys }

func TestSkewedQueryUsesCustomPeerKeys(t *testing.T) {
	w := NewWorld(WithoutGroups())
	a, b, c := NewCSet[p2CompA](), &customOrderedPeer{CSet: NewCSet[p2CompB]()}, NewCSet[p2CompC]()
	var want []EntityIndex
	for i := range 512 {
		e := w.NewEntity()
		b.Add(e, &p2CompB{})
		c.Add(e, &p2CompC{})
		if i%3 != 0 {
			b.keys = append(b.keys, e.Index())
		}
		if i%64 == 0 {
			a.Add(e, &p2CompA{})
			if i%3 != 0 {
				want = append(want, e.Index())
			}
		}
	}
	var got []EntityIndex
	for key := range intersectPools(w, a, []ComponentSet{b, a, c}) {
		got = append(got, key)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}
