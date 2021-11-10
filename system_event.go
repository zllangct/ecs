package ecs

import "time"

type Event struct {
	Frame uint64
	Delta time.Duration
}

type InitReceiver interface {
	Init()
}

type PreStartReceiver interface {
	PreStart(event Event)
}

type StartReceiver interface {
	Start(event Event)
}

type PostStartReceiver interface {
	PostStart(event Event)
}

type PreUpdateReceiver interface {
	PreUpdate(event Event)
}

type UpdateReceiver interface {
	Update(event Event)
}

type PostUpdateReceiver interface {
	PostUpdate(event Event)
}

type PreDestroyReceiver interface {
	PreDestroy(event Event)
}

type DestroyReceiver interface {
	Destroy(event Event)
}

type PostDestroyReceiver interface {
	PostDestroy(event Event)
}
