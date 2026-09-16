package ecs

import "testing"

type groupAccessRecorder struct {
	sawA     map[EntityIndex]int32
	buddyB   map[EntityIndex]int32
	queryHit map[EntityIndex]bool
}

type groupAccessProbe struct{ rec *groupAccessRecorder }

func (s *groupAccessProbe) Update(ctx *SystemContext, event Event) error {
	s.rec.sawA = map[EntityIndex]int32{}
	s.rec.buddyB = map[EntityIndex]int32{}
	s.rec.queryHit = map[EntityIndex]bool{}
	for idx, a := range ctx.GetComponents[groupCompA, *groupCompA]() {
		s.rec.sawA[idx] = a.V
		if b, ok := ctx.GetBuddy[groupCompB, *groupCompB](idx); ok {
			s.rec.buddyB[idx] = b.V
		}
	}
	q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
	for idx := range q.Iter() {
		s.rec.queryHit[idx] = true
	}
	return nil
}

// GetComponents[A] 覆盖表内+残段；GetBuddy[B] 表内直取；Query{A,B} 以组表为主迭代
func TestGroupAccess_Iteration(t *testing.T) {
	w := NewWorld()
	rec := &groupAccessRecorder{}
	err := w.RegisterLight(LightSystem((&groupAccessProbe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11))) // 表内
	e2 := w.NewEntity(WithComponents(newGroupCompA(2)))                    // 残段（只有A）
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.sawA[e1.Index()] != 1 || rec.sawA[e2.Index()] != 2 {
		t.Fatalf("GetComponents[A] = %v", rec.sawA)
	}
	if rec.buddyB[e1.Index()] != 11 {
		t.Fatalf("GetBuddy[B] for e1 = %d", rec.buddyB[e1.Index()])
	}
	if _, ok := rec.buddyB[e2.Index()]; ok {
		t.Fatal("e2 has no B")
	}
	if !rec.queryHit[e1.Index()] || rec.queryHit[e2.Index()] {
		t.Fatalf("query result wrong: %v", rec.queryHit)
	}
}

// 表内组件经 GetBuddy 写入，值落表（下行从表读回）
func TestGroupAccess_WriteThrough(t *testing.T) {
	w := NewWorld()
	rec := &groupAccessRecorder{}
	err := w.RegisterLight(LightSystem((&groupAccessProbe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	got, ok := w.GetComponent[groupCompA, *groupCompA](e)
	if !ok {
		t.Fatal("missing A")
	}
	got.V = 42 // 直写表内存
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.sawA[e.Index()] != 42 {
		t.Fatalf("saw A = %d, want 42", rec.sawA[e.Index()])
	}
}
