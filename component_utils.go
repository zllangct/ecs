package ecs

import "reflect"

func GetType[T any]() reflect.Type {
	return reflect.TypeOf(*new(T))
}

