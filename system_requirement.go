package ecs

import "reflect"

const (
	ComponentReadWrite ComponentPermission = 0
	ComponentReadOnly  ComponentPermission = 1
)

type ComponentPermission uint8

type IRequirement interface {
	Type() reflect.Type
	getPermission() ComponentPermission
}

type ReadOnly[T ComponentObject] struct{}

func (r *ReadOnly[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (r *ReadOnly[T]) getPermission() ComponentPermission {
	return ComponentReadOnly
}
