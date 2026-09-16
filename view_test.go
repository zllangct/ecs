package ecs

import (
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// groupCompBView 手写只读视图（codegen 视图的等价物）
type groupCompBView struct{ b *groupCompB }

func (v groupCompBView) ComponentPacketIdentifier() rockmem.PacketIdentifier {
	return 99999999991002
}
func (v groupCompBView) FromPtr(p unsafe.Pointer) groupCompBView {
	return groupCompBView{b: (*groupCompB)(p)}
}
func (v groupCompBView) V() int32 { return v.b.V }

type viewRecorder struct {
	pairs map[EntityIndex][2]int32
	count int
}

// View2 双写版本 probe
type view2Probe struct{ rec *viewRecorder }

func (s *view2Probe) Update(ctx *SystemContext, event Event) error {
	s.rec.pairs = map[EntityIndex][2]int32{}
	s.rec.count = 0
	for a, b := range View2[groupCompA, groupCompB](ctx) {
		s.rec.count++
		s.rec.pairs[0] = [2]int32{a.V, b.V} // 无 index 语义，仅记录值
		a.V++                               // 写穿透验证
	}
	return nil
}

// View2RO probe：A 可写 + B 只读视图，记录实体与值
type view2ROProbe struct{ rec *viewRecorder }

func (s *view2ROProbe) Update(ctx *SystemContext, event Event) error {
	s.rec.pairs = map[EntityIndex][2]int32{}
	s.rec.count = 0
	q := ctx.NewQuery(WithComp[groupCompA, *groupCompA](), WithComp[groupCompB, *groupCompB]())
	_ = q
	for a, bv := range View2RO[groupCompA, groupCompBView](ctx) {
		s.rec.count++
		s.rec.pairs[EntityIndex(s.rec.count)] = [2]int32{a.V, bv.V()}
		a.V += 10
	}
	return nil
}

// View2：只命中组表行（集齐 A,B 的实体），写穿透生效
func TestView2_Lockstep(t *testing.T) {
	w := NewWorld()
	rec := &viewRecorder{}
	err := w.RegisterLight(LightSystem((&view2Probe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2))) // 残段不命中
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 1 {
		t.Fatalf("View2 count = %d, want 1", rec.count)
	}
	if rec.pairs[0] != [2]int32{1, 11} {
		t.Fatalf("pairs = %v", rec.pairs)
	}
	// 写穿透：a.V++ 应落到表内存
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 2 {
		t.Fatalf("write-through failed: A.V = %d, want 2", got.V)
	}
}

// View2RO：只读视图版本
func TestView2RO_Lockstep(t *testing.T) {
	w := NewWorld()
	rec := &viewRecorder{}
	err := w.RegisterLight(LightSystem((&view2ROProbe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 1 {
		t.Fatalf("View2RO count = %d, want 1", rec.count)
	}
	if rec.pairs[EntityIndex(1)] != [2]int32{1, 11} {
		t.Fatalf("pairs = %v", rec.pairs)
	}
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 11 {
		t.Fatalf("write-through failed: A.V = %d, want 11", got.V)
	}
}

// View2 权限：B 为只读依赖时，可写 View2 静默返回空
func TestView2_ReadonlyRejected(t *testing.T) {
	w := NewWorld()
	rec := &viewRecorder{}
	err := w.RegisterLight(LightSystem((&view2Probe{rec: rec}).Update),
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
	if rec.count != 0 {
		t.Fatalf("View2 with readonly B should be empty, got %d", rec.count)
	}
}

// View2 回退：未建组（WithoutGroups）时仍返回正确结果
func TestView2_FallbackNonGrouped(t *testing.T) {
	w := NewWorld(WithoutGroups())
	rec := &viewRecorder{}
	err := w.RegisterLight(LightSystem((&view2Probe{rec: rec}).Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(11)))
	w.NewEntity(WithComponents(newGroupCompA(2)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if rec.count != 1 {
		t.Fatalf("fallback count = %d, want 1", rec.count)
	}
	if rec.pairs[0] != [2]int32{1, 11} {
		t.Fatalf("fallback pairs = %v", rec.pairs)
	}
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 2 {
		t.Fatalf("fallback write-through failed: A.V = %d", got.V)
	}
}
