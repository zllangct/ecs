package ecs

import "time"

type Event struct {
	Frame uint64
	Delta time.Duration
}

type InitReceiver interface {
	Init()
}

type SyncBeforeStartReceiver interface {
	SyncBeforeStart(event Event)
}

type StartReceiver interface {
	Start(event Event)
}

type SyncAfterStartReceiver interface {
	SyncAfterStart(event Event)
}

type SyncBeforePreUpdateReceiver interface {
	SyncBeforePreUpdate(event Event)
}

type PreUpdateReceiver interface {
	PreUpdate(event Event)
}

type SyncAfterPreUpdateReceiver interface {
	SyncAfterPreUpdate(event Event)
}

type SyncBeforeUpdateReceiver interface {
	SyncBeforeUpdate(event Event)
}

type UpdateReceiver interface {
	Update(event Event)
}

type SyncAfterUpdateReceiver interface {
	SyncAfterUpdate(event Event)
}

type SyncBeforePostUpdateReceiver interface {
	SyncBeforePostUpdate(event Event)
}

type PostUpdateReceiver interface {
	PostUpdate(event Event)
}

type SyncAfterPostUpdateReceiver interface {
	SyncAfterPostUpdate(event Event)
}

type SyncBeforeDestroyReceiver interface {
	SyncBeforeDestroy(event Event)
}

type DestroyReceiver interface {
	Destroy(event Event)
}

type SyncAfterPostDestroyReceiver interface {
	SyncAfterDestroy(event Event)
}
