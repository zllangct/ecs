package main

import (
	"github.com/zllangct/ecs"
)

type MoveSystem struct {
	ecs.System[MoveSystem]
}

func (m *MoveSystem) Init(si ecs.SystemInitializer) {
	m.SetRequirements(si, &Position{}, &ecs.ReadOnly[Movement]{})
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta

	count := 0
	iter := ecs.GetInterestedComponents[Movement](m)
	for move := iter.Begin(); !iter.End(); move = iter.Next() {
		entity := move.Owner()
		pos := ecs.GetRelatedComponent[Position](m, entity)
		if pos == nil {
			continue
		}

		pos.X = pos.X + int(float64(move.Dir[0]*move.V)*delta.Seconds())
		pos.Y = pos.Y + int(float64(move.Dir[1]*move.V)*delta.Seconds())
		pos.Z = pos.Z + int(float64(move.Dir[2]*move.V)*delta.Seconds())

		oldMoveV := move.V
		_ = oldMoveV

		// test for read only component, changes are not effective
		move.V = move.V + 1

		count++
		//ecs.Log.Info("target id:", entity, " delta:", delta, " current position:", pos.X, pos.Y, pos.Z,
		//	" move speed:", oldMoveV)
	}
}
