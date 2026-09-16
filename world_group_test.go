package ecs

import (
	"testing"

	rockmem "github.com/zllangct/rockmem/golang"
)

type groupCompA struct{ V int32 }

func (c *groupCompA) NewComponentSet() ComponentSet                { return NewCSet[groupCompA]() }
func (c *groupCompA) PacketIdentifier() rockmem.PacketIdentifier   { return 99999999991001 }
func (c *groupCompA) IsNomadic() bool                              { return false }
func (c *groupCompA) IsDisposable() bool                           { return false }

type groupCompB struct{ V int32 }

func (c *groupCompB) NewComponentSet() ComponentSet                { return NewCSet[groupCompB]() }
func (c *groupCompB) PacketIdentifier() rockmem.PacketIdentifier   { return 99999999991002 }
func (c *groupCompB) IsNomadic() bool                              { return false }
func (c *groupCompB) IsDisposable() bool                           { return false }

type groupCompD struct{ V int32 }

func (c *groupCompD) NewComponentSet() ComponentSet                { return NewCSet[groupCompD]() }
func (c *groupCompD) PacketIdentifier() rockmem.PacketIdentifier   { return 99999999991003 }
func (c *groupCompD) IsNomadic() bool                              { return false }
func (c *groupCompD) IsDisposable() bool                           { return true }

func newGroupCompA(v int32) *groupCompA { return &groupCompA{V: v} }
func newGroupCompB(v int32) *groupCompB { return &groupCompB{V: v} }
func newGroupCompD(v int32) *groupCompD { return &groupCompD{V: v} }

func init() {
	RegisterComponent[groupCompA]("test_group")
	RegisterComponent[groupCompB]("test_group")
	RegisterComponent[groupCompD]("test_group")
}

func groupProbeDeps() []SystemOption {
	return []SystemOption{
		WithDeps(Dep[groupCompA](), Dep[groupCompB]()),
	}
}

// groupProbeLight 空 Update，仅承载依赖/组声明
func groupProbeLight(ctx *SystemContext, event Event) error { return nil }

func TestWorld_GroupDeclaredAndBuilt(t *testing.T) {
	w := NewWorld()
	err := w.RegisterLight(LightSystem(groupProbeLight),
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
	if !w.archetypes.sealed {
		t.Fatal("registry should be sealed after first Update")
	}
	a, ok := w.archetypes.owner(GetIntType[groupCompA, *groupCompA]())
	if !ok {
		t.Fatal("group not built for groupCompA")
	}
	if _, ok := w.archetypes.owner(GetIntType[groupCompB, *groupCompB]()); !ok {
		t.Fatal("group not built for groupCompB")
	}
	// Task 6 之前：数据仍在 CSet，GetComponent 回落正确
	got, ok := w.GetComponent[groupCompA, *groupCompA](e)
	if !ok || got.V != 1 {
		t.Fatalf("GetComponent = %+v %v", got, ok)
	}
	_ = a
}

func TestWorld_GroupRejectsDisposable(t *testing.T) {
	w := NewWorld()
	err := w.RegisterLight(LightSystem(groupProbeLight),
		WithDeps(Dep[groupCompA](), Dep[groupCompD]()),
		WithGroup(Dep[groupCompA](), Dep[groupCompD]()),
	)
	if err == nil {
		t.Fatal("disposable component in group should fail registration")
	}
}

func TestWorld_GroupRejectsUnregistered(t *testing.T) {
	w := NewWorld()
	err := w.RegisterLight(LightSystem(groupProbeLight),
		WithGroup(NewItDependency(12345678901234), NewItDependency(12345678901235)),
	)
	if err == nil {
		t.Fatal("unregistered component in group should fail registration")
	}
}
