package game

import "github.com/zllangct/ecs"

type PlayerPosition struct {
	SessionID int
	Pos       Position
}

type SyncSystem struct {
	ecs.System[SyncSystem]
}

func (m *SyncSystem) Init(si ecs.SystemInitConstraint) {
	m.SetRequirements(si, &Position{}, &PlayerComponent{})
}

func (m *SyncSystem) PostUpdate(event ecs.Event) {
	p := ecs.GetComponentAll[Position](m)
	for i := p.Begin(); !p.End(); i = p.Next() {
		pc := ecs.GetRelated[PlayerComponent](m, i.Owner())
		if pc == nil {
			continue
		}
		SendToClient(pc.SessionID, PlayerPosition{
			SessionID: pc.SessionID,
			Pos:       *i,
		})
	}
}
