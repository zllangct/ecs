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
	StopCallback         func(world *ecsWorld)
}

func NewDefaultWorldConfig() *WorldConfig {
	return &WorldConfig{
		HashCount:            runtime.NumCPU() * 4,
		DefaultFrameInterval: time.Millisecond * 33,
	}
}

type IWorld interface {
	Run()
	GetStatus() WorldStatus
	GetID() int64
	NewEntity() *EntityInfo

	register(system ISystem)
	registerForT(system interface{}, order ...Order)
	getComponents(typ reflect.Type) interface{}
	getNewComponents(typ reflect.Type) []OperateInfo
	getEntityInfo(id Entity) *EntityInfo
	getSystem(sys reflect.Type) (ISystem, bool)

	b14d94e462795b8bd42a0bf62ae90826()
}

type ecsWorld struct {
	//mutex
	mutex sync.Mutex
	//id
	id int64
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

	wStop chan struct{}
	//do some work for world cleaning
	stopHandler func(world *ecsWorld)
}

func newWorld(runtime *ecsRuntime, config *WorldConfig) *ecsWorld {
	world := &ecsWorld{
		id:         UniqueID(),
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

func (w *ecsWorld) b14d94e462795b8bd42a0bf62ae90826() {}

func (w *ecsWorld) GetID() int64 {
	return w.id
}

// Run start ecs world
func (w *ecsWorld) Run() {
	go w.run()
}

func (w *ecsWorld) run() {
	if Runtime.status() != StatusRunning {
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
		case <-w.wStop:
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

func (w *ecsWorld) stop() {
	w.wStop <- struct{}{}
}

func (w *ecsWorld) GetStatus() WorldStatus {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.status
}

// Register register system
func (w *ecsWorld) register(system ISystem) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.systemFlow.register(system)
	w.systems.Store(reflect.TypeOf(system), system)
}

func (w *ecsWorld) registerForT(system interface{}, order ...Order) {
	sys := system.(ISystem)
	if len(order) > 0 {
		sys.setOrder(order[0])
	}
	w.register(system.(ISystem))
}

func (w *ecsWorld) getSystem(sys reflect.Type) (ISystem, bool) {
	s, ok := w.systems.Load(sys)
	return s.(ISystem), ok
}

// AddEntity entity operate : add
func (w *ecsWorld) addEntity(entity *EntityInfo) {
	w.entities.add(entity)
}

func (w *ecsWorld) getEntityInfo(id Entity) *EntityInfo {
	return w.entities.getInfo(id)
}

// deleteEntity entity operate : delete
func (w *ecsWorld) deleteEntity(info *EntityInfo) {
	w.entities.delete(info.entity)
}

func (w *ecsWorld) getNewComponentsAll() map[reflect.Type][]OperateInfo {
	return w.components.GetNewComponentsAll()
}

func (w *ecsWorld) getNewComponents(typ reflect.Type) []OperateInfo {
	return w.components.GetNewComponents(typ)
}

func (w *ecsWorld) getComponents(typ reflect.Type) interface{} {
	return w.components.GetCollection(typ)
}

func (w *ecsWorld) NewEntity() *EntityInfo {
	return newEntityInfo(w)
}

func (w *ecsWorld) addComponent(info *EntityInfo, component IComponent) {
	w.components.TempTemplateOperate(info, component, CollectionOperateAdd)
}

func (w *ecsWorld) deleteComponent(info *EntityInfo, component IComponent) {
	w.components.TempTemplateOperate(info, component, CollectionOperateDelete)
}

func (w *ecsWorld) AddFreeComponent(component IComponent) {
	w.addComponent(nil, component)
}
