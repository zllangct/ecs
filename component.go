package ecs

import (
	"fmt"
	"reflect"
	"unsafe"
)

type IComponent interface {
	Owner() *EntityInfo //TODO 切换到 Entity
	Type() reflect.Type
	ID() int64
	Template() IComponent

	setOwner(owner *EntityInfo)
	setID(id int64)

	ins() IComponent
	newCollection() interface{}
	addToCollection(collection interface{}) IComponent
	deleteFromCollection(collection interface{})
}

const (
	h4 = uint8(240)
	l4 = uint8(15)
	zero = uint8(0)
)

type ComponentState uint8
const (
	ComponentInvalid ComponentState = iota
	ComponentActive
)

type ComponentType uint8
const (
	ComponentTypeNormal ComponentState = iota
	ComponentTypeOnce
	ComponentTypeFree
	ComponentTypeFreeAndOnce
)

type Component[T any] struct {
	owner     *EntityInfo
	id        int64
	realType  reflect.Type
	st     	  uint8
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

func (c *Component[T]) newCollection() interface{} {
	return NewCollection[T]()
}

func (c *Component[T]) setOwner(entity *EntityInfo) {
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

func (c *Component[T]) ins() (com IComponent) {
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(c)
	return
}

func (c *Component[T]) setState(state ComponentState) {
	c.st = (c.st & l4) | (uint8(state) << 4)
}

func (c *Component[T]) getState() ComponentState {
	return ComponentState(c.st & h4 >> 4)
}

func (c *Component[T]) setType(typ ComponentType) {
	c.st = (c.st & h4) | uint8(typ)
}

func (c *Component[T]) getType() ComponentType {
	return ComponentType(c.st & l4)
}

func (c *Component[T]) Invalidate() {
	c.setState(ComponentInvalid)
}

func (c *Component[T]) Active() {
	c.setState(ComponentActive)
}

func (c *Component[T]) Remove() {
	if c.owner == nil {
		return
	}
	c.owner.Remove(c)
}

func (c *Component[T]) Template() IComponent {
	return c
}

func (c *Component[T]) Owner() *EntityInfo {
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
