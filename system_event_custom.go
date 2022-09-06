package ecs

type CustomEventName string
type CustomEventHandler func(UtilityGetter, []interface{})

type CustomEvent struct {
	Event CustomEventName
	Args  []interface{}
}
