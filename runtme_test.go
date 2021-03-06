package ecs

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"reflect"
	"testing"
	"time"
)

//position component
type Position struct {
	ComponentBase
	X int
	Y int
	Z int
}

func NewPosition() *Position {
	c := Position{}
	c.SetRealType(reflect.TypeOf(c))
	return &c
}

//movement component
type Movement struct {
	ComponentBase
	v   int
	dir []int
}

func NewMovement() *Movement {
	c := Movement{}
	c.SetRealType(reflect.TypeOf(c))
	return &c
}

//move system
type MoveSystemData struct {
	movement *Movement
	position *Position
}

type MoveSystem struct {
	SystemBase
	components map[uint64]MoveSystemData
	logger     IInternalLogger
}

func (p *MoveSystem) Init(runtime *Runtime) {
	p.SystemBase.Init(runtime)
	p.logger = runtime.logger
	p.SetType(reflect.TypeOf(p))
	p.SetRequirements(Position{}, Movement{})
	p.components = map[uint64]MoveSystemData{}
}

func (p *MoveSystem) Filter(com IComponent, op ComponentOperate) {
	if p.IsConcerned(com) {
		owner := com.GetOwner()
		switch op {
		case COMPONENT_OPERATE_ADD:
			p.components[owner.ID()] = MoveSystemData{
				movement: owner.GetComponent(Movement{}).(*Movement),
				position: owner.GetComponent(Position{}).(*Position),
			}
		case COMPONENT_OPERATE_DELETE:
			delete(p.components, owner.ID())
		}
	}
}

func (p *MoveSystem) Update(event Event) {
	delta := event.Delta
	for _, comInfos := range p.components {
		position := comInfos.position
		move := comInfos.movement
		position.X = position.X + int(float64(move.dir[0]*move.v)*delta.Seconds())
		position.Y = position.Y + int(float64(move.dir[1]*move.v)*delta.Seconds())
		position.Z = position.Z + int(float64(move.dir[2]*move.v)*delta.Seconds())

		p.logger.Info("target id:", position.GetOwner().ID(), " current position:", position.X, position.Y, position.Z)
	}
}

//hp component
type HealthPoint struct {
	ComponentBase
	HP int
}

//damage system
type DamageSystemData struct {
	movement *Movement
	position *Position
}

type DamageSystem struct {
	SystemBase
	components map[uint64]DamageSystemData
}

func (p *DamageSystem) Init(runtime *Runtime) {
	p.SystemBase.Init(runtime)
	p.SetType(reflect.TypeOf(p))
	p.SetRequirements(Position{}, HealthPoint{})
	p.components = map[uint64]DamageSystemData{}
}

func (p *DamageSystem) Filter(com IComponent, op ComponentOperate) {
	if p.IsConcerned(com) {
		owner := com.GetOwner()
		switch op {
		case COMPONENT_OPERATE_ADD:
			p.components[owner.ID()] = DamageSystemData{
				movement: owner.GetComponent(&HealthPoint{}).(*Movement),
				position: owner.GetComponent(&Position{}).(*Position),
			}
		case COMPONENT_OPERATE_DELETE:
			delete(p.components, owner.ID())
		}
	}
}

func (p *DamageSystem) Update(event Event) {

}

//main function
func TestRuntime0(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:8888", nil))
	}()

	rt := NewRuntime()
	go rt.Run()

	rt.Register(&MoveSystem{})
	//rt.Register(&DamageSystem{})

	NewEntity(rt)
	NewEntity(rt)
	NewEntity(rt)

	p1 := NewPosition()
	p1.X = 100
	p1.Y = 100
	p1.Z = 100
	m1 := NewMovement()
	m1.v = 1000
	m1.dir = []int{1, 0, 0}
	NewEntity(rt).AddComponent(p1, m1)
	p2 := NewPosition()
	p2.X = 100
	p2.Y = 100
	p2.Z = 100
	m2 := NewMovement()
	m2.v = 2000
	m2.dir = []int{0, 1, 0}
	NewEntity(rt).AddComponent(p2, m2)
	time.Sleep(time.Second * 3)
	//<-make(chan struct{})
}

func TestRuntime1(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe("localhost:8888", nil))
	}()

	rt := NewRuntime()
	go rt.Run()

	rt.Register(&MoveSystem{})
	rt.Register(&DamageSystem{})

	//初始化实体
	for i := 0; i < 100000; i++ {
		p1 := NewPosition()
		p1.X = 100
		p1.Y = 100
		p1.Z = 100
		m1 := NewMovement()
		m1.v = 1000
		m1.dir = []int{1, 0, 0}
		NewEntity(rt).AddComponent(p1, m1)
		p2 := NewPosition()
		p2.X = 100
		p2.Y = 100
		p2.Z = 100
		m2 := NewMovement()
		m2.v = 2000
		m2.dir = []int{0, 1, 0}
		NewEntity(rt).AddComponent(p2, m2)
	}
	time.Sleep(time.Second * 3)
	//<-make(chan struct{})
}
