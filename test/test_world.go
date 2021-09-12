package main

import (
	"ecs"
	"time"
)

//position component
type Position struct {
	ecs.Component[Position]
	X int
	Y int
	Z int
}

//movement component
type Movement struct {
	ecs.Component[Movement]
	v   int
	dir []int
}

////move system
type MoveSystemData struct {
	movement *Movement
	position *Position
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	components map[int64]MoveSystemData
	logger     ecs.IInternalLogger
}

func (p *MoveSystem) Init() {
	//p.logger = p.GetWorld().logger
	p.SetRequirements(&Position{}, &Movement{})
	p.components = map[int64]MoveSystemData{}
}

func (p *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta
	for _, comInfos := range p.components {
		position := comInfos.position
		move := comInfos.movement
		position.X = position.X + int(float64(move.dir[0]*move.v)*delta.Seconds())
		position.Y = position.Y + int(float64(move.dir[1]*move.v)*delta.Seconds())
		position.Z = position.Z + int(float64(move.dir[2]*move.v)*delta.Seconds())

		p.logger.Info("target id:", position.Owner().ID(), " current position:", position.X, position.Y, position.Z)
	}
}

//hp component
type HealthPoint struct {
	ecs.Component[HealthPoint]
	HP int
}

//damage system
type DamageSystemData struct {
	movement *Movement
	position *Position
}

type DamageSystem struct {
	ecs.System[DamageSystem]
	components map[int64]DamageSystemData
}

func (p *DamageSystem) Init() {
	p.SetRequirements(Position{}, HealthPoint{})
	p.components = map[int64]DamageSystemData{}
}

// Update will be called every frame
func (p *DamageSystem) Update(event ecs.Event) {

}

//main function
func Runtime0() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8888", nil))
	//}()

	rt := ecs.Runtime
	rt.Run()

	world := rt.NewWorld()
	world.Run()

	world.Register(&MoveSystem{})

	ee1 := world.NewEntity()
	ee2 := world.NewEntity()
	ee3 := world.NewEntity()

	println(ee1.ID(), ee2.ID(), ee3.ID())

	p1 := &Position{}
	p1.X = 100
	p1.Y = 100
	p1.Z = 100
	m1 := &Movement{}
	m1.v = 1000
	m1.dir = []int{1, 0, 0}
	world.NewEntity().AddByTemplate(p1, m1)
	p2 := &Position{}
	p2.X = 100
	p2.Y = 100
	p2.Z = 100
	m2 := &Movement{}
	m2.v = 2000
	m2.dir = []int{0, 1, 0}
	e2 := world.NewEntity()
	e2.AddByTemplate(p2, m2)

	world.Info("e2:", e2.ID())

	for {
		time.Sleep(time.Second * 3)
	}
}

func Runtime1() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8888", nil))
	//}()

	//rt := Runtime
	//world := rt.NewWorld()
	//go rt.Run()
	//
	//world.Register(&MoveSystem{})
	//world.Register(&DamageSystem{})
	//
	////初始化实体
	//for i := 0; i < 100000; i++ {
	//	p1 := NewPosition()
	//	p1.X = 100
	//	p1.Y = 100
	//	p1.Z = 100
	//	m1 := NewMovement()
	//	m1.v = 1000
	//	m1.dir = []int{1, 0, 0}
	//	world.NewEntity().AddComponent(p1, m1)
	//	p2 := NewPosition()
	//	p2.X = 100
	//	p2.Y = 100
	//	p2.Z = 100
	//	m2 := NewMovement()
	//	m2.v = 2000
	//	m2.dir = []int{0, 1, 0}
	//	world.NewEntity().AddComponent(p2, m2)
	//}
	//time.Sleep(time.Second * 3)
	//<-make(chan struct{})
}

func main() {
	Runtime0()
}