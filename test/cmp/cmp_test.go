package cmp

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zllangct/ecs"
	karmem "github.com/zllangct/ecs/karmem"
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

func BenchmarkGet(b *testing.B) {
	const Size = 10000

	set := ecs.NewUSet[testdata.Point]()
	setp := ecs.NewUSet[testdata.PointViewer]()

	writer := karmem.NewWriter(16)

	p := &testdata.Point{}
	p.WriteAsRoot(writer)

	reader := karmem.NewReader(writer.Bytes())
	pv := testdata.NewPointViewer(reader, 0)

	for i := 0; i < Size; i++ {
		set.Add(p)
		setp.Add(pv)
	}

	b.Run("p", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			point := set.Get(int64(i % set.Len()))
			x := point.X
			y := point.Y
			z := point.Z
			_, _, _ = x, y, z
		}
	})

	b.Run("pv", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pointV := setp.Get(int64(i % set.Len()))
			x := pointV.X()
			y := pointV.Y()
			z := pointV.Z()
			_, _, _ = x, y, z
		}
	})
}

func TestED(t *testing.T) {
	p := testdata.Name{
		Value: ecs.NewFixed16("hello"),
		Points: [2]testdata.Point{
			testdata.Point{1, 1, 1},
		},
	}

	writer := karmem.NewWriter(0)
	p.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	r := testdata.NewNameViewer(reader, 0)
	if r.Value() != p.Value.String() {
		t.Error("not equal")
	}

	w := testdata.NewNameSource(reader, 0)
	_ = w
	w.SetValue("world1")
	if r.Value() != "world1" {
		t.Error("not equal")
	}
	w.SetValue("hello")

	w.SetPoints([2]testdata.Point{
		testdata.Point{2, 2, 2},
	})
	if r.Points()[0].X() != 2 {
		t.Error("not equal")
	}

	w.SetPoints([2]testdata.Point{
		{1, 1, 1},
	})

	w.SetArr([2]int32{1, 1})
	if w.Arr()[0] != 1 {
		t.Error("not equal")
	}
	w.SetArr([2]int32{})

	ref3new := testdata.Name{}
	ref3new.ReadAsRoot(reader)

	opt := cmp.Comparer(func(x, y ecs.Fixed16) bool {
		if x.String() == y.String() {
			return true
		}
		return false
	})

	if !cmp.Equal(p, ref3new, opt) {
		t.Error("not equal")
	}
}

func TestEDI(t *testing.T) {
	p := testdata.Point{
		X: 1,
		Y: 1,
		Z: 1,
	}

	writer := karmem.NewWriter(0)
	p.WriteAsRoot(writer)
	reader := karmem.NewReader(writer.Bytes())

	r := testdata.NewPointViewer(reader, 0)
	if r.X() != p.X {
		t.Error("not equal")
	}

	w := testdata.NewPointSource(reader, 0)
	_ = w
	w.SetX(3)

	ref3new := testdata.Point{}
	ref3new.ReadAsRoot(reader)

	if !cmp.Equal(ref3new.X, float32(3)) {
		t.Error("not equal")
	}
}

func TestName(t *testing.T) {
	var a [16]byte
	b := [16]byte{1, 2, 3}

	copy(a[:], b[:])

	fmt.Printf("%v, %v\n", a, b)
}
