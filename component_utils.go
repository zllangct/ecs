package ecs

import (
	"encoding/json"
	"reflect"
)

func GetType[T any]() reflect.Type {
	return TypeOf[T]()
}

func ObjectToString(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
