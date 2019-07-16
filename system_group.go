package main

import "reflect"

type SystemGroup struct {
	Flags []reflect.Type
	systems []ISystem
}

func (this *SystemGroup)Add(system ISystem)  {

}

func (this *SystemGroup)Update()  {
	for _,system := range this.systems {
		system.Update()
	}
}