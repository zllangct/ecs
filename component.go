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

type TComponent[T IComponent] interface {
	GetComponent() *T
}

type ComponentBase[T any] struct {
	lock     sync.Mutex
	owner    *Entity
	realType reflect.Type
}

func (c *ComponentBase) setOwner(entity *Entity) {
	c.owner = entity
}

func (c *ComponentBase) GetOwner() *Entity {
	return c.owner
}

func (c *ComponentBase) GetBase() *ComponentBase {
	return c
}

func (c *ComponentBase) SetRealType(t reflect.Type) {
	c.realType = t
}

func (c *ComponentBase) GetType() reflect.Type {
	return c.realType
}

func NewTContainer[T IComponent](){

}
