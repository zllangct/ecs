package ecs

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

type WorldStatus int

type WorldConfig struct {
	HashCount            int           //容器桶数量
	DefaultFrameInterval time.Duration //帧间隔
}

func NewDefaultWorldConfig() *WorldConfig {
	return &WorldConfig{
		HashCount:            runtime.NumCPU() * 4,
		DefaultFrameInterval: time.Millisecond * 33,
	}
}

type World struct {
	//mutex
	mutex sync.Mutex
	//world status
	status WorldStatus
	//config
	config *WorldConfig
	//runtime
	runtime *ecsRuntime
	//system flow,all systems
	systemFlow *systemFlow
	systems    sync.Map
	//all components
	components *ComponentCollection
	//all entities
	entities *EntityCollection

	stop chan struct{}
	//do some work for world cleaning
	stopHandler func(world *World)
}

func NewWorld(runtime *ecsRuntime, config *WorldConfig) *World {
	world := &World{
		systemFlow: nil,
		config:     config,
		components: NewComponentCollection(config.HashCount),
		entities:   NewEntityCollection(config.HashCount),
		status:     StatusInit,
	}

	if world.config.DefaultFrameInterval <= 0 {
		world.config.DefaultFrameInterval = time.Millisecond * 33
	}

	if world.config.HashCount == 0 {
		world.config.HashCount = runtime.config.CpuNum
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
	if Runtime.Status() != StatusRunning {
		Log.Error("runtime is not running")
		return
	}

	w.mutex.Lock()
	if w.status != StatusInit {
		Log.Error("this world is already running.")
		return
	}
	frameInterval := w.config.DefaultFrameInterval
	w.status = StatusRunning
	w.mutex.Unlock()

	Log.Info("start world success")

	defer func() {
		w.mutex.Lock()
		w.status = StatusStop
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
				w.stopHandler(w)
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

func (w *World) Stop() {
	w.stop <- struct{}{}
}

func (w *World) SetStopHandler(handler func(world *World)) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.stopHandler = handler
}

func (w *World) GetStatus() WorldStatus {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.status
}

// Register register system
func (w *World) Register(system ISystem) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.systemFlow.register(system)
	w.systems.Store(reflect.TypeOf(system), system)
}

func (w *World) GetSystem(sys reflect.Type) (ISystem, bool) {
	s, ok := w.systems.Load(sys)
	return s.(ISystem), ok
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

func (w *World) ComponentAttach(target *Entity, com IComponent) {
	w.components.TempTemplateOperate(target, com.Template(), CollectionOperateAdd)
}

func (w *World) ComponentRemove(target *Entity, com IComponent) {
	w.components.TempTemplateOperate(target, com.Template(), CollectionOperateDelete)
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

func (w *World) NewEntity() *Entity {
	return NewEntity(w)
}

func GetSystem[T ISystem](w *World) (ISystem, bool) {
	return w.GetSystem(TypeOf[T]())
}
