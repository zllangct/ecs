package ecs

import "time"

type Event struct {
	Frame uint64
	Delta time.Duration
}

type InitReceiver interface {
	Init(ctx *SystemInitContext) error
}

type SyncBeforeStartReceiver interface {
	SyncBeforeStart(ctx *SystemContext, event Event) error
}

type StartReceiver interface {
	Start(ctx *SystemContext, event Event) error
}

type SyncAfterStartReceiver interface {
	SyncAfterStart(ctx *SystemContext, event Event) error
}

type SyncBeforePreUpdateReceiver interface {
	SyncBeforePreUpdate(ctx *SystemContext, event Event) error
}

type PreUpdateReceiver interface {
	PreUpdate(ctx *SystemContext, event Event) error
}

type SyncAfterPreUpdateReceiver interface {
	SyncAfterPreUpdate(ctx *SystemContext, event Event) error
}

type SyncBeforeUpdateReceiver interface {
	SyncBeforeUpdate(ctx *SystemContext, event Event) error
}

type UpdateReceiver interface {
	Update(ctx *SystemContext, event Event) error
}

type SyncAfterUpdateReceiver interface {
	SyncAfterUpdate(ctx *SystemContext, event Event) error
}

type SyncBeforePostUpdateReceiver interface {
	SyncBeforePostUpdate(ctx *SystemContext, event Event) error
}

type PostUpdateReceiver interface {
	PostUpdate(ctx *SystemContext, event Event) error
}

type SyncAfterPostUpdateReceiver interface {
	SyncAfterPostUpdate(ctx *SystemContext, event Event) error
}

type SyncBeforeDestroyReceiver interface {
	SyncBeforeDestroy(ctx *SystemContext, event Event) error
}

type DestroyReceiver interface {
	Destroy(ctx *SystemContext, event Event) error
}

type SyncAfterPostDestroyReceiver interface {
	SyncAfterDestroy(ctx *SystemContext, event Event) error
}

// SyncAfterDestroyReceiver 是 SyncAfterPostDestroyReceiver 的规范别名。
type SyncAfterDestroyReceiver = SyncAfterPostDestroyReceiver

type systemCallback func(*SystemContext, Event) error

// Bind stage receivers once when the system is registered.
func bindCallbacks(system any) (callbacks [StageMaxIndex + 1]systemCallback) {
	if r, ok := system.(SyncBeforeStartReceiver); ok {
		callbacks[StageSyncBeforeStart] = r.SyncBeforeStart
	}
	if r, ok := system.(StartReceiver); ok {
		callbacks[StageStart] = r.Start
	}
	if r, ok := system.(SyncAfterStartReceiver); ok {
		callbacks[StageSyncAfterStart] = r.SyncAfterStart
	}
	if r, ok := system.(SyncBeforePreUpdateReceiver); ok {
		callbacks[StageSyncBeforePreUpdate] = r.SyncBeforePreUpdate
	}
	if r, ok := system.(PreUpdateReceiver); ok {
		callbacks[StagePreUpdate] = r.PreUpdate
	}
	if r, ok := system.(SyncAfterPreUpdateReceiver); ok {
		callbacks[StageSyncAfterPreUpdate] = r.SyncAfterPreUpdate
	}
	if r, ok := system.(SyncBeforeUpdateReceiver); ok {
		callbacks[StageSyncBeforeUpdate] = r.SyncBeforeUpdate
	}
	if r, ok := system.(UpdateReceiver); ok {
		callbacks[StageUpdate] = r.Update
	}
	if r, ok := system.(SyncAfterUpdateReceiver); ok {
		callbacks[StageSyncAfterUpdate] = r.SyncAfterUpdate
	}
	if r, ok := system.(SyncBeforePostUpdateReceiver); ok {
		callbacks[StageSyncBeforePostUpdate] = r.SyncBeforePostUpdate
	}
	if r, ok := system.(PostUpdateReceiver); ok {
		callbacks[StagePostUpdate] = r.PostUpdate
	}
	if r, ok := system.(SyncAfterPostUpdateReceiver); ok {
		callbacks[StageSyncAfterPostUpdate] = r.SyncAfterPostUpdate
	}
	if r, ok := system.(SyncBeforeDestroyReceiver); ok {
		callbacks[StageSyncBeforeDestroy] = r.SyncBeforeDestroy
	}
	if r, ok := system.(DestroyReceiver); ok {
		callbacks[StageDestroy] = r.Destroy
	}
	if r, ok := system.(SyncAfterDestroyReceiver); ok {
		callbacks[StageSyncAfterDestroy] = r.SyncAfterDestroy
	}
	return callbacks
}
