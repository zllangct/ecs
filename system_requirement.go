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
	checkSet(initializer *SystemInitializer) IComponentSet
}

type ReadOnly[T ComponentObject] struct{}

func (r *ReadOnly[T]) Type() reflect.Type {
	return TypeOf[T]()
}

func (r *ReadOnly[T]) getPermission() ComponentPermission {
	return ComponentReadOnly
}

func (r *ReadOnly[T]) checkSet(initializer *SystemInitializer) IComponentSet {
	ins := new(Component[T])
	return initializer.sys.World().getComponentCollection().checkSet(ins)
}
