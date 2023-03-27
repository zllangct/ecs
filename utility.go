package ecs

import (
	"reflect"
	"unsafe"
)

type IUtilityGetter interface {
	getWorld() IWorld
}

type UtilityGetter struct {
	world *IWorld
}

func (g UtilityGetter) getWorld() IWorld {
	return *g.world
}

type UtilityObject interface {
	__UtilityIdentification()
}

type UtilityPointer[T UtilityObject] interface {
	IUtility
	*T
}

type IUtility interface {
	Type() reflect.Type
	getPointer() unsafe.Pointer
	setSystem(sys ISystem)
	setWorld(world IWorld)
}

type utilityIdentification struct{}

func (u utilityIdentification) __UtilityIdentification() {}

type Utility[T UtilityObject] struct {
	utilityIdentification
	w   IWorld
	sys ISystem
}

func (u *Utility[T]) setWorld(w IWorld) {
	u.w = w
}

func (u *Utility[T]) setSystem(sys ISystem) {
	u.sys = sys
}

func (u *Utility[T]) GetSystem() ISystem {
	return u.sys
}

func (u *Utility[T]) getPointer() unsafe.Pointer {
	return unsafe.Pointer(u)
}

func (u *Utility[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (u *Utility[T]) GetRequirements() map[reflect.Type]IRequirement {
	return u.sys.GetRequirements()
}

func (u *Utility[T]) SystemPause() {
	u.sys.pause()
}

func (u *Utility[T]) SystemResume() {
	u.sys.resume()
}

func (u *Utility[T]) SystemStop() {
	u.sys.stop()
}
