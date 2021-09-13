package main

import (
	"ecs"
	"reflect"
	"time"
)

//position component
type Position struct {
	ecs.Component[Position]
	X int
	Y int
	Z int
	v   int
	dir []int
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	logger     ecs.IInternalLogger
}

func (m *MoveSystem) Init() {
	//m.logger = m.GetWorld().logger
	m.SetRequirements(&Position{})
}

func (m *MoveSystem) Filter(ls map[reflect.Type][]ecs.ComponentOptResult) {
	if len(ls) > 0 {
		m.World().Info("new component:", len(ls))
	}
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta

	nc := m.GetInterestedNew()
	m.Filter(nc)

	//way 1:
	//cs := m.GetInterested(ecs.GetType[Position]()).(*ecs.Collection[Position])

	//way 2:
	cs := ecs.GetInterestedComponents[Position](m)

	for iter := ecs.NewIterator(cs); !iter.End(); iter.Next(){
		c := iter.Val()
		c.X = c.X + int(float64(c.dir[0]*c.v)*delta.Seconds())
		c.Y = c.Y + int(float64(c.dir[1]*c.v)*delta.Seconds())
		c.Z = c.Z + int(float64(c.dir[2]*c.v)*delta.Seconds())

		m.logger.Info("target id:", c.Owner().ID(), " current position:", c.X, c.Y, c.Z)
	}
}

//hp component
type HealthPoint struct {
	ecs.Component[HealthPoint]
	HP int
}

type Force struct {
	ecs.Component[Force]
	AttackRange int
	PhysicalBaseAttack int
	Strength int
	CriticalChange int
	CriticalMultiple int
}

type Action struct {
	ecs.Component[Action]
	Type int
}

type DamageSystem struct {
	ecs.System[DamageSystem]
	actions []ecs.ComponentOptResult
}

func (d *DamageSystem) Init() {
	d.SetRequirements(Position{}, HealthPoint{}, Force{}, Action{})
}


func (d *DamageSystem) Filter(ls map[reflect.Type][]ecs.ComponentOptResult) {
	if len(ls) == 0 {
		return
	}

	as, ok := ls[ecs.GetType[Action]()]
	if !ok {
		return
	}

	d.actions = as
}

// Update will be called every frame
func (d *DamageSystem) Update(event ecs.Event) {
	nc := d.GetInterestedNew()
	d.Filter(nc)

	//todo logic

}

//main function
func Runtime0() {
	// pprof
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8888", nil))
	//}()

	// 创建运行时
	rt := ecs.Runtime
	rt.Run()

	// 创建世界
	world := rt.NewWorld()
	world.Run()

	// 注册系统
	world.Register(&MoveSystem{})

	// 创建实体并添加组件
	ee1 := world.NewEntity()
	ee2 := world.NewEntity()
	ee3 := world.NewEntity()

	println(ee1.ID(), ee2.ID(), ee3.ID())

	p1 := &Position{
		X: 100,
		Y: 100,
		Z: 100,
		v: 2000,
		dir: []int{1,0,0},
	}
	world.NewEntity().AddByTemplate(p1)
	p2 := &Position{
		X: 100,
		Y: 100,
		Z: 100,
		v: 2000,
		dir: []int{0,1,0},
	}
	e2 := world.NewEntity()
	e2.AddByTemplate(p2)

	world.Info("e2:", e2.ID())

	for {
		time.Sleep(time.Second * 3)
	}
}

func main() {
	Runtime0()
}