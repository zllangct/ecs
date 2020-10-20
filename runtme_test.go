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

//movement component
type Movement struct {
	ComponentBase
	v   int
	dir []int
}

//move system
type MoveSystemData struct {
	movement *Movement
	position *Position
}

type MoveSystem struct {
	SystemBase
	components map[uint64]MoveSystemData
}

func (p *MoveSystem) Init(runtime *Runtime) {
	p.SystemBase.Init(runtime)
	p.SetType(reflect.TypeOf(p))
	p.SetRequirements(&Position{}, &Movement{})
	p.components = map[uint64]MoveSystemData{}
}

func (p *MoveSystem) Filter(com IComponent, op ComponentOperate) {
	if p.IsConcerned(com) {
		owner := com.GetOwner()
		switch op {
		case COMPONENT_OPERATE_ADD:
			p.components[owner.ID()] = MoveSystemData{
				movement: owner.GetComponent(&Movement{}).(*Movement),
				position: owner.GetComponent(&Position{}).(*Position),
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

		//println("current position:", position.X, position.Y, position.Z)
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
	p.SetRequirements(&Position{}, &HealthPoint{})
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
	rt.Register(&DamageSystem{})

	entity1 := NewEntity(rt)
	entity1.AddComponent(
		&Position{X: 100, Y: 100, Z: 100},
		&Movement{
			v:   1000,
			dir: []int{1, 0, 0},
		},
	)
	entity2 := NewEntity(rt)
	entity2.AddComponent(
		&Position{X: 100, Y: 100, Z: 100},
		&Movement{
			v:   2000,
			dir: []int{0, 1, 0},
		},
	)
	time.Sleep(time.Second * 20)
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
		NewEntity(rt).AddComponent(
			&Position{X: 100, Y: 100, Z: 100},
			&Movement{
				v:   1000,
				dir: []int{1, 0, 0},
			},
		)
		NewEntity(rt).AddComponent(
			&Position{X: 100, Y: 100, Z: 100},
			&Movement{
				v:   2000,
				dir: []int{0, 1, 0},
			},
		)
	}
	//time.Sleep(time.Second * 30)
	<-make(chan struct{})
}
