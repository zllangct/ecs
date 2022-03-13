package game

import (
	"github.com/zllangct/ecs"
)

type InputSystem struct {
	ecs.System[InputSystem, *InputSystem]
}

func (is *InputSystem) Init() {
	// is.SetRequirements(&Movement{},&MoveChange{})
	ecs.AddRequireComponent2[Movement, MoveChange](is)
}

func (is *InputSystem) PreUpdate(event ecs.Event) {
	iterMC := ecs.GetInterestedComponents[MoveChange](is)
	var mov *Movement
	for mc := iterMC.Begin(); !iterMC.End(); mc = iterMC.Next() {
		info := ecs.GetEntityInfo(is.World(), mc.Entity)
		if info == nil {
			continue
		}
		mov = ecs.GetRelatedComponent[Movement](is, info)
		if mov != nil {
			ecs.Log.Infof("move changed: old: %+v, new: %+v", mov, mc)
			mov.V = mc.V
			mov.Dir = mc.Dir
		}
	}
}
