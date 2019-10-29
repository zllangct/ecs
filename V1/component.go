package main

import "sync"

//state of component's life circle
type ComponentState  int
const (
	COMPONENT_STATE_NONE ComponentState = iota
	COMPONENT_STATE_NORMAL
	COMPONENT_STATE_CLOSING
	COMPONENT_STATE_CLOSED
)


type IComponent interface {
	GetOwner() *Entity
	GetBase() *ComponentBase
	setOwner(*Entity)
}


type ComponentBase struct {
	sync.Mutex
	owner *Entity
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




