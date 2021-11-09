package ecs

import (
	"fmt"
	"reflect"
	"unsafe"
)

type IComponent interface {
	Owner() *EntityInfo
	Type() reflect.Type
	ID() int64
	Template() IComponent

	setOwner(owner *EntityInfo)
	setID(id int64)
	setState(state ComponentState)
	getState() ComponentState
	getComponentType() ComponentType
	//init()

	instance() IComponent
	newCollection() interface{}
	addToCollection(collection interface{}) IComponent
	deleteFromCollection(collection interface{})
}

type IComponentTemplate interface {
	b26897f1d0938a734939da7d2f8aaedc()
}

type IFreeComponentTemplate interface {
	sf8sg6s9s0s0df5g55dg8sd7f()
	toIComponent(com interface{}) IComponent
}

const (
	h4   = uint8(240)
	l4   = uint8(15)
	zero = uint8(0)
)

type ComponentState uint8

const (
	ComponentStateInvalid ComponentState = iota
	ComponentStateActive
	ComponentStateDisable
)

type ComponentType uint8

const (
	ComponentTypeNormal ComponentType = iota
	ComponentTypeDisposable
	ComponentTypeFree
	ComponentTypeFreeDisposable
)

type FreeComponent[T any] struct {
	Component[T]
}

func (f *FreeComponent[T]) getComponentType() ComponentType {
	return ComponentTypeFree
}

func (f FreeComponent[T]) sf8sg6s9s0s0df5g55dg8sd7f() {}

func (f FreeComponent[T]) toIComponent(com interface{}) IComponent {
	return com.(IComponent)
}

type DisposableComponent[T any] struct {
	Component[T]
}

func (f *DisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeDisposable
}

type FreeDisposableComponent[T any] struct {
	Component[T]
}

func (f *FreeDisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeFreeDisposable
}

func (f FreeDisposableComponent[T]) toIComponent(com interface{}) IComponent {
	return com.(IComponent)
}

func (f *FreeDisposableComponent[T]) sf8sg6s9s0s0df5g55dg8sd7f() {}

type Component[T any] struct {
	owner    *EntityInfo
	id       int64
	realType reflect.Type
	st       uint8
}

func (c Component[T]) b26897f1d0938a734939da7d2f8aaedc() {}

func (c *Component[T]) init() {
	c.setType(c.getComponentType())
	c.setState(ComponentStateInvalid)
}

func (c *Component[T]) getComponentType() ComponentType {
	return ComponentTypeNormal
}

func (c *Component[T]) addToCollection(collection interface{}) IComponent {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return nil
	}
	id, ins := cc.Add(c.rawInstance())
	com := c.toIComponent(ins)
	com.setID(id)
	com.setState(ComponentStateActive)
	return com
}

func (c *Component[T]) toIComponent(com interface{}) IComponent {
	return com.(IComponent)
}

func (c *Component[T]) deleteFromCollection(collection interface{}) {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return
	}
	cc.Remove(c.ID())
	c.setState(ComponentStateDisable)
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

func (c *Component[T]) rawInstance() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T]) instance() (com IComponent) {
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
	c.setState(ComponentStateDisable)
}

func (c *Component[T]) Active() {
	c.setState(ComponentStateActive)
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

func (c *Component[T]) ToString() string {
	return fmt.Sprintf("%+v", c.rawInstance())
}
