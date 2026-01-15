package cmp

import (
	"fmt"
	"testing"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/test/cmp/testdata"
)

type _testStandardSystem struct {
}

func (t *_testStandardSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	for i, p := range ecs.GetComponents[testdata.Point](ctx) {
		fmt.Printf("StandardExample, EnityIndex: %d, P:%v\n", i, p)
	}
	return nil
}

func (t *_testStandardSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("StandardExample"),
		ecs.WithDep[testdata.Point](),
		ecs.WithDep[testdata.Position](),
	)
	return nil
}

//go:generate go run ../../cmd/karmem/main.go build -golang -o ./testdata/ testdata/components.km
//go:generate go run ../../cmd/karmem/main.go fmt -s testdata/components.km
func TestCmpMain(t *testing.T) {
	world := ecs.NewWorld(ecs.WithWorldAutoOptimize())

	e := world.NewEntity()

	point1 := testdata.Point{
		X: 1, Y: 2, Z: 3,
	}

	pos1 := testdata.Position{
		X: 7,
	}

	name1 := testdata.Name{
		Value: ecs.NewFixed16("hello"),
	}

	e.Add(&point1)
	e.Add(&pos1)
	e.Add(&name1)

	// get and walk components by GetComponents
	sys := func(ctx *ecs.SystemContext, event ecs.Event) error {
		for i, p := range ecs.GetComponents[testdata.Point](ctx) {
			pos, ok := ecs.GetBuddy[testdata.Position](ctx, i)
			if ok {
				fmt.Printf("LightExample 1, EnityIndex: %d, P:%v, Pos:%v\n", i, p, pos)
			} else {
				fmt.Printf("LightExample 1, EnityIndex: %d, P:%v, Pos not found\n", i, p)
			}
		}
		return nil
	}

	// query
	sys2 := func(ctx *ecs.SystemContext, event ecs.Event) error {
		r := ecs.Query(ctx, ecs.WithComp[testdata.Point](), ecs.WithComp[testdata.Position]())
		for idx, _ := range r.Iter() {
			p, _ := ecs.GetBuddy[testdata.Point](ctx, idx)
			pos, _ := ecs.GetBuddy[testdata.Position](ctx, idx)
			name, _ := ecs.GetBuddy[testdata.Name](ctx, idx)
			fmt.Printf("LightExample 2, EnityIndex: %d, P:%v, Pos:%v, Name:%s\n", idx, p, pos, name.Value.String())
		}

		return nil
	}

	err := world.RegisterLight(sys,
		ecs.WithName("LightExample"),
		ecs.WithDep[testdata.Point](),
		ecs.WithDep[testdata.Position](),
	)
	if err != nil {
		t.Error(err)
	}

	err = world.RegisterLight(sys2,
		ecs.WithName("LightExample"),
		ecs.WithDep[testdata.Point](),
		ecs.WithDep[testdata.Position](ecs.ReadWrite),
		ecs.WithDep[testdata.Name](),
	)
	if err != nil {
		t.Error(err)
	}

	err = world.RegisterStandard(&_testStandardSystem{})
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 2; i++ {
		err = world.Update()
		if err != nil {
			t.Error(fmt.Errorf("%v", err))
			if e, ok := err.(*ecs.ExecuteErrors); ok {
				t.Error(fmt.Errorf("%v", e.SubErrs))
			}
		}
	}
}

func TestName(t *testing.T) {
	var a [16]byte
	b := [16]byte{1, 2, 3}

	copy(a[:], b[:])

	fmt.Printf("%v, %v\n", a, b)
}
