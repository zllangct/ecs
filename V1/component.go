package main

import "sync"

const (
	COMPONENT_STATE_NONE = iota
	COMPONENT_STATE_NORMAL
	COMPONENT_STATE_CLOSED
)

type ComponentState = int

type IComponent interface {

}

type Reference struct {
	Sys ISystem
	Index *int
}

type ComponentBase struct {
	sync.RWMutex
	owner IEntity
	state ComponentState
}