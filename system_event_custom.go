package ecs

type SystemCustomEventName string
type SysEventHandler func([]interface{})

type SystemCustomEvent struct {
	Event SystemCustomEventName
	Args  []interface{}
}
