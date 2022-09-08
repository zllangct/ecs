package game

import "github.com/zllangct/ecs"

type EmptySystem struct {
	ecs.System[EmptySystem]
	isPreStart  bool
	isStart     bool
	isPostStart bool
}

func (e *EmptySystem) Init(si ecs.SystemInitializer) {
	ecs.Log.Info("empty system init")
}

func (e *EmptySystem) Start(event ecs.Event) {
	ecs.Log.Info("empty system start")
}

func (e *EmptySystem) PreUpdate(event ecs.Event) {
	if e.isPreStart {
		return
	}
	e.isPreStart = true
	ecs.Log.Info("empty system pre update")
}

func (e *EmptySystem) Update(event ecs.Event) {
	if e.isStart {
		return
	}
	e.isStart = true
	ecs.Log.Info("empty system update")
}

func (e *EmptySystem) PostUpdate(event ecs.Event) {
	if e.isPostStart {
		return
	}
	e.isPostStart = true
	ecs.Log.Info("empty system post update")
}
