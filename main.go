package main

import (
	"./ecs"
)

type NameComponent struct {
	ecs.ComponentBase
	Name string
}

type CountComponent struct {
	ecs.ComponentBase
	Count int
}

type AddressComponent struct {
	ecs.ComponentBase
	Address string
}

type NameSystem struct {
	ecs.SystemBase
}

func (this *NameSystem) Init() {

}

func (*NameSystem) ComponentRequire() []ecs.IComponent {
	return []ecs.IComponent{&NameComponent{}}
}

func (*NameSystem) Update() {
	panic("implement me")
}

type OtherSystem struct {
	ecs.SystemBase
}

func (*OtherSystem) Init() {

}

func (*OtherSystem) ComponentRequire() []ecs.IComponent {
	return []ecs.IComponent{&AddressComponent{},&CountComponent{}}
}

func (this *OtherSystem) Update() {
	//for _, com := range this.components {
	//	otherc:=()com
	//}
}

func main() {
	//runtime := Runtime{
	//	UpdateInterval:500,
	//}
	//
	//runtime.Filter(&NameSystem{})
	//runtime.Filter(&OtherSystem{})
	//
	//runtime.Run()

	e:=ecs.Entity{}

	for i := 0; i < 10; i++ {
		e.AddComponent(&NameComponent{})
	}

	m:=make(map[int]int)
	ecs.ConcurrentTest(func() {
		//e.Components()[0]=&NameComponent{Name:"1"}
		m[1]=1
	}, func() {
		//cs := e.Components()
		//cs = cs[1:]
		m[1]=2
	})

	println(ecs.NextUID())



	ch:=make(chan struct{})
	<-ch
}
