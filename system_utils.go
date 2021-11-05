package ecs

import (
	"unsafe"
)

func GetInterestedComponents[T any](s ISystem) *Collection[T] {
	typ := GetType[T]()
	if _, ok := s.Requirements()[typ]; !ok {
		Log.Error("not require, typ:", typ)
		return nil
	}
	if s.World() == nil {
		Log.Error("world is nil")
	}
	c := s.World().getComponents(typ)
	if c == nil {
		return nil
	}
	return c.(*Collection[T])
}

func CheckComponent[T any](s ISystem, entity *Entity) *T {
	return getComponentWithSystem[T](s, entity)
}

func getComponentWithSystem[T any](s ISystem, entity *Entity) *T {
	c := entity.getComponentByType(TypeOf[T]())
	return (*T)(unsafe.Pointer((*iface)(unsafe.Pointer(&c)).data))
}
