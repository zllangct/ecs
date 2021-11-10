package game

import (
	"github.com/zllangct/ecs"
)

type InputSystem struct {
	ecs.System[InputSystem]
}

func (is *InputSystem) Init() {
	// is.SetRequirements(&Movement{},&MoveChange{})
	ecs.AddRequireComponent2[Movement, MoveChange](is)
}

func (is *InputSystem) PreUpdate(event ecs.Event) {
	iterMC := ecs.GetInterestedComponents[MoveChange](is)
	if iterMC == nil {
		return
	}
	var mov *Movement
	for mc := iterMC.Begin(); !iterMC.End(); iterMC.Next() {
		mov = ecs.CheckComponent[Movement](is, mc.Owner())
		if mov != nil {
			ecs.Log.Infof("move changed: old: %+v, new: %+v", mov, mc)
			mov.V = mc.V
			mov.Dir = mc.Dir
		}
	}
}
