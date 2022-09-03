package ecs

import (
	"reflect"
	"runtime"
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
	Debug             bool //Debug模式
	IsMetrics         bool
	IsMetricsPrint    bool
	CpuNum            int    //使用的最大cpu数量
	MaxPoolThread     uint32 //线程池最大线程数量
	MaxPoolJobQueue   uint32 //线程池最大任务队列长度
	HashCount         int    //容器桶数量
	CollectionVersion int
	FrameInterval     time.Duration //帧间隔
	StopCallback      func(world *ecsWorld)
}

func NewDefaultWorldConfig() *WorldConfig {
	return &WorldConfig{
		Debug:           true,
		IsMetrics:       true,
		IsMetricsPrint:  false,
		CpuNum:          runtime.NumCPU(),
		MaxPoolThread:   uint32(runtime.NumCPU() * 2),
		MaxPoolJobQueue: 10,
		HashCount:       runtime.NumCPU() * 4,
		FrameInterval:   time.Millisecond * 33,
	}
}

type IWorld interface {
	GetStatus() WorldStatus
	GetID() int64
	NewEntity() EntityInfo
	GetEntityInfo(id Entity) EntityInfo
	AddFreeComponent(component IComponent)
	RegisterSystem(system ISystem)
	GetSyncLauncher() *SyncWorldLauncher
	GetAsyncLauncher() *AsyncWorldLauncher
	Optimize(t time.Duration, force bool) // TODO move to launcher

	getSystem(sys reflect.Type) (ISystem, bool)
	update()
	setStatus(status WorldStatus)
	addComponent(entity Entity, component IComponent)
	getComponentSet(typ reflect.Type) IComponentSet
	getComponentSetByIntType(typ uint16) IComponentSet
	getComponentCollection() IComponentCollection
	registerForT(system interface{}, order ...Order)
}

type ecsWorld struct {
	//id
	id int64
	//world status
	status WorldStatus
	//config
	config *WorldConfig
	//system flow,all systems
	systemFlow *systemFlow
	//all components
	components IComponentCollection
	//all entities
	entities *EntitySet
	//optimizer
	optimizer *optimizer
	//entity id generator
	idGenerator *EntityIDGenerator

	gate IGate

	workPool *Pool
	metrics  *Metrics

	frame           uint64
	ts              time.Time
	delta           time.Duration
	pureUpdateDelta time.Duration
}

func NewWorld(config *WorldConfig) *ecsWorld {
	world := &ecsWorld{
		id:         LocalUniqueID(),
		systemFlow: nil,
		config:     config,
		entities:   NewEntityCollection(),
		status:     WorldStatusInit,
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

	world.idGenerator = NewEntityIDGenerator(1024, 10)

	world.metrics = NewMetrics(world.config.IsMetrics, world.config.IsMetricsPrint)

	world.components = NewComponentCollection(world, config.HashCount)
	world.optimizer = newOptimizer(world)

	if world.config.FrameInterval <= 0 {
		world.config.FrameInterval = time.Millisecond * 33
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

func (w *ecsWorld) update() {
	if w.status != WorldStatusRunning {
		w.status = WorldStatusRunning
	}
	e := Event{Delta: w.delta, Frame: w.frame}
	now := time.Now()
	w.systemFlow.run(e)
	w.frame++
	w.delta = time.Since(w.ts)
	w.pureUpdateDelta = time.Since(now)
	w.ts = time.Now()
}

func (w *ecsWorld) GetSyncLauncher() *SyncWorldLauncher {
	return newSyncWorldLauncher(w)
}

func (w *ecsWorld) GetAsyncLauncher() *AsyncWorldLauncher {
	return newAsyncWorldLauncher(w)
}

func (w *ecsWorld) Optimize(t time.Duration, force bool) {
	w.optimizer.optimize(t, force)
}

func (w *ecsWorld) setStatus(status WorldStatus) {
	w.status = status
}

func (w *ecsWorld) GetStatus() WorldStatus {
	return w.status
}

func (w *ecsWorld) GetMetrics() *Metrics {
	return w.metrics
}

func (w *ecsWorld) RegisterSystem(system ISystem) {
	w.systemFlow.register(system)
}

func (w *ecsWorld) registerForT(system interface{}, order ...Order) {
	sys := system.(ISystem)
	if len(order) > 0 {
		sys.setOrder(order[0])
	}
	w.RegisterSystem(system.(ISystem))
}

func (w *ecsWorld) getSystem(sys reflect.Type) (ISystem, bool) {
	s, ok := w.systemFlow.systems[sys]
	if ok {
		return s.(ISystem), ok
	}
	return nil, ok
}

func (w *ecsWorld) addJob(job func(), hashKey ...uint32) {
	w.workPool.Add(job, hashKey...)
}

func (w *ecsWorld) addEntity(entity Entity) {
	w.entities.Add(entity, nil)
}

func (w *ecsWorld) GetEntityInfo(entity Entity) EntityInfo {
	return EntityInfo{world: w, entity: entity}
}

func (w *ecsWorld) deleteEntity(entity Entity) {
	w.entities.Remove(entity)
}

func (w *ecsWorld) getComponentSet(typ reflect.Type) IComponentSet {
	return w.components.getComponentSet(typ)
}

func (w *ecsWorld) getComponentSetByIntType(typ uint16) IComponentSet {
	return w.components.getComponentSetByIntType(typ)
}

func (w *ecsWorld) getComponentCollection() IComponentCollection {
	return w.components
}

func (w *ecsWorld) NewEntity() EntityInfo {
	info := EntityInfo{world: w, entity: w.idGenerator.NewID()}
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

func IGateToInstance[T GateObject](gate any) *T {
	g, ok := gate.(*T)
	if ok {
		return nil
	}
	return g
}
