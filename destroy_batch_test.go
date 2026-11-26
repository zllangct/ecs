package ecs

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestWorldDestroyBatchPreservesDelayWindow(t *testing.T) {
	w := NewWorld()
	reference := NewEntityIDGenerator(1024, 10)
	var entities []Entity
	for range 140 {
		id := w.NewEntity()
		if id != reference.NewID() {
			t.Fatal("initial IDs differ")
		}
		entities = append(entities, id)
	}
	runBatch := func(valid []Entity) {
		// Duplicate requests and the stale handle for a living entity must
		// have no effect. Apply valid requests in canonical order to the
		// public, one-at-a-time API as an independent reference.
		for _, id := range slices.Backward(valid) {
			w.DestroyEntity(id)
			w.DestroyEntity(id)
		}
		w.DestroyEntity(Entity(1<<32) | entities[139])
		for _, id := range valid {
			reference.mergeHint = 0 // Original merge starts at the head for every window.
			reference.FreeID(id)
		}
		if err := w.Step(time.Millisecond); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(w.idGenerator.Marshal(), reference.Marshal()) {
			t.Fatal("serialized freelist or delay buffer differs from individual frees")
		}
	}
	// The previous frame leaves high indices pending. The next merge
	// starts below them and crosses several complete windows.
	runBatch(entities[110:117])
	snapshot, err := w.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	w, err = NewWorldFromData(snapshot)
	if err != nil {
		t.Fatal(err)
	}
	runBatch(entities[:36])
	for range 17 {
		if got, want := w.NewEntity(), reference.NewID(); got != want {
			t.Fatalf("resumed reuse sequence %d, want %d", got, want)
		}
	}
	runBatch(entities[50:93])
	for range 100 {
		if got, want := w.NewEntity(), reference.NewID(); got != want {
			t.Fatalf("final reuse sequence %d, want %d", got, want)
		}
	}
}

func TestEntityFreeMergeHintInvalidation(t *testing.T) {
	g := NewEntityIDGenerator(64, 5)
	var ids []Entity
	for range 48 {
		ids = append(ids, g.NewID())
	}
	allAllocated := g.Marshal()
	for _, id := range ids[:5] {
		g.FreeID(id)
	}
	if g.mergeHint == 0 {
		t.Fatal("merge did not retain a predecessor")
	}
	if err := g.Unmarshal(allAllocated); err != nil {
		t.Fatal(err)
	}
	if g.mergeHint != 0 {
		t.Fatal("restore retained a predecessor from another freelist")
	}
	for _, id := range ids[10:15] {
		g.FreeID(id)
	}
	for i := 10; i < 15; i++ {
		got := g.NewID().toReuseID()
		if got.index != ids[i].Index() || got.reuse != 1 || g.mergeHint != 0 {
			t.Fatalf("bad reuse or cached predecessor after allocation: %+v", got)
		}
	}
}

func TestWorldDestroyBatchStorageAndOrder(t *testing.T) {
	run := func(reverse bool) *SerializableWorldData {
		w := NewWorld()
		observedNomadic := int32(0)
		var transientOwner Entity
		if err := w.RegisterLight(func(ctx *SystemContext, evt Event) error {
			if evt.Frame == 1 {
				ctx.AddComponents(transientOwner, newGroupCompD(8))
				w.AddNomadic(&nomadicComp{Seq: 31})
			}
			for _, n := range ctx.GetComponents[nomadicComp]() {
				observedNomadic = n.Seq
			}
			return nil
		}, WithDep[nomadicComp](),
			WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
			WithGroup(Dep[gCompC](), Dep[gCompD]())); err != nil {
			t.Fatal(err)
		}
		all := []Component{newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}, &gCompD{V: 4}, &gCompE{V: 5}}
		remove := []Entity{
			w.NewEntity(WithComponents(all...)),
			w.NewEntity(WithComponents(newGroupCompB(22), &gCompC{V: 33})),
			w.NewEntity(),
		}
		first := w.NewEntity(WithComponents(all...))
		second := w.NewEntity(WithComponents(all...))
		transientOwner = remove[0]
		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		if _, ok := w.GetComponent[groupCompD](transientOwner); !ok {
			t.Fatal("disposable component missing before destroy")
		}
		if reverse {
			slices.Reverse(remove)
		}
		for _, id := range remove {
			w.DestroyEntity(id)
		}
		// This addition is pending at destruction and must not create an
		// orphan. The unrelated nomadic value must reach the system.
		info, _ := w.GetEntityInfo(remove[1])
		info.Add(newGroupCompD(9))
		if err := w.Update(); err != nil {
			t.Fatal(err)
		}
		if observedNomadic != 31 {
			t.Fatalf("nomadic value %d, want 31", observedNomadic)
		}
		for _, id := range remove {
			if _, ok := w.GetEntityInfo(id); ok {
				t.Fatalf("entity %d survived", id)
			}
			for _, set := range w.components {
				if set.Exist(id.Index()) {
					t.Fatalf("residual component for destroyed entity %d", id)
				}
			}
			for _, a := range w.archetypes.groups {
				if _, ok := a.rowOf(id.Index()); ok {
					t.Fatalf("group row for destroyed entity %d", id)
				}
			}
		}
		for _, id := range []Entity{first, second} {
			if a, ok := w.GetComponent[groupCompA](id); !ok || a.V != 1 {
				t.Fatalf("changed surviving group row for %d", id)
			}
		}
		// The snapshot stores component sets. Check live group order as well.
		rows := []EntityIndex{second.Index(), first.Index()}
		for _, group := range w.archetypes.groups {
			if !slices.Equal(group.entities, rows) {
				t.Fatalf("surviving group order %v, want %v", group.entities, rows)
			}
		}
		snap, err := w.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		return snap
	}
	if !reflect.DeepEqual(run(false), run(true)) {
		t.Fatal("batch enqueue order changed the snapshot")
	}
}
