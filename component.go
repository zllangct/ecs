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
	ID() int64

	setOwner(owner *EntityInfo)
	setState(state ComponentState)
	getState() ComponentState
	getComponentType() ComponentType

	instance() IComponent
	newCollection() interface{}
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

type FreeComponent[T FreeComponentObject, TP FreeComponentPointer[T]] struct {
	Component[T, TP]
}

func (f *FreeComponent[T, TP]) getComponentType() ComponentType {
	return ComponentTypeFree
}

func (f FreeComponent[T, TP]) freeComponentIdentification() {}

type DisposableComponent[T DisposableComponentObject, TP DisposableComponentPointer[T]] struct {
	Component[T, TP]
}

func (f *DisposableComponent[T, TP]) getComponentType() ComponentType {
	return ComponentTypeDisposable
}

func (f DisposableComponent[T, TP]) disposableComponentIdentification() {}

type FreeDisposableComponent[T FreeDisposableComponentObject, TP FreeDisposableComponentPointer[T]] struct {
	Component[T, TP]
}

func (f *FreeDisposableComponent[T, TP]) getComponentType() ComponentType {
	return ComponentTypeFreeDisposable
}

func (f FreeDisposableComponent[T, TP]) freeComponentIdentification() {}

func (f FreeDisposableComponent[T, TP]) disposableComponentIdentification() {}

type Component[T ComponentObject, TP ComponentPointer[T]] struct {
	id       int64
	st       uint8
	owner    *EntityInfo
	realType reflect.Type
}

func (c Component[T, TP]) componentIdentification() {}

func (c Component[T, TP]) getEntity() Entity {
	return Entity(c.id)
}

func (c *Component[T, TP]) init() {
	c.setType(c.getComponentType())
	c.setState(ComponentStateInvalid)
}

func (c *Component[T, TP]) getComponentType() ComponentType {
	return ComponentTypeNormal
}

func (c *Component[T, TP]) addToCollection(collection interface{}) IComponent {
	cc, ok := collection.(*Collection[T, TP])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return nil
	}
	_, ins := cc.Add(c.rawInstance())
	insP := TP(ins)
	insP.setState(ComponentStateActive)
	*c.rawInstance() = *insP
	return insP
}

func (c *Component[T, TP]) deleteFromCollection(collection interface{}) {
	cc, ok := collection.(*Collection[T, TP])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return
	}
	ret := cc.Remove(c.ID())
	if ret != nil {
		TP(ret).setState(ComponentStateDisable)
	}
	return
}

func (c *Component[T, TP]) newCollection() interface{} {
	return NewCollection[T, TP]()
}

func (c *Component[T, TP]) setOwner(entity *EntityInfo) {
	c.owner = entity
	c.id = int64(entity.Entity())
}

func (c *Component[T, TP]) ID() int64 {
	return c.id
}

func (c *Component[T, TP]) rawInstance() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T, TP]) instance() (com IComponent) {
	(*iface)(unsafe.Pointer(&com)).data = unsafe.Pointer(c)
	return
}

func (c *Component[T, TP]) setState(state ComponentState) {
	c.st = (c.st & l4) | (uint8(state) << 4)
}

func (c *Component[T, TP]) getState() ComponentState {
	return ComponentState(c.st & h4 >> 4)
}

func (c *Component[T, TP]) setType(typ ComponentType) {
	c.st = (c.st & h4) | uint8(typ)
}

func (c *Component[T, TP]) getType() ComponentType {
	return ComponentType(c.st & l4)
}

func (c *Component[T, TP]) invalidate() {
	c.setState(ComponentStateInvalid)
}

func (c *Component[T, TP]) active() {
	c.setState(ComponentStateActive)
}

func (c *Component[T, TP]) remove() {
	if c.owner == nil {
		return
	}
	c.owner.Remove(c)
}

func (c *Component[T, TP]) Owner() *EntityInfo {
	return c.owner
}

func (c *Component[T, TP]) Type() reflect.Type {
	if c.realType == nil {
		c.realType = TypeOf[T]()
	}
	return c.realType
}

func (c *Component[T, TP]) ToString() string {
	return fmt.Sprintf("%+v", c.rawInstance())
}
