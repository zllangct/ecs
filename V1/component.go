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
	GetBase() *ComponentBase
	Init()
}

type ComponentBase struct {
	sync.Mutex
	owner *Entity
	state ComponentState
}

func (p *ComponentBase)Init()  {

}

func (p *ComponentBase) GetBase() *ComponentBase {
	return p
}

func (p *ComponentBase) GetState() ComponentState {
	p.Lock()
	defer p.Unlock()
	return p.state
}

func (p *ComponentBase) IsAlive() bool {
	p.Lock()
	defer p.Unlock()
	return p.state != COMPONENT_STATE_CLOSED
}

