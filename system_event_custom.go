package ecs

type SysEventHandler func(...interface{})

type SystemCustomEvent struct {
	Event string
	Args  []interface{}
}
