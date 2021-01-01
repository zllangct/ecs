package ecs

import (
	"reflect"
	"sync"
)

type IComponent interface {
	setOwner(*Entity)
	GetOwner() *Entity
	GetBase() *ComponentBase
	GetType() reflect.Type
}

type IComponentType interface {
	typeFlag()
}

type ComponentBase struct {
	lock  sync.Mutex
	owner *Entity
	typ   reflect.Type
}

func (p *ComponentBase) setOwner(entity *Entity) {
	p.owner = entity
}

func (p *ComponentBase) GetOwner() *Entity {
	return p.owner
}

func (p *ComponentBase) GetBase() *ComponentBase {
	return p
}

func (p *ComponentBase) SetType(t reflect.Type) {
	p.typ = t
}

func (p *ComponentBase) GetType() reflect.Type {
	return p.typ
}

func (p ComponentBase) typeFlag() {
}
