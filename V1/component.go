package main

import "sync"

type ComponentState  int
const (
	COMPONENT_STATE_NONE ComponentState = iota
	COMPONENT_STATE_NORMAL
	COMPONENT_STATE_CLOSED
)


type IComponent interface {
	IsAlive() bool
	Init()
}

type ComponentBase struct {
	sync.RWMutex
	owner *Entity
	state ComponentState
}

func (p *ComponentBase)Init()  {

}

func (p *ComponentBase)IsAlive() bool {
	p.RLock()
	defer p.RUnlock()
	return p.state != COMPONENT_STATE_CLOSED
}

