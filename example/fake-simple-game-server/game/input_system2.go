package game

import (
	"github.com/zllangct/ecs"
)

type MoveChange2 struct {
	ecs.DisposableComponent[MoveChange2]
	Entity ecs.Entity
	V      int
	Dir    []int
}

type InputSystem2 struct {
	ecs.System[InputSystem2]
}

func (is *InputSystem2) Init(si ecs.SystemInitializer) {
	is.SetRequirements(si, &Movement{}, &MoveChange{})
}

func (is *InputSystem2) PreUpdate(event ecs.Event) {
	iterMC := ecs.GetInterestedComponents[MoveChange2](is)
	var mov *Movement
	for mc := iterMC.Begin(); !iterMC.End(); mc = iterMC.Next() {
		mov = ecs.GetRelatedComponent[Movement](is, mc.Owner())
		if mov != nil {
			ecs.Log.Infof("move changed: old: %+v, new: %+v", mov, mc)
			mov.V = mc.V
			mov.Dir = mc.Dir
		}
	}
}
