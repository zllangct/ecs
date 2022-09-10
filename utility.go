package ecs

import (
	"reflect"
	"unsafe"
)

type IUtilityGetter interface {
	getWorld() iWorldBase
}

type UtilityGetter struct {
	world *iWorldBase
}

func (g UtilityGetter) getWorld() iWorldBase {
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
	Async(fn func(*SystemApi))
	Type() reflect.Type
	getPointer() unsafe.Pointer
	setSystem(sys ISystem)
	setWorld(world iWorldBase)
}

type Utility[T UtilityObject] struct {
	w   iWorldBase
	sys ISystem
}

func (u *Utility[T]) setWorld(w iWorldBase) {
	u.w = w
}

func (u *Utility[T]) setSystem(sys ISystem) {
	u.sys = sys
}

func (u Utility[T]) __UtilityIdentification() {}

func (u *Utility[T]) GetSystem() ISystem {
	return u.sys
}

func (u *Utility[T]) Async(fn func(*SystemApi)) {
	u.sys.doAsync(fn)
}

func (u *Utility[T]) Sync(fn func(*SystemApi)) {
	u.sys.doSync(fn)
}

func (u *Utility[T]) getPointer() unsafe.Pointer {
	return unsafe.Pointer(u)
}

func (u *Utility[T]) Type() reflect.Type {
	return TypeOf[T]()
}

type SystemApi struct {
	sys ISystem
}

func (s *SystemApi) GetRequirements() map[reflect.Type]IRequirement {
	return s.sys.GetRequirements()
}

func (s *SystemApi) Pause() {

}

func (s *SystemApi) Resume() {

}

func (s *SystemApi) Stop() {

}
