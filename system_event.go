package ecs

import "time"

type Event struct {
	Delta time.Duration
}

type IEventInit interface {
	Init()
}

type IEventStart interface {
	Start(event Event)
}

type IEventPostStart interface {
	PostStart(event Event)
}

type IEventUpdate interface {
	Update(event Event)
}

type IEventPostUpdate interface {
	PostUpdate(event Event)
}

type IEventDestroy interface {
	Destroy(event Event)
}

type IEventPostDestroy interface {
	PostDestroy(event Event)
}
