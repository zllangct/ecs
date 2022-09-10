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
	ComponentTypeFreeMask       ComponentType = 1 << 7
	ComponentTypeDisposableMask ComponentType = 1 << 6
)

const (
	ComponentTypeNormal         ComponentType = 0
	ComponentTypeDisposable                   = 1 | ComponentTypeDisposableMask
	ComponentTypeFree                         = 2 | ComponentTypeFreeMask
	ComponentTypeFreeDisposable               = 3 | ComponentTypeFreeMask | ComponentTypeDisposableMask
)

type EmptyComponent struct {
	Component[EmptyComponent]
}

type IComponent interface {
	Owner() Entity
	Type() reflect.Type

	setOwner(owner Entity)
	setState(state ComponentState)
	setIntType(typ uint16)
	setSeq(seq uint32)
	getState() ComponentState
	getIntType() uint16
	getComponentType() ComponentType
	getPermission() ComponentPermission
	check(initializer SystemInitConstraint)
	getSeq() uint32
	newCollection(meta *ComponentMetaInfo) IComponentSet
	addToCollection(ct ComponentType, p unsafe.Pointer)
	deleteFromCollection(collection interface{})
	isValidComponentType() bool

	debugAddress() unsafe.Pointer
}

type ComponentObject interface {
	__ComponentIdentification()
	OwnerEntity() Entity
}

type FreeComponentObject interface {
	ComponentObject
	__FreeComponentIdentification()
}

type DisposableComponentObject interface {
	ComponentObject
	__DisposableComponentIdentification()
}

type FreeDisposableComponentObject interface {
	ComponentObject
	__FreeComponentIdentification()
	__DisposableComponentIdentification()
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

func (f FreeComponent[T]) __FreeComponentIdentification() {}

type DisposableComponent[T DisposableComponentObject] struct {
	Component[T]
}

func (f *DisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeDisposable
}

func (f DisposableComponent[T]) __DisposableComponentIdentification() {}

type FreeDisposableComponent[T FreeDisposableComponentObject] struct {
	Component[T]
}

func (f *FreeDisposableComponent[T]) getComponentType() ComponentType {
	return ComponentTypeFreeDisposable
}

func (f FreeDisposableComponent[T]) __FreeComponentIdentification() {}

func (f FreeDisposableComponent[T]) __DisposableComponentIdentification() {}

type Component[T ComponentObject] struct {
	st    uint8
	o1    uint8
	it    uint16
	seq   uint32
	owner Entity
}

func (c Component[T]) __ComponentIdentification() {}

func (c Component[T]) OwnerEntity() Entity {
	return c.owner
}

func (c *Component[T]) getComponentType() ComponentType {
	return ComponentTypeNormal
}

func (c *Component[T]) addToCollection(ct ComponentType, p unsafe.Pointer) {
	cc := (*ComponentSet[T])(p)
	var ins *T
	if ct&ComponentTypeFreeMask > 0 {
		ins, _ = cc.UnorderedCollection.Add(c.rawInstance())
	} else {
		ins = cc.Add(c.rawInstance(), c.owner)
	}

	if ins != nil {
		(*Component[T])(unsafe.Pointer(ins)).setState(ComponentStateActive)
	}
}

func (c *Component[T]) deleteFromCollection(collection interface{}) {
	cc, ok := collection.(*ComponentSet[T])
	if !ok {
		Log.Info("add to collection, collecion is nil")
		return
	}
	c.setState(ComponentStateDisable)
	cc.Remove(c.owner)
	return
}

func (c *Component[T]) newCollection(meta *ComponentMetaInfo) IComponentSet {
	return NewComponentSet[T](meta)
}

func (c *Component[T]) setOwner(entity Entity) {
	c.owner = entity
}

func (c *Component[T]) rawInstance() *T {
	return (*T)(unsafe.Pointer(c))
}

func (c *Component[T]) instance() IComponent {
	return any((*T)(unsafe.Pointer(c))).(IComponent)
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

func (c *Component[T]) setIntType(typ uint16) {
	c.it = typ
}

func (c *Component[T]) getIntType() uint16 {
	return c.it
}

func (c *Component[T]) setSeq(seq uint32) {
	c.seq = seq
}

func (c *Component[T]) getSeq() uint32 {
	return c.seq
}

func (c *Component[T]) invalidate() {
	c.setState(ComponentStateInvalid)
}

func (c *Component[T]) active() {
	c.setState(ComponentStateActive)
}

func (c *Component[T]) Owner() Entity {
	return c.owner
}

func (c *Component[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (c *Component[T]) getPermission() ComponentPermission {
	return ComponentReadWrite
}

func (c *Component[T]) check(initializer SystemInitConstraint) {
	if initializer.isValid() {
		panic("out of initialization stage")
	}
	ins := c.instance()
	if !ins.isValidComponentType() {
		panic("invalid component type")
	}
	sys := initializer.getSystem()
	sys.World().getOrCreateComponentMetaInfo(ins)
	sys.World().getComponentCollection().checkSet(ins)
}

func (c *Component[T]) isValidComponentType() bool {
	typ := c.Type()
	fieldNum := typ.NumField()
	if fieldNum < 1 {
		return false
	}
	for i := 1; i < fieldNum; i++ {
		subType := typ.Field(i).Type
		if !IsPureValueType(subType) {
			return false
		}
	}
	return true
}

func (c *Component[T]) debugAddress() unsafe.Pointer {
	return unsafe.Pointer(c)
}

func (c *Component[T]) ToString() string {
	return fmt.Sprintf("%+v", c.rawInstance())
}
