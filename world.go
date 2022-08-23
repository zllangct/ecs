package ecs

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

type WorldStatus int

const (
	WorldStatusInit WorldStatus = iota
	WorldStatusRunning
	WorldStatusPause
	WorldStatusStop
)

type WorldConfig struct {
	Debug                bool //Debug模式
	IsMetrics            bool
	IsMetricsPrint       bool
	CpuNum               int    //使用的最大cpu数量
	MaxPoolThread        uint32 //线程池最大线程数量
	MaxPoolJobQueue      uint32 //线程池最大任务队列长度
	HashCount            int    //容器桶数量
	CollectionVersion    int
	DefaultFrameInterval time.Duration //帧间隔
	StopCallback         func(world *ecsWorld)
}

func NewDefaultWorldConfig() *WorldConfig {
	return &WorldConfig{
		Debug:                true,
		IsMetrics:            true,
		IsMetricsPrint:       false,
		CpuNum:               runtime.NumCPU(),
		MaxPoolThread:        uint32(runtime.NumCPU() * 2),
		MaxPoolJobQueue:      10,
		HashCount:            runtime.NumCPU() * 4,
		DefaultFrameInterval: time.Millisecond * 33,
	}
}

type IWorld interface {
	Update()
	GetStatus() WorldStatus
	GetID() int64
	NewEntity() EntityInfo
	GetEntityInfo(id Entity) EntityInfo
	AddFreeComponent(component IComponent)
	Register(system ISystem)
	GetSystem(sys reflect.Type) (ISystem, bool)
	Optimize(t time.Duration, force bool)

	addComponent(entity Entity, component IComponent)
	getComponents(typ reflect.Type) IComponentSet
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
	//system flow,all systems
	systemFlow *systemFlow
	//all components
	//components *ComponentCollection
	components IComponentCollection
	//all entities
	entities *EntityCollection
	//sibling cache
	siblingCache *siblingCache
	//optimizer
	optimizer *optimizer

	workPool *Pool
	metrics  *Metrics

	frame           uint64
	ts              time.Time
	delta           time.Duration
	pureUpdateDelta time.Duration

	wStop chan struct{}
	//do some work for world cleaning
	stopHandler func(world *ecsWorld)
}

func NewWorld(config *WorldConfig) *ecsWorld {
	world := &ecsWorld{
		id:         LocalUniqueID(),
		systemFlow: nil,
		config:     config,
		entities:   NewEntityCollection(),
		status:     WorldStatusInit,
		wStop:      make(chan struct{}),
		ts:         time.Now(),
	}

	if world.config.MaxPoolThread <= 0 {
		world.config.MaxPoolThread = uint32(runtime.NumCPU())
	}

	if world.config.MaxPoolJobQueue <= 0 {
		world.config.MaxPoolJobQueue = 20
	}

	world.workPool = NewPool(config.MaxPoolThread, config.MaxPoolJobQueue)
	world.workPool.Start()

	world.components = NewComponentCollection(world, config.HashCount)
	world.optimizer = newOptimizer(world)
	world.siblingCache = newSiblingCache(world, 1024)

	if world.config.DefaultFrameInterval <= 0 {
		world.config.DefaultFrameInterval = time.Millisecond * 33
	}

	if world.config.HashCount == 0 {
		world.config.HashCount = config.CpuNum
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

func (w *ecsWorld) Update() {
	e := Event{Delta: w.delta, Frame: w.frame}
	now := time.Now()
	w.systemFlow.run(e)
	w.frame++
	w.delta = time.Since(w.ts)
	w.pureUpdateDelta = time.Since(now)
	w.ts = time.Now()
}

func (w *ecsWorld) Optimize(t time.Duration, force bool) {
	w.optimizer.optimize(t, force)
}

func (w *ecsWorld) GetStatus() WorldStatus {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.status
}

func (w *ecsWorld) GetMetrics() *Metrics {
	return w.metrics
}

func (w *ecsWorld) Destroy() {
	//TODO
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

func (w *ecsWorld) addJob(job func(), hashKey ...uint32) {
	w.workPool.Add(job, hashKey...)
}

// AddEntity entity operate : add
func (w *ecsWorld) addEntity(entity Entity) {
	w.entities.Add(entity)
}

func (w *ecsWorld) GetEntityInfo(entity Entity) EntityInfo {
	return EntityInfo{world: w, entity: entity}
}

// deleteEntity entity operate : delete
func (w *ecsWorld) deleteEntity(entity Entity) {
	w.entities.Remove(entity)
}

func (w *ecsWorld) getComponents(typ reflect.Type) IComponentSet {
	return w.components.getCollection(typ)
}

func (w *ecsWorld) NewEntity() EntityInfo {
	info := EntityInfo{world: w, entity: newEntity()}
	return info
}

func (w *ecsWorld) addComponent(entity Entity, component IComponent) {
	w.components.operate(CollectionOperateAdd, entity, component)
}

func (w *ecsWorld) deleteComponent(entity Entity, component IComponent) {
	w.components.operate(CollectionOperateDelete, entity, component)
}

func (w *ecsWorld) AddFreeComponent(component IComponent) {
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
	default:
		Log.Errorf("component not free type, %s", component.Type().String())
		return
	}
	w.addComponent(0, component)
}
