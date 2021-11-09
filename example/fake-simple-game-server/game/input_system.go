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
	var mov *Movement
	for mc := iterMC.Begin(); !iterMC.End(); iterMC.Next() {
		mov = ecs.CheckComponent[Movement](is, mc.Owner())
		mov.V = mc.V
		mov.Dir = mc.Dir
	}

	//TODO 验证一次性组件的删除是否正常
}
