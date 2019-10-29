package main

import "reflect"

//position component
type Position struct {
	ComponentBase
	position []int

}

//movement component
type Movement struct {
	ComponentBase
	v int
	dir []int
}

//move system
type MoveSystem struct {
	SystemBase
}

func (p *MoveSystem)Init()  {
	p.SystemBase.Init()
	p.SetType(reflect.TypeOf(p))
	p.SetRequirements(&Position{},&Movement{})
}

func (p *MoveSystem) Filter() {
	panic("implement me")
}

func (p *MoveSystem) SystemUpdate() {

}

//hp component
type HealthPoint struct {
	ComponentBase
	HP int
}

//damage system
type DamageSystem struct {
	SystemBase
}

func (p *DamageSystem)Init()  {
	p.SystemBase.Init()
	p.SetType(reflect.TypeOf(p))
	p.SetRequirements(&Position{},&HealthPoint{})
}

func (p *DamageSystem) Filter() {
	panic("implement me")
}

func (p *DamageSystem) SystemUpdate() {
	panic("implement me")
}

//main function
func main()  {
	rt:=NewRuntime()
	rt.Run()

	rt.Register(&MoveSystem{})
	rt.Register(&DamageSystem{})

	entity:=NewEntity(rt)
	entity.AddComponent(
		&Position{
			position: []int{1,1,1},
		},
		&Movement{
			v:             0,
			dir:           []int{1,0,0},
		},
	)
}
