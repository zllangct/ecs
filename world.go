package ecs

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

type WorldStatus int

type WorldConfig struct {
	HashCount            int //容器桶数量
	CollectionVersion    int
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
	GetEntityInfo(id Entity) *EntityInfo
	AddFreeComponent(component IComponent)
	Register(system ISystem)
	GetSystem(sys reflect.Type) (ISystem, bool)
	Optimize(t time.Duration)

	getComponents(typ reflect.Type) interface{}
	registerForT(system interface{}, order ...Order)
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
	//all components
	//components *ComponentCollection
	components IComponentCollection
	//all entities
	entities *EntityCollection
	//optimizer
	optimizer *optimizer

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
		wStop:      make(chan struct{}),
	}
	world.optimizer = newOptimizer(world)

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

func (w *ecsWorld) GetID() int64 {
	return w.id
}

// Run start ecs world
func (w *ecsWorld) Run() {
	go w.run()
}

func doFrameForBenchmark(w IWorld, frame uint64, lastDelta time.Duration) {
	world := w.(*ecsWorld)
	world.update(Event{Delta: lastDelta, Frame: frame})
}

func (w *ecsWorld) update(event Event) {
	w.systemFlow.run(event)
}

func (w *ecsWorld) Optimize(t time.Duration) {
	w.optimizer.optimize(t)
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
	var frame uint64
	//main loop
	for {
		select {
		case <-w.wStop:
			w.mutex.Lock()
			if w.stopHandler != nil {
				w.stopHandler(w)
			}
			w.systemFlow.stop()
			w.mutex.Unlock()
			return
		default:
		}

		ts = time.Now()
		w.update(Event{Delta: delta, Frame: frame})
		frame++
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
func (w *ecsWorld) Register(system ISystem) {
	w.systemFlow.register(system)
}

func (w *ecsWorld) registerForT(system interface{}, order ...Order) {
	sys := system.(ISystem)
	if len(order) > 0 {
		sys.setOrder(order[0])
	}
	w.Register(system.(ISystem))
}

func (w *ecsWorld) GetSystem(sys reflect.Type) (ISystem, bool) {
	s, ok := w.systemFlow.systems.Load(sys)
	if ok {
		return s.(ISystem), ok
	}
	return nil, ok
}

// AddEntity entity operate : add
func (w *ecsWorld) addEntity(entity *EntityInfo) {
	w.entities.add(entity)
}

func (w *ecsWorld) GetEntityInfo(id Entity) *EntityInfo {
	return w.entities.getInfo(id)
}

// deleteEntity entity operate : delete
func (w *ecsWorld) deleteEntity(info *EntityInfo) {
	w.entities.delete(info.entity)
}

func (w *ecsWorld) getComponents(typ reflect.Type) interface{} {
	return w.components.getCollection(typ)
}

func (w *ecsWorld) NewEntity() *EntityInfo {
	return newEntityInfo(w)
}

func (w *ecsWorld) addComponent(info *EntityInfo, component IComponent) {
	w.components.operate(CollectionOperateAdd, info, component)
}

func (w *ecsWorld) deleteComponent(info *EntityInfo, component IComponent) {
	w.components.operate(CollectionOperateDelete, info, component)
}

func (w *ecsWorld) AddFreeComponent(component IComponent) {
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
	default:
		Log.Errorf("component not free type, %s", component.Type().String())
		return
	}
	w.addComponent(nil, component)
}
