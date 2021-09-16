package ecs

import (
	"reflect"
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
	ins := *new(T)
	c := entity.getComponentByType(reflect.TypeOf(ins))
	return (*T)(unsafe.Pointer((*iface)(unsafe.Pointer(&c)).data))
}
