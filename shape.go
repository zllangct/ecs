package ecs

import (
	"reflect"
	"unsafe"
)

type ShapeObject interface {
	shapeBaseIdentifier()
}

type ShapeObjectPointer[T ShapeObject] interface {
	parse(info *EntityInfo, types []IRequirement) bool
	eleTypes() []reflect.Type
	*T
}

type ShapeBase struct{}

func (s ShapeBase) shapeBaseIdentifier() {}

type Shape2[T1, T2 ComponentObject] struct {
	ShapeBase
	C1 *T1
	C2 *T2
}

func (s *Shape2[T1, T2]) eleTypes() []reflect.Type {
	return []reflect.Type{TypeOf[T1](), TypeOf[T2]()}
}

func (s *Shape2[T1, T2]) parse(info *EntityInfo, types []IRequirement) bool {
	c1 := info.getComponentByTypeInSystem(types[0].Type())
	if c1 == nil {
		return false
	}
	s.C1 = (*T1)((*iface)(unsafe.Pointer(&c1)).data)

	c2 := info.getComponentByTypeInSystem(types[1].Type())
	if c2 == nil {
		return false
	}
	s.C2 = (*T2)((*iface)(unsafe.Pointer(&c2)).data)

	return true
}

type Shape3[T1, T2, T3 ComponentObject] struct {
	ShapeBase
	C1 *T1
	C2 *T2
	C3 *T3
}

func (s *Shape3[T1, T2, T3]) eleTypes() []reflect.Type {
	return []reflect.Type{TypeOf[T1](), TypeOf[T2](), TypeOf[T3]()}
}

func (s *Shape3[T1, T2, T3]) parse(info *EntityInfo, types []IRequirement) bool {
	c1 := info.getComponentByTypeInSystem(types[0].Type())
	if c1 == nil {
		return false
	}
	s.C1 = (*T1)((*iface)(unsafe.Pointer(&c1)).data)

	c2 := info.getComponentByTypeInSystem(types[1].Type())
	if c2 == nil {
		return false
	}
	s.C2 = (*T2)((*iface)(unsafe.Pointer(&c2)).data)

	c3 := info.getComponentByTypeInSystem(types[2].Type())
	if c3 == nil {
		return false
	}
	s.C3 = (*T3)((*iface)(unsafe.Pointer(&c3)).data)

	return true
}
