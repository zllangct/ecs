package ecs

import "time"

type SyncWorld struct {
	worldBase
}

func NewSyncWorld(config *WorldConfig) *SyncWorld {
	w := &SyncWorld{}
	w.worldBase.init(config)
	return w
}

func (w *SyncWorld) Startup() {
	w.startup()
}

func (w *SyncWorld) Update() {
	w.update()
}

func (w *SyncWorld) Optimize(t time.Duration, force bool) {}

func (w *SyncWorld) Stop() {
	w.stop()
}

func (w *SyncWorld) NewEntity() *EntityInfo {
	return w.newEntity()
}

func (w *SyncWorld) GetEntityInfo(id Entity) (*EntityInfo, bool) {
	return w.getEntityInfo(id)
}

func (w *SyncWorld) getWorld() iWorldBase {
	return w
}
