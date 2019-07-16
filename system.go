package main

import (
	"reflect"
)

type ISystem interface {
	Init()
	ComponentRequire()[]IComponent
	Update()
}

type SystemBase struct {
	components map[reflect.Type]IComponent
}

func (this *SystemBase)AddComponent(component IComponent) IComponent {
	return this.components[reflect.TypeOf(component)]
}



