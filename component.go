package ecs

import (
	"reflect"
	"sync"
)

type IComponent interface {
	setOwner(*Entity)
	GetOwner() *Entity
	GetBase() *ComponentBase
	GetRealType() reflect.Type
}

type IComponentType interface {
	typeFlag()
}

type ComponentBase struct {
	lock     sync.Mutex
	owner    *Entity
	realType reflect.Type
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

func (p *ComponentBase) SetRealType(t reflect.Type) {
	p.realType = t
}

func (p *ComponentBase) GetRealType() reflect.Type {
	return p.realType
}

func (p ComponentBase) typeFlag() {
}
