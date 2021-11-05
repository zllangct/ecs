package ecs

type SysEventHandler func(...interface{})

type SystemCustomEvent struct {
	Event string
	Args  []interface{}
}

type CustomEvent[T any] struct {
	Event string
	Args  T
}

type Event1[T1 any] TS1[T1]
type Event2[T1, T2 any] TS2[T1, T2]
type Event3[T1, T2, T3 any] TS3[T1, T2, T3]

func NewEvent1[T1 any](event string, a1 T1) {

}

func NewEvent2[T1, T2 any](event string, a1 T1, a2 T2) {

}

func NewEvent3[T1, T2, T3 any](event string, a1 T1, a2 T2, a3 T3) {

}
