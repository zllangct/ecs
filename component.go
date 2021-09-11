package ecs

import (
	"reflect"
	"sync"
	"unsafe"
)

type IComponent interface {
	Owner() *Entity
	Type() reflect.Type
	setID(id int64)
	setOwner(entity *Entity)
}

type ITComponent[T any] interface {
	IComponent
	NewContainer() *Collection
	ToIComponent() IComponent
}

type Component[T any] struct {
	lock     sync.Mutex
	owner    *Entity
	id 		 int64
	realType reflect.Type
}

func (c *Component[T]) setOwner(entity *Entity) {
	c.owner = entity
}

func (c *Component[T]) setID(id int64){
	c.id = id
}

func (c *Component[T]) Ins() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T]) Owner() *Entity {
	return c.owner
}

func (c *Component[T]) Type() reflect.Type {
	if c.realType == nil {
		c.realType = reflect.TypeOf(*(new(T)))
	}
	return c.realType
}

func (c *Component[T]) ToIComponent() IComponent {
	return IComponent(c)
}

func (c *Component[T]) NewContainer() *Collection {
	return NewCollection(int(reflect.TypeOf(*new(T)).Size()))
}
