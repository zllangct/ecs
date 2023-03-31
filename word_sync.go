package ecs

import "time"

type SyncWorld struct {
	ecsWorld
}

func NewSyncWorld(config *WorldConfig) *SyncWorld {
	w := &SyncWorld{}
	w.ecsWorld.init(config)
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

func (w *SyncWorld) NewEntity() Entity {
	return w.newEntity().Entity()
}

func (w *SyncWorld) DestroyEntity(entity Entity) {
	info, ok := w.getEntityInfo(entity)
	if !ok {
		return
	}
	info.Destroy(w)
}

func (w *SyncWorld) Add(entity Entity, components ...IComponent) {
	info, ok := w.getEntityInfo(entity)
	if !ok {
		return
	}
	info.Add(w, components...)
}

func (w *SyncWorld) Remove(entity Entity, components ...IComponent) {
	info, ok := w.getEntityInfo(entity)
	if !ok {
		return
	}
	info.Remove(w, components...)
}

func (w *SyncWorld) getWorld() IWorld {
	return w
}
