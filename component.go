package ecs

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

type IComponent interface {
	Owner() *Entity
	Type() reflect.Type
	ID() int64
	Ins() IComponent
	Template() IComponent

	setOwner(owner *Entity)
	setID(id int64)

	addToCollection(collection interface{}) IComponent
	deleteFromCollection(collection interface{})

	NewCollection() interface{}
}

type Component[T any] struct {
	lock      sync.Mutex
	owner     *Entity
	id        int64
	realType  reflect.Type
	operation map[string]func() []interface{}
}

func (c *Component[T]) addToCollection(collection interface{}) IComponent {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return nil
	}
	id, ins := cc.Add(c.RawIns())
	c.setID(id)
	var com IComponent
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(ins)
	return com
}

func (c *Component[T]) deleteFromCollection(collection interface{}) {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return
	}
	cc.Remove(c.ID())
	return
}

func (c *Component[T]) NewCollection() interface{} {
	return NewCollection[T]()
}

func (c *Component[T]) setOwner(entity *Entity) {
	c.owner = entity
}

func (c *Component[T]) setID(id int64) {
	c.id = id
}

func (c *Component[T]) ID() int64 {
	return c.id
}

func (c *Component[T]) RawIns() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T]) Ins() (com IComponent) {
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(c)
	return
}

func (c *Component[T]) Template() IComponent {
	return c
}

func (c *Component[T]) Owner() *Entity {
	return c.owner
}

func (c *Component[T]) Type() reflect.Type {
	if c.realType == nil {
		c.realType = TypeOf[T]()
	}
	return c.realType
}

func (c *Component[T]) String() string {
	return fmt.Sprintf("%+v", c.RawIns())
}
