package ecs

import (
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
)

type WorldStatus int

const (
	WorldStatusInitializing WorldStatus = iota
	WorldStatusInitialized
	WorldStatusRunning
	WorldStatusStop
)

var mainThreadDebug = false
var mainThreadID int64 = -1

func EnableMainThreadDebug() {
	mainThreadDebug = true
}

func checkMainThread() {
	if id := atomic.LoadInt64(&mainThreadID); id != goroutineID() && id > 0 {
		panic("not main thread")
	}
}

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

type iWorldBase interface {
	GetStatus() WorldStatus
	GetID() int64
	NewEntity() *EntityInfo
	GetEntityInfo(id Entity) (*EntityInfo, bool)
	AddFreeComponent(component IComponent)
	RegisterSystem(system ISystem)
	RegisterComponent(component IComponent)
	GetMetrics() *Metrics
	GetUtilityGetter() UtilityGetter

	getComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo
	optimize(t time.Duration, force bool)
	getSystem(sys reflect.Type) (ISystem, bool)
	update()
	setStatus(status WorldStatus)
	addComponent(entity Entity, component IComponent)
	getComponentSet(typ reflect.Type) IComponentSet
	getComponentSetByIntType(typ uint16) IComponentSet
	getComponentCollection() IComponentCollection
	getComponentMeta() *componentMeta
	getOrCreateComponentMetaInfo(component IComponent) *ComponentMetaInfo
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

	componentMeta *componentMeta

	utilities map[reflect.Type]interface{}

	workPool *Pool
	metrics  *Metrics

	frame           uint64
	ts              time.Time
	delta           time.Duration
	pureUpdateDelta time.Duration
}

func (w *ecsWorld) init(config *WorldConfig) *ecsWorld {
	w.id = LocalUniqueID()
	w.systemFlow = nil
	w.config = config
	w.entities = NewEntityCollection()
	w.ts = time.Now()

	if w.config.MaxPoolThread <= 0 {
		w.config.MaxPoolThread = uint32(runtime.NumCPU())
	}

	if w.config.MaxPoolJobQueue <= 0 {
		w.config.MaxPoolJobQueue = 20
	}

	w.workPool = NewPool(config.MaxPoolThread, config.MaxPoolJobQueue)

	w.idGenerator = NewEntityIDGenerator(1024, 10)

	w.componentMeta = NewComponentMeta()
	w.utilities = make(map[reflect.Type]interface{})

	w.metrics = NewMetrics(w.config.IsMetrics, w.config.IsMetricsPrint)

	w.components = NewComponentCollection(w, config.HashCount)
	w.optimizer = newOptimizer(w)

	if w.config.FrameInterval <= 0 {
		w.config.FrameInterval = time.Millisecond * 33
	}

	if w.config.HashCount == 0 {
		w.config.HashCount = config.CpuNum
	}

	//initialise system flow
	sf := newSystemFlow(w)
	w.systemFlow = sf

	w.setStatus(WorldStatusInitialized)

	return w
}

func (w *ecsWorld) GetID() int64 {
	return w.id
}

func (w *ecsWorld) SwitchMainThread() {
	atomic.StoreInt64(&mainThreadID, goroutineID())
}

func (w *ecsWorld) startup() {
	if w.GetStatus() != WorldStatusInitialized {
		panic("world is not initialized or already running.")
	}

	// TODO 系统在初始化阶段已注册完毕，此处打印系统执行顺序和并行关系

	w.SwitchMainThread()
	w.workPool.Start()

	w.setStatus(WorldStatusRunning)
}

func (w *ecsWorld) update() {
	if mainThreadDebug {
		checkMainThread()
	}
	w.SwitchMainThread()
	if w.status != WorldStatusRunning {
		panic("world is not running, must startup first.")
	}
	e := Event{Delta: w.delta, Frame: w.frame}
	start := time.Now()
	w.systemFlow.run(e)
	now := time.Now()
	w.delta = now.Sub(w.ts)
	w.pureUpdateDelta = now.Sub(start)
	w.ts = now
	w.frame++
}

func (w *ecsWorld) optimize(t time.Duration, force bool) {
	w.optimizer.optimize(t, force)
}

func (w *ecsWorld) setStatus(status WorldStatus) {
	w.status = status
}

func (w *ecsWorld) GetUtilityGetter() UtilityGetter {
	ug := UtilityGetter{}
	iw := iWorldBase(w)
	ug.world = &iw
	return ug
}

func (w *ecsWorld) GetStatus() WorldStatus {
	return w.status
}

func (w *ecsWorld) GetMetrics() *Metrics {
	return w.metrics
}

func (w *ecsWorld) RegisterSystem(system ISystem) {
	if mainThreadDebug {
		checkMainThread()
	}
	w.systemFlow.register(system)
}

func (w *ecsWorld) RegisterComponent(component IComponent) {
	if mainThreadDebug {
		checkMainThread()
	}
	w.componentMeta.GetOrCreateComponentMetaInfo(component)
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

func (w *ecsWorld) addEntity(info EntityInfo) *EntityInfo {
	return w.entities.Add(info)
}

func (w *ecsWorld) GetEntityInfo(entity Entity) (*EntityInfo, bool) {
	return w.entities.GetEntityInfo(entity)
}

func (w *ecsWorld) deleteEntity(entity Entity) {
	w.entities.Remove(entity)
}

func (w *ecsWorld) getComponentSet(typ reflect.Type) IComponentSet {
	return w.components.getComponentSet(typ)
}

func (w *ecsWorld) getComponentSetByIntType(it uint16) IComponentSet {
	return w.components.getComponentSetByIntType(it)
}

func (w *ecsWorld) getComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo {
	return w.componentMeta.GetComponentMetaInfoByType(typ)
}

func (w *ecsWorld) getComponentCollection() IComponentCollection {
	return w.components
}

func (w *ecsWorld) getComponentMeta() *componentMeta {
	return w.componentMeta
}

func (w *ecsWorld) getOrCreateComponentMetaInfo(component IComponent) *ComponentMetaInfo {
	return w.componentMeta.GetOrCreateComponentMetaInfo(component)
}

func (w *ecsWorld) NewEntity() *EntityInfo {
	info := EntityInfo{world: w, entity: w.idGenerator.NewID()}
	return w.addEntity(info)
}

func (w *ecsWorld) addComponent(entity Entity, component IComponent) {
	typ := component.Type()
	if !w.componentMeta.Exist(typ) {
		w.componentMeta.CreateComponentMetaInfo(component.Type(), component.getComponentType())
	}
	w.components.operate(CollectionOperateAdd, entity, component)
}

func (w *ecsWorld) deleteComponent(entity Entity, component IComponent) {
	w.components.operate(CollectionOperateDelete, entity, component)
}

func (w *ecsWorld) deleteComponentByIntType(entity Entity, it uint16) {
	w.components.deleteOperate(CollectionOperateDelete, entity, it)
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

func (w *ecsWorld) ClearFreeComponentSetByType(typ reflect.Type) {
	free := w.componentMeta.GetFreeTypes()
	if it, ok := free[typ]; ok {
		w.components.deleteOperate(CollectionOperateDeleteAll, 0, it)
	}
}
