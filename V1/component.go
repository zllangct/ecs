package main

import "sync"

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




