package ecs

import (
	"reflect"
	"sync"
)

type IComponent interface {
	setOwner(*Entity)
	setType(reflect.Type)
	GetOwner() *Entity
	GetBase() *ComponentBase
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

func (p *ComponentBase) setType(t reflect.Type) {
	p.typ = t
}

func CreateComponent(com interface{}) uint64 {
	//typ := reflect.TypeOf(com)
	return 0
}
