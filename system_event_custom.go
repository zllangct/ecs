package ecs

type CustomEventName string
type CustomEventHandler func([]interface{})

type CustomEvent struct {
	Event CustomEventName
	Args  []interface{}
}
