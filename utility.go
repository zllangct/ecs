package ecs

import (
	"reflect"
	"unsafe"
)

type UtilityGetter struct {
	world *iWorldBase
}

func (g *UtilityGetter) getWorld() iWorldBase {
	return *g.world
}

type UtilityObject interface {
	__UtilityIdentification()
}

type IUtility interface {
	Do(fn func(*SystemApi))
	getPointer() unsafe.Pointer
	setSystem(sys ISystem)
	setWorld(world iWorldBase)
}

type Utility[T SystemObject] struct {
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
	if u.sys == nil {
		s, ok := u.w.getSystem(TypeOf[T]())
		if ok {
			u.sys = s
		}
	}
	return u.sys
}

func (u *Utility[T]) Do(fn func(*SystemApi)) {
	u.sys.doSync(fn)
}

func (u *Utility[T]) getPointer() unsafe.Pointer {
	return unsafe.Pointer(u)
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
