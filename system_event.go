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

func implsCheck(system any) uint16 {
	var impls uint16
	imp := false
	_, imp = system.(SyncBeforeStartReceiver)
	if imp {
		impls = impls | 1<<StageSyncBeforeStart
	}
	_, imp = system.(StartReceiver)
	if imp {
		impls = impls | 1<<StageStart
	}
	_, imp = system.(SyncAfterStartReceiver)
	if imp {
		impls = impls | 1<<StageSyncAfterStart
	}
	_, imp = system.(SyncBeforePreUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncBeforePreUpdate
	}
	_, imp = system.(PreUpdateReceiver)
	if imp {
		impls = impls | 1<<StagePreUpdate
	}
	_, imp = system.(SyncAfterPreUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncAfterPreUpdate
	}
	_, imp = system.(SyncBeforeUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncBeforeUpdate
	}
	_, imp = system.(UpdateReceiver)
	if imp {
		impls = impls | 1<<StageUpdate
	}
	_, imp = system.(SyncAfterUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncAfterUpdate
	}
	_, imp = system.(SyncBeforePostUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncBeforePostUpdate
	}
	_, imp = system.(PostUpdateReceiver)
	if imp {
		impls = impls | 1<<StagePostUpdate
	}
	_, imp = system.(SyncAfterPostUpdateReceiver)
	if imp {
		impls = impls | 1<<StageSyncAfterPostUpdate
	}
	_, imp = system.(SyncBeforeDestroyReceiver)
	if imp {
		impls = impls | 1<<StageSyncBeforeDestroy
	}
	_, imp = system.(DestroyReceiver)
	if imp {
		impls = impls | 1<<StageDestroy
	}
	_, imp = system.(SyncAfterPostDestroyReceiver)
	if imp {
		impls = impls | 1<<StageSyncAfterDestroy
	}
	return impls
}
