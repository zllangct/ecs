package game

import "github.com/zllangct/ecs"

type EmptySystem struct {
	ecs.System[EmptySystem]
}

func (e *EmptySystem) Init() {
	ecs.Log.Info("empty system init")
}

func (e *EmptySystem) Start(event ecs.Event) {
	ecs.Log.Info("empty system start")
}

func (e *EmptySystem) PreUpdate(event ecs.Event) {
	ecs.Log.Info("empty system pre update")
}

func (e *EmptySystem) Update(event ecs.Event) {
	ecs.Log.Info("empty system update")
}

func (e *EmptySystem) PostUpdate(event ecs.Event) {
	ecs.Log.Info("empty system post update")
}
