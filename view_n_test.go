package ecs

import (
	"testing"
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

// ============================================================================
// View3-View8 测试组件（groupCompA/B 复用 world_group_test.go，D 已被占为 disposable）
// ============================================================================

type gCompC struct{ V int32 }
type gCompD struct{ V int32 }
type gCompE struct{ V int32 }
type gCompF struct{ V int32 }
type gCompG struct{ V int32 }
type gCompH struct{ V int32 }

func (c *gCompC) NewComponentSet() ComponentSet              { return NewCSet[gCompC]() }
func (c *gCompC) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992001 }
func (c *gCompC) IsNomadic() bool                            { return false }
func (c *gCompC) IsDisposable() bool                         { return false }

func (c *gCompD) NewComponentSet() ComponentSet              { return NewCSet[gCompD]() }
func (c *gCompD) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992002 }
func (c *gCompD) IsNomadic() bool                            { return false }
func (c *gCompD) IsDisposable() bool                         { return false }

func (c *gCompE) NewComponentSet() ComponentSet              { return NewCSet[gCompE]() }
func (c *gCompE) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992003 }
func (c *gCompE) IsNomadic() bool                            { return false }
func (c *gCompE) IsDisposable() bool                         { return false }

func (c *gCompF) NewComponentSet() ComponentSet              { return NewCSet[gCompF]() }
func (c *gCompF) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992004 }
func (c *gCompF) IsNomadic() bool                            { return false }
func (c *gCompF) IsDisposable() bool                         { return false }

func (c *gCompG) NewComponentSet() ComponentSet              { return NewCSet[gCompG]() }
func (c *gCompG) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992005 }
func (c *gCompG) IsNomadic() bool                            { return false }
func (c *gCompG) IsDisposable() bool                         { return false }

func (c *gCompH) NewComponentSet() ComponentSet              { return NewCSet[gCompH]() }
func (c *gCompH) PacketIdentifier() rockmem.PacketIdentifier { return 99999999992006 }
func (c *gCompH) IsNomadic() bool                            { return false }
func (c *gCompH) IsDisposable() bool                         { return false }

func init() {
	RegisterComponent[gCompC]("test_group_n")
	RegisterComponent[gCompD]("test_group_n")
	RegisterComponent[gCompE]("test_group_n")
	RegisterComponent[gCompF]("test_group_n")
	RegisterComponent[gCompG]("test_group_n")
	RegisterComponent[gCompH]("test_group_n")
}

// 只读视图
type gCompCView struct{ b *gCompC }
type gCompDView struct{ b *gCompD }
type gCompEView struct{ b *gCompE }
type gCompFView struct{ b *gCompF }
type gCompGView struct{ b *gCompG }
type gCompHView struct{ b *gCompH }

func (v gCompCView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992001 }
func (v gCompCView) FromPtr(p unsafe.Pointer) gCompCView                 { return gCompCView{(*gCompC)(p)} }
func (v gCompCView) V() int32                                            { return v.b.V }

func (v gCompDView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992002 }
func (v gCompDView) FromPtr(p unsafe.Pointer) gCompDView                 { return gCompDView{(*gCompD)(p)} }
func (v gCompDView) V() int32                                            { return v.b.V }

func (v gCompEView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992003 }
func (v gCompEView) FromPtr(p unsafe.Pointer) gCompEView                 { return gCompEView{(*gCompE)(p)} }
func (v gCompEView) V() int32                                            { return v.b.V }

func (v gCompFView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992004 }
func (v gCompFView) FromPtr(p unsafe.Pointer) gCompFView                 { return gCompFView{(*gCompF)(p)} }
func (v gCompFView) V() int32                                            { return v.b.V }

func (v gCompGView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992005 }
func (v gCompGView) FromPtr(p unsafe.Pointer) gCompGView                 { return gCompGView{(*gCompG)(p)} }
func (v gCompGView) V() int32                                            { return v.b.V }

func (v gCompHView) ComponentPacketIdentifier() rockmem.PacketIdentifier { return 99999999992006 }
func (v gCompHView) FromPtr(p unsafe.Pointer) gCompHView                 { return gCompHView{(*gCompH)(p)} }
func (v gCompHView) V() int32                                            { return v.b.V }

// ============================================================================
// View3 lockstep：值正确、残段排除、写穿透
// ============================================================================

type view3Probe struct{ count, sumA, sumB, sumC int }

func (s *view3Probe) Update(ctx *SystemContext, event Event) error {
	s.count, s.sumA, s.sumB, s.sumC = 0, 0, 0, 0
	for t := range View3[groupCompA, groupCompB, gCompC](ctx) {
		s.count++
		s.sumA += int(t.V1.V)
		s.sumB += int(t.V2.V)
		s.sumC += int(t.V3.V)
		t.V1.V++
	}
	return nil
}

func TestView3_Lockstep(t *testing.T) {
	w := NewWorld()
	probe := &view3Probe{}
	err := w.RegisterLight(LightSystem(probe.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}))
	w.NewEntity(WithComponents(newGroupCompA(4), newGroupCompB(5))) // 残段（缺C）排除
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if probe.count != 1 || probe.sumA != 1 || probe.sumB != 2 || probe.sumC != 3 {
		t.Fatalf("count=%d sums=%d,%d,%d", probe.count, probe.sumA, probe.sumB, probe.sumC)
	}
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 2 {
		t.Fatalf("write-through failed: A.V = %d", got.V)
	}
}

// ============================================================================
// View8 lockstep：8 列指针运算正确性
// ============================================================================

type view8Probe struct {
	count int
	sums  [8]int
}

func (s *view8Probe) Update(ctx *SystemContext, event Event) error {
	s.count = 0
	s.sums = [8]int{}
	for t := range View8[groupCompA, groupCompB, gCompC, gCompD, gCompE, gCompF, gCompG, gCompH](ctx) {
		s.count++
		s.sums[0] += int(t.V1.V)
		s.sums[1] += int(t.V2.V)
		s.sums[2] += int(t.V3.V)
		s.sums[3] += int(t.V4.V)
		s.sums[4] += int(t.V5.V)
		s.sums[5] += int(t.V6.V)
		s.sums[6] += int(t.V7.V)
		s.sums[7] += int(t.V8.V)
		t.V8.V++
	}
	return nil
}

func view8Deps() []SystemOption {
	return []SystemOption{
		WithDeps(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC](), Dep[gCompD](),
			Dep[gCompE](), Dep[gCompF](), Dep[gCompG](), Dep[gCompH]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC](), Dep[gCompD](),
			Dep[gCompE](), Dep[gCompF](), Dep[gCompG](), Dep[gCompH]()),
	}
}

func TestView8_Lockstep(t *testing.T) {
	w := NewWorld()
	probe := &view8Probe{}
	if err := w.RegisterLight(LightSystem(probe.Update), view8Deps()...); err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(
		newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}, &gCompD{V: 4},
		&gCompE{V: 5}, &gCompF{V: 6}, &gCompG{V: 7}, &gCompH{V: 8},
	))
	// 缺 H 的残段排除
	w.NewEntity(WithComponents(
		newGroupCompA(10), newGroupCompB(10), &gCompC{V: 10}, &gCompD{V: 10},
		&gCompE{V: 10}, &gCompF{V: 10}, &gCompG{V: 10},
	))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if probe.count != 1 {
		t.Fatalf("count = %d, want 1", probe.count)
	}
	for i, want := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		if probe.sums[i] != want {
			t.Fatalf("sums[%d] = %d, want %d (all %v)", i, probe.sums[i], want, probe.sums)
		}
	}
	got, _ := w.GetComponent[gCompH, *gCompH](e1)
	if got.V != 9 {
		t.Fatalf("write-through failed: H.V = %d", got.V)
	}
}

// ============================================================================
// View3RO / View8RO lockstep
// ============================================================================

type view3ROProbe struct{ count, sumA, sumB, sumC int }

func (s *view3ROProbe) Update(ctx *SystemContext, event Event) error {
	s.count, s.sumA, s.sumB, s.sumC = 0, 0, 0, 0
	for t := range View3RO[groupCompA, groupCompBView, gCompCView](ctx) {
		s.count++
		s.sumA += int(t.V1.V)
		s.sumB += int(t.V2.V())
		s.sumC += int(t.V3.V())
		t.V1.V += 10
	}
	return nil
}

func TestView3RO_Lockstep(t *testing.T) {
	w := NewWorld()
	probe := &view3ROProbe{}
	err := w.RegisterLight(LightSystem(probe.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly), Dep[gCompC](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	e1 := w.NewEntity(WithComponents(newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}))
	w.NewEntity(WithComponents(newGroupCompA(4), newGroupCompB(5)))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if probe.count != 1 || probe.sumA != 1 || probe.sumB != 2 || probe.sumC != 3 {
		t.Fatalf("count=%d sums=%d,%d,%d", probe.count, probe.sumA, probe.sumB, probe.sumC)
	}
	got, _ := w.GetComponent[groupCompA, *groupCompA](e1)
	if got.V != 11 {
		t.Fatalf("write-through failed: A.V = %d", got.V)
	}
}

type view8ROProbe struct {
	count int
	sums  [8]int
}

func (s *view8ROProbe) Update(ctx *SystemContext, event Event) error {
	s.count = 0
	s.sums = [8]int{}
	for t := range View8RO[groupCompA, groupCompBView, gCompCView, gCompDView, gCompEView, gCompFView, gCompGView, gCompHView](ctx) {
		s.count++
		s.sums[0] += int(t.V1.V)
		s.sums[1] += int(t.V2.V())
		s.sums[2] += int(t.V3.V())
		s.sums[3] += int(t.V4.V())
		s.sums[4] += int(t.V5.V())
		s.sums[5] += int(t.V6.V())
		s.sums[6] += int(t.V7.V())
		s.sums[7] += int(t.V8.V())
	}
	return nil
}

func TestView8RO_Lockstep(t *testing.T) {
	w := NewWorld()
	probe := &view8ROProbe{}
	err := w.RegisterLight(LightSystem(probe.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly), Dep[gCompC](ReadOnly), Dep[gCompD](ReadOnly),
			Dep[gCompE](ReadOnly), Dep[gCompF](ReadOnly), Dep[gCompG](ReadOnly), Dep[gCompH](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC](), Dep[gCompD](),
			Dep[gCompE](), Dep[gCompF](), Dep[gCompG](), Dep[gCompH]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	w.NewEntity(WithComponents(
		newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}, &gCompD{V: 4},
		&gCompE{V: 5}, &gCompF{V: 6}, &gCompG{V: 7}, &gCompH{V: 8},
	))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if probe.count != 1 {
		t.Fatalf("count = %d, want 1", probe.count)
	}
	for i, want := range []int{1, 2, 3, 4, 5, 6, 7, 8} {
		if probe.sums[i] != want {
			t.Fatalf("sums[%d] = %d, want %d", i, probe.sums[i], want)
		}
	}
}

// ============================================================================
// View4-View7 / View4RO-View7RO 冒烟（同一模板的中间元数，验证实例化与基本行为）
// ============================================================================

type viewSmokeProbe struct{ count int }

func (s *viewSmokeProbe) Update(ctx *SystemContext, event Event) error {
	s.count = 0
	for range View4[groupCompA, groupCompB, gCompC, gCompD](ctx) {
		s.count++
	}
	for range View5[groupCompA, groupCompB, gCompC, gCompD, gCompE](ctx) {
		s.count++
	}
	for range View6[groupCompA, groupCompB, gCompC, gCompD, gCompE, gCompF](ctx) {
		s.count++
	}
	for range View7[groupCompA, groupCompB, gCompC, gCompD, gCompE, gCompF, gCompG](ctx) {
		s.count++
	}
	for range View4RO[groupCompA, groupCompBView, gCompCView, gCompDView](ctx) {
		s.count++
	}
	for range View5RO[groupCompA, groupCompBView, gCompCView, gCompDView, gCompEView](ctx) {
		s.count++
	}
	for range View6RO[groupCompA, groupCompBView, gCompCView, gCompDView, gCompEView, gCompFView](ctx) {
		s.count++
	}
	for range View7RO[groupCompA, groupCompBView, gCompCView, gCompDView, gCompEView, gCompFView, gCompGView](ctx) {
		s.count++
	}
	return nil
}

func TestViewN_Smoke(t *testing.T) {
	w := NewWorld()
	probe := &viewSmokeProbe{}
	// 声明一个大组 {A..H}，View4-7 取其前缀子集（同属一张组表即可 lockstep）
	err := w.RegisterLight(LightSystem(probe.Update),
		WithDeps(Dep[groupCompA](), Dep[groupCompB](ReadOnly), Dep[gCompC](ReadOnly), Dep[gCompD](ReadOnly),
			Dep[gCompE](ReadOnly), Dep[gCompF](ReadOnly), Dep[gCompG](ReadOnly), Dep[gCompH](ReadOnly)),
		WithGroup(Dep[groupCompA](), Dep[groupCompB](), Dep[gCompC](), Dep[gCompD](),
			Dep[gCompE](), Dep[gCompF](), Dep[gCompG](), Dep[gCompH]()),
	)
	if err != nil {
		t.Fatal(err)
	}
	// View4-7（全可写）权限不满足（B..G 只读）→ 空；View4RO-7RO 应各命中 1
	w.NewEntity(WithComponents(
		newGroupCompA(1), newGroupCompB(2), &gCompC{V: 3}, &gCompD{V: 4},
		&gCompE{V: 5}, &gCompF{V: 6}, &gCompG{V: 7}, &gCompH{V: 8},
	))
	if err := w.Update(); err != nil {
		t.Fatal(err)
	}
	if probe.count != 4 {
		t.Fatalf("count = %d, want 4 (View4RO-7RO 各命中 1；View4-7 权限拒绝)", probe.count)
	}
}
