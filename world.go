package ecs

import (
	"reflect"
	"sync"
	"time"
)

type World struct {
	//mutex
	mutex sync.Mutex
	//frame interval
	frameInterval time.Duration
	//runtime
	runtime *ecsRuntime
	//system flow,all systems
	systemFlow *systemFlow
	//all components
	components *ComponentCollection
	//all entities
	entities *EntityCollection
	//logger
	logger IInternalLogger
}

func NewWorld(runtime *ecsRuntime) *World {
	//default config
	config := NewDefaultRuntimeConfig()
	world := &World{
		systemFlow: nil,
		components: NewComponentCollection(config.HashCount),
		entities:   NewEntityCollection(config.HashCount),
		logger:     runtime.logger,
	}
	//initialise system flow
	sf := newSystemFlow(world)
	world.systemFlow = sf
	//generate world
	return world
}

//start ecs world
func (w *World) Run() {
	//main loop
	frameInterval := w.frameInterval
	var ts time.Time
	var delta time.Duration
	for {
		ts = time.Now()
		w.systemFlow.run(delta)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			time.Sleep(frameInterval - delta)
			delta = frameInterval
		}
	}
}

func (w *World) AddJob(handler func(JobContext, ...interface{}), args ...interface{}) {
	w.runtime.AddJob(handler, args...)
}

//register system
func (w *World) Register(system ISystem) {
	w.systemFlow.register(system)
}

//entity operate : add
func (w *World) AddEntity(entity *Entity) {
	w.entities.add(entity)
}

//entity operate : delete
func (w *World) DeleteEntity(entity *Entity) {
	w.entities.delete(entity)
}

//entity operate : delete
func (w *World) DeleteEntityByID(id uint64) {
	w.entities.deleteByID(id)
}

func (w *World) ComponentAttach(target *Entity, com IComponent) {
	w.components.TempComponentOperate(target, com, COLLECTION_OPERATE_ADD)
}

func (w *World) ComponentRemove(target *Entity, com IComponent) {
	w.components.TempComponentOperate(target, com, COLLECTION_OPERATE_DELETE)
}

func (w *World) GetAllComponents() ComponentCollectionIter {
	return w.components.GetAllComponents()
}

func (w *World) Error(v ...interface{}) {
	if w.logger != nil {
		w.logger.Error(v...)
	}
}

func (w *World) getNewComponentsAll() []CollectionOperateInfo {
	return w.components.GetNewComponentsAll()
}

func (w *World) getNewComponents(op CollectionOperate, typ reflect.Type) []CollectionOperateInfo {
	return w.components.GetNewComponents(op, typ)
}

func (w *World) NewEntity() *Entity{
	return NewEntity(w)
}