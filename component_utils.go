package ecs

import (
	"encoding/json"
	"reflect"
	"unsafe"
)

func GetType[T ComponentObject]() reflect.Type {
	return TypeOf[T]()
}

func getComponentOwnerEntity(p unsafe.Pointer) Entity {
	return *(*Entity)(p)
}

func ObjectToString(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
