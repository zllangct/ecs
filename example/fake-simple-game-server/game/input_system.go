package game

import "github.com/zllangct/ecs"

type InputSystem struct {
	ecs.System[InputSystem]
}

func (m *InputSystem) Init() {
	m.SetRequirements(&MoveChange{}, &Movement{})
	m.EventRegister("Change", m.Change)
}

func (m *InputSystem) Change(in ...interface{}) {
	dir := in[0].([]int)
	v := in[1].(int)
	_, _ = dir, v

}

func (m *InputSystem) Update(event ecs.Event) {

}
