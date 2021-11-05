package game

import "github.com/zllangct/ecs"

type SyncSystem struct {
	ecs.System[SyncSystem]
	logger     ecs.ILogger
}

func (m *SyncSystem) Init() {
	m.SetRequirements(&Position{})
}

func (m *SyncSystem) Update(event ecs.Event) {

}