package ecs

import "reflect"

type IComponent interface {
	Type()reflect.Type
	Init(typ reflect.Type)
}

type ComponentBase struct {
	Entity *Entity
	typ reflect.Type
}

func (p *ComponentBase)Type() reflect.Type {
	return p.typ
}

func (p *ComponentBase)Init(typ reflect.Type)  {
	p.typ=typ
}



