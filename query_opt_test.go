package ecs

import "testing"

type queryOptRecorder struct {
	getA  map[EntityIndex]int32
	getB  map[EntityIndex]int32
	hits  int
	roNil bool // 只读依赖经 QueryGet 可写访问应返回 nil
}

type queryOptProbe struct {
	rec      *queryOptRecorder
	readonly bool
}

func (s *queryOptProbe) Update(ctx *SystemContext, event Event) error {
	s.rec.getA = map[EntityIndex]int32{}
	s.rec.getB = map[EntityIndex]int32{}
	s.rec.hits = 0
	// 同一帧内两次 NewQuery（验证缓存路径不破坏行为）
	for round := 0; round < 2; round++ {
		q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
		for index := range q.Iter() {
			s.rec.hits++
			if a, ok := QueryGet[groupCompA, *groupCompA](&q, index); ok {
				s.rec.getA[index] = a.V
			}
			if b, ok := QueryGet[groupCompB, *groupCompB](&q, index); ok {
				s.rec.getB[index] = b.V
			} else if s.readonly {
				s.rec.roNil = true
			}
		}
	}
	return nil
}

// QueryGet：组表内实体直取命中，值正确
func TestQueryGet_GroupTable(t *testing.T) {
	w := NewWorld()
	rec := &queryOptRecorder{}
	err := w.RegisterLight(LightSystem((&queryOptProbe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2))) // 残段，不应命中查询
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.hits != 2 { // 两轮 × 1 实体
		t.Fatalf("hits = %d, want 2", rec.hits)
	}
	if rec.getA[e1.Index()] != 1 || rec.getB[e1.Index()] != 11 {
		t.Fatalf("QueryGet values: A=%d B=%d", rec.getA[e1.Index()], rec.getB[e1.Index()])
	}
}

// QueryGet：只读依赖不可经可写访问器获取（静默失败约定）
func TestQueryGet_ReadonlyRejected(t *testing.T) {
	w := NewWorld()
	rec := &queryOptRecorder{}
	probe := &queryOptProbe{rec: rec, readonly: true}
	err := w.RegisterLight(LightSystem(probe.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if !rec.roNil {
		t.Fatal("readonly buddy via QueryGet should return nil")
	}
}

// QueryGet：非组（纯 CSet）路径同样工作
func TestQueryGet_CSetPath(t *testing.T) {
	w := NewWorld(WithoutGroups())
	rec := &queryOptRecorder{}
	err := w.RegisterLight(LightSystem((&queryOptProbe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()), // 被 WithoutGroups 忽略
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(7), newGroupCompB(77)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.getA[e1.Index()] != 7 || rec.getB[e1.Index()] != 77 {
		t.Fatalf("QueryGet CSet path: A=%d B=%d", rec.getA[e1.Index()], rec.getB[e1.Index()])
	}
}
