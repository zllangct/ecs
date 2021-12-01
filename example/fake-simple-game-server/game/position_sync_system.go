package game

import "github.com/zllangct/ecs"

type PlayerPosition struct {
	SessionID int
	Pos       Position
}

type SyncSystem struct {
	ecs.System[SyncSystem]
}

func (m *SyncSystem) Init() {
	ecs.AddRequireComponent2[Position, PlayerComponent](m)
}

func (m *SyncSystem) PostUpdate(event ecs.Event) {
	p := ecs.GetInterestedComponents[Position](m)
	if p == nil {
		return
	}

	for i := p.Begin(); !p.End(); i = p.Next() {
		pc := ecs.GetRelatedComponent[PlayerComponent](m, i.Owner())
		if pc == nil {
			continue
		}
		SendToClient(pc.SessionID, PlayerPosition{
			SessionID: pc.SessionID,
			Pos:       *i,
		})
	}
}
