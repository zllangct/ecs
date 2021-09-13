package ecs

import (
	"reflect"
	"sync"
	"time"
)

type WorldStatus int

type World struct {
	//mutex
	mutex sync.Mutex
	//world status
	status WorldStatus
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

	stop chan struct{}
	//do some work for world cleaning
	stopHandler func()
}

func NewWorld(runtime *ecsRuntime) *World {
	//default config
	config := NewDefaultRuntimeConfig()
	world := &World{
		systemFlow: nil,
		frameInterval: config.FrameInterval,
		components: NewComponentCollection(config.HashCount),
		entities:   NewEntityCollection(config.HashCount),
		logger:     runtime.logger,
		status: STATUS_INIT,
	}
	//initialise system flow
	sf := newSystemFlow(world)
	world.systemFlow = sf
	//generate world
	return world
}

// Run start ecs world
func (w *World) Run() {
	go w.run()
}

func (w *World) run() {
	if Runtime.Status() != STATUS_RUNNING {
		w.logger.Error("runtime is not running")
		return
	}

	w.mutex.Lock()
	if w.status != STATUS_INIT {
		w.logger.Info("this world is already running.")
		return
	}
	frameInterval := w.frameInterval
	w.status = STATUS_RUNNING
	w.mutex.Unlock()

	w.logger.Info("start world success")

	defer func() {
		w.mutex.Lock()
		w.status = STATUS_STOP
		w.mutex.Unlock()
	}()

	var ts time.Time
	var delta time.Duration
	//main loop
	for {
		select {
		case <-w.stop:
			w.mutex.Lock()
			if w.stopHandler != nil {
				w.stopHandler()
			}
			w.mutex.Unlock()
			return
		default:
		}

		ts = time.Now()
		w.systemFlow.run(delta)
		delta = time.Since(ts)
		//w.Info(delta, frameInterval - delta)
		if frameInterval-delta > 0 {
			time.Sleep(frameInterval - delta)
			delta = frameInterval
		}
	}
}

func (w *World) Stop()  {
	w.stop<- struct{}{}
}

func (w *World) SetStopHandler(handler func()){
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.stopHandler = handler
}

func (w *World) GetStatus() WorldStatus {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.status
}

func (w *World) AddJob(handler func(JobContext, ...interface{}), args ...interface{}) {
	Runtime.AddJob(handler, args...)
}

// Register register system
func (w *World) Register(system ISystem) {
	w.systemFlow.register(system)
}

// AddEntity entity operate : add
func (w *World) AddEntity(entity *Entity) {
	w.entities.add(entity)
}

// DeleteEntity entity operate : delete
func (w *World) DeleteEntity(entity *Entity) {
	w.entities.delete(entity)
}

// DeleteEntityByID entity operate : delete
func (w *World) DeleteEntityByID(id int64) {
	w.entities.deleteByID(id)
}

func (w *World) ComponentTemplateAttach(target *Entity, com IComponentTemplate) {
	w.components.TempTemplateOperate(target, com, COLLECTION_OPERATE_ADD)
}

func (w *World) ComponentAttach(target *Entity, com IComponent) {
	w.components.TempTemplateOperate(target, com.Template(), COLLECTION_OPERATE_ADD)
}

func (w *World) ComponentRemove(target *Entity, com IComponent) {
	w.components.TempTemplateOperate(target, com.Template(), COLLECTION_OPERATE_DELETE)
}

func (w *World) Error(v ...interface{}) {
	if w.logger != nil {
		w.logger.Error(v...)
	}
}

func (w *World) Info(v ...interface{}) {
	if w.logger != nil {
		w.logger.Info(v...)
	}
}

func (w *World) Fatal(v ...interface{}) {
	if w.logger != nil {
		w.logger.Fatal(v...)
	}
}

func (w *World) getNewComponentsAll() map[reflect.Type][]ComponentOptResult {
	return w.components.GetNewComponentsAll()
}

func (w *World) getNewComponents(typ reflect.Type) []ComponentOptResult {
	return w.components.GetNewComponents(typ)
}

func (w *World) getComponents(typ reflect.Type) interface{} {
	return w.components.GetCollection(typ)
}

func (w *World) NewEntity() *Entity{
	return NewEntity(w)
}
