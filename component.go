package ecs

import (
	"fmt"
	"reflect"
	"unsafe"
)

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

type IComponent interface {
	Owner() *EntityInfo
	Type() reflect.Type

	setOwner(owner *EntityInfo)
	setState(state ComponentState)
	getState() ComponentState
	getComponentType() ComponentType
	getPermission() ComponentPermission

	newCollection() ICollection
	addToCollection(collection interface{}) IComponent
	deleteFromCollection(collection interface{})
}

type ComponentObject interface {
	componentIdentification()
	getEntity() Entity
}

type FreeComponentObject interface {
	ComponentObject
	freeComponentIdentification()
}

type DisposableComponentObject interface {
	ComponentObject
	disposableComponentIdentification()
}

type FreeDisposableComponentObject interface {
	ComponentObject
	freeComponentIdentification()
	disposableComponentIdentification()
}

type ComponentPointer[T ComponentObject] interface {
	IComponent
	*T
}

type FreeComponentPointer[T FreeComponentObject] interface {
	IComponent
	*T
}

type DisposableComponentPointer[T DisposableComponentObject] interface {
	IComponent
	*T
}

type FreeDisposableComponentPointer[T FreeDisposableComponentObject] interface {
	IComponent
	*T
}

type FreeComponent[T FreeComponentObject] struct {
	Component[T]
}

func (f *FreeComponent[T]) getComponentType() ComponentType {
	return ComponentTypeFree
}

func (f FreeComponent[T]) freeComponentIdentification() {}

type DisposableComponent[T DisposableComponentObject] struct {
	Component[T]
}

func (f *DisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeDisposable
}

func (f DisposableComponent[T]) disposableComponentIdentification() {}

type FreeDisposableComponent[T FreeDisposableComponentObject] struct {
	Component[T]
}

func (f *FreeDisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeFreeDisposable
}

func (f FreeDisposableComponent[T]) freeComponentIdentification() {}

func (f FreeDisposableComponent[T]) disposableComponentIdentification() {}

type Component[T ComponentObject] struct {
	st    uint8
	o1    uint8
	o2    uint16
	seq   uint32
	owner *EntityInfo
}

func (c Component[T]) componentIdentification() {}

func (c Component[T]) getEntity() Entity {
	return c.owner.Entity()
}

func (c *Component[T]) init() {
	c.setType(c.getComponentType())
	c.setState(ComponentStateInvalid)
}

func (c *Component[T]) getComponentType() ComponentType {
	return ComponentTypeNormal
}

func conv[T any](in T, p *T) IComponent {
	return nil
}

func (c *Component[T]) addToCollection(collection interface{}) IComponent {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return nil
	}
	ins, _ := cc.Add(c.rawInstance(), int64(c.owner.entity))
	*c.rawInstance() = *ins
	i := interface{}(ins).(IComponent)
	i.setState(ComponentStateActive)
	return i
}

func (c *Component[T]) deleteFromCollection(collection interface{}) {
	cc, ok := collection.(*Collection[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return
	}
	c.setState(ComponentStateDisable)
	cc.Remove(int64(c.owner.Entity()))
	return
}

func (c *Component[T]) newCollection() ICollection {
	return NewCollection[T]()
}

func (c *Component[T]) setOwner(entity *EntityInfo) {
	c.owner = entity
}

func (c *Component[T]) rawInstance() *T {
	return (*T)(unsafe.Pointer(c))
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

func (c *Component[T]) invalidate() {
	c.setState(ComponentStateInvalid)
}

func (c *Component[T]) active() {
	c.setState(ComponentStateActive)
}

func (c *Component[T]) remove() {
	if c.owner == nil {
		return
	}
	c.owner.Remove(c)
}

func (c *Component[T]) Owner() *EntityInfo {
	return c.owner
}

func (c *Component[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (c *Component[T]) getPermission() ComponentPermission {
	return ComponentReadWrite
}

func (c *Component[T]) ToString() string {
	return fmt.Sprintf("%+v", c.rawInstance())
}
