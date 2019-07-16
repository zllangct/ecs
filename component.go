package main

import "reflect"

type IComponent interface {
	Type()reflect.Type
	Init(typ reflect.Type)
}

type ComponentBase struct {
	Entity *Entity
	typ reflect.Type
}

func (this *ComponentBase)Type() reflect.Type {
	return this.typ
}

func (this *ComponentBase)Init(typ reflect.Type)  {
	this.typ=typ
}



