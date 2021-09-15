package ecs

import (
	"reflect"
	"sync"
	"unsafe"
	"fmt"
)

type IComponent interface {
	Owner() *Entity
	Type() reflect.Type
	ID() int64
	Instance() IComponent
	Template() IComponent

	setOwner(owner *Entity)
	SetID(id int64)
	ComponentType() reflect.Type
	AddToCollection(collection interface{}) IComponent
	NewCollection() interface{}
}

type Component[T any] struct {
	lock     sync.Mutex
	owner    *Entity
	id 		 int64
	realType  reflect.Type
	operation map[string]func()[]interface{}
}

func (c *Component[T]) SetID(id int64)  {
	c.id = id
}

func (c *Component[T]) ComponentType() reflect.Type {
	return reflect.TypeOf(*new(T))
}

func (c *Component[T]) AddToCollection(collection interface{}) IComponent {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return nil
	}
	Log.Info("AddToCollection:", )
	_, ins := cc.Add(c.Ins())
	var com IComponent
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(ins)
	return com
}

func (c *Component[T]) NewCollection() interface{} {
	return NewCollection[T]()
}

func (c *Component[T]) setOwner(entity *Entity) {
	c.owner = entity
}

func (c *Component[T]) setID(id int64){
	c.id = id
}

func (c *Component[T]) ID() int64{
	return c.id
}

func (c *Component[T]) Ins() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T]) Instance() IComponent {
	var com IComponent
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(c)
	return com
}

func (c *Component[T]) Template() IComponent {
	return (c)
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

func (c *Component[T]) String() string {
	return fmt.Sprintf("%+v", c.Ins())
}




