package ecs

import (
	"reflect"
	"unsafe"
)

const (
	ComponentReadWrite ComponentPermission = 0
	ComponentReadOnly  ComponentPermission = 1
)

type ComponentPermission uint8

type IRequirement interface {
	Type() reflect.Type
	getPermission() ComponentPermission
	check(initializer SystemInitializer)
}

type ReadOnly[T ComponentObject] struct{}

func (r *ReadOnly[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (r *ReadOnly[T]) getPermission() ComponentPermission {
	return ComponentReadOnly
}

func (r *ReadOnly[T]) check(initializer SystemInitializer) {
	ins := any((*T)(unsafe.Pointer(r))).(IComponent)
	ins.check(initializer)
}
