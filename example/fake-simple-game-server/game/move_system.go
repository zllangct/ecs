package game

import (
	"github.com/zllangct/ecs"
)

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	timeScale float64
}

func (m *MoveSystem) Init() {
	m.SetRequirements(&Position{}, &Movement{})
	m.EventRegister("UpdateTimeScale", m.UpdateTimeScale)
}

func (m *MoveSystem) UpdateTimeScale(timeScale []interface{}) {
	ecs.Log.Info("time scale change to ", timeScale[0])
	m.timeScale = timeScale[0].(float64)
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta

	iterPos := ecs.GetInterestedComponents[Position](m)
	if iterPos == nil {
		return
	}
	iterMov := ecs.GetInterestedComponents[Movement](m)
	if iterMov == nil {
		return
	}

	d := map[ecs.Entity]*MoveSystemData{}

	for iter := iterPos; !iter.End(); iter.Next() {
		position := iter.Val()
		owner := position.Owner()
		movement := ecs.CheckComponent[Movement](m, owner)
		if movement == nil {
			continue
		}

		d[position.Owner().Entity()] = &MoveSystemData{P: position, M: movement}
	}

	for e, data := range d {
		if data.M == nil || data.P == nil {
			continue
		}
		data.P.X = data.P.X + int(float64(data.M.Dir[0]*data.M.V)*delta.Seconds())
		data.P.Y = data.P.Y + int(float64(data.M.Dir[1]*data.M.V)*delta.Seconds())
		data.P.Z = data.P.Z + int(float64(data.M.Dir[2]*data.M.V)*delta.Seconds())

		ecs.Log.Info("target id:", e, "delta:", delta, " current position:", data.P.X, data.P.Y, data.P.Z)
	}

}
