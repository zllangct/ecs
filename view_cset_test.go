package ecs

import (
	"iter"
	"slices"
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
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
				usage := ctx.info.getOptReporter().shapeUsage
				shape := FixedCompound2{GetIntType[p2CompA](), GetIntType[p2CompB]()}
				before := usage[shape]
				first := 0
				for range rows {
					first++
					break
				}
				if first != 1 {
					t.Fatalf("early break visited %d rows", first)
				}
				if calls := usage[shape] - before; layout == "CSet" && calls != 1 {
					t.Fatalf("early break reported %d scans, want 1", calls)
				}
				for a, bv := range rows {
					a.V += bv.V()
				}
				if calls := usage[shape] - before; layout == "CSet" && calls != 2 {
					t.Fatalf("two iterations reported %d scans, want 2", calls)
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
	t.Run("dense", func(t *testing.T) { checkSuspendedView2ROCSet(t, true) })
	t.Run("ordered-merge", func(t *testing.T) { checkSuspendedView2ROCSet(t, false) })
}

func checkSuspendedView2ROCSet(t *testing.T, dense bool) {
	w := NewWorld(WithoutGroups())
	// In the ordered-merge case, A={1,3,5,7,9} and B={1,3,4,7,9}.
	var initial [10]Entity
	for i := range initial {
		var comps []Component
		switch {
		case dense || i == 1 || i == 3 || i == 7 || i == 9:
			comps = []Component{&p2CompA{}, &p2CompB{}}
		case i == 4:
			comps = []Component{&p2CompB{}}
		case i == 5:
			comps = []Component{&p2CompA{}}
		}
		initial[i] = w.NewEntity(WithComponents(comps...))
	}
	var next func() (*p2CompA, p2CompBView, bool)
	var stop func()
	if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		if next == nil {
			next, stop = iter.Pull2(View2RO[p2CompA, p2CompBView](ctx))
			if _, _, ok := next(); !ok {
				t.Fatal("active system could not read a row")
			}
			return nil
		}
		if _, ok := ctx.GetBuddyReadOnly[p2CompBView](initial[5].Index()); !ok {
			t.Fatal("pending add was not flushed")
		}
		remaining := []int{3, 7} // The suspended merge's peer slice excludes the new B at 5.
		if dense {
			remaining = []int{1, 2}
		}
		for _, i := range remaining {
			a, b, ok := next()
			wantA, _ := ctx.GetBuddy[p2CompA](initial[i].Index())
			wantB, _ := ctx.GetBuddyReadOnly[p2CompBView](initial[i].Index())
			if !ok || a != wantA || b != wantB {
				t.Fatalf("resumed iterator did not read entity %d from the current pools", i)
			}
		}
		return nil
	}, WithDeps(Dep[p2CompA](), Dep[p2CompB](ReadOnly))); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	defer stop()
	info, _ := w.GetEntityInfo(initial[5])
	info.Add(&p2CompB{})
	// Force both dense arrays to grow while the iterator is suspended.
	for range 128 {
		w.NewEntity(WithComponents(&p2CompA{}, &p2CompB{}))
	}
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := next(); ok { // Both iterators still have matching rows.
		t.Fatal("suspended iterator resumed access outside the system callback")
	}
}

func TestView2ROCSetOrderAfterSwapRemove(t *testing.T) {
	w := NewWorld(WithoutGroups())
	var entities [8]Entity
	for i := range entities {
		var comps []Component
		if i > 0 {
			comps = append(comps, newGroupCompA(int32(i)))
		}
		if i < 7 {
			comps = append(comps, newGroupCompB(10*int32(i)))
		}
		entities[i] = w.NewEntity(WithComponents(comps...))
	}

	var rows iter.Seq2[*groupCompA, groupCompBSnapshot]
	var got []int32
	if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		if rows == nil {
			rows = View2RO[groupCompA, groupCompBSnapshot](ctx)
		}
		got = got[:0]
		for a, b := range rows {
			if b.V() != 10*a.V || b.doubled != 20*a.V {
				t.Fatalf("mismatched columns at A=%d: B=%d doubled=%d", a.V, b.V(), b.doubled)
			}
			got = append(got, a.V)
		}
		return nil
	}, WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly))); err != nil {
		t.Fatal(err)
	}
	update := func(want ...int32) {
		t.Helper()
		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("iteration order %v, want %v", got, want)
		}
	}
	// Initially both pools are dense, but start at different entity IDs.
	update(1, 2, 3, 4, 5, 6)
	info, _ := w.GetEntityInfo(entities[2])
	info.Remove(newGroupCompB(0))
	update(1, 6, 3, 4, 5) // B is smaller, with its last row moved to the hole.
	info, _ = w.GetEntityInfo(entities[3])
	info.Remove(newGroupCompA(0))
	update(1, 4, 5, 6) // Equal size: A has the lower component ID and leads the intersection.
}

type viewNomadicPeer struct{ nomadicComp }

func (*viewNomadicPeer) PacketIdentifier() rockmem.PacketIdentifier { return 99999999991004 }
func (*viewNomadicPeer) NewComponentSet() ComponentSet              { return NewCSet[viewNomadicPeer]() }
func (viewNomadicPeer) ComponentPacketIdentifier() rockmem.PacketIdentifier {
	return (*viewNomadicPeer)(nil).PacketIdentifier()
}
func (viewNomadicPeer) FromPtr(p unsafe.Pointer) viewNomadicPeer { return *(*viewNomadicPeer)(p) }

func TestView2ROCSetRequiresLiveEntity(t *testing.T) {
	w := NewWorld(WithoutGroups())
	w.AddNomadic(&nomadicComp{}, &viewNomadicPeer{})
	called := false
	if err := w.RegisterLight(func(ctx *SystemContext, _ Event) error {
		called = true
		for _, id := range []ComponentIntType{GetIntType[nomadicComp](), GetIntType[viewNomadicPeer]()} {
			pool, ok := w.getComponentSet(id)
			if !ok || !pool.Exist(0) {
				t.Fatal("expected both nomadic pools to have an ownerless row at key zero")
			}
		}
		for range View2RO[nomadicComp, viewNomadicPeer](ctx) {
			t.Fatal("two matching pool keys yielded a row with no live entity")
		}
		return nil
	}, WithDeps(Dep[nomadicComp](), Dep[viewNomadicPeer](ReadOnly))); err != nil {
		t.Fatal(err)
	}
	if err := w.Update(); err != nil || !called {
		t.Fatalf("expected system to scan nomadic components, called=%v, err=%v", called, err)
	}
}
