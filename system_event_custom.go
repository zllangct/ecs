package ecs

type CustomEventName string
type CustomEventHandler func(*GateApi, []interface{})

type CustomEvent struct {
	Event CustomEventName
	Args  []interface{}
}
