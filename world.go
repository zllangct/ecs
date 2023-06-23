package ecs

import (
	"reflect"
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

type WorldStatus int

const (
	WorldStatusInitializing WorldStatus = iota
	WorldStatusInitialized
	WorldStatusRunning
	WorldStatusStop
)

type WorldConfig struct {
	Debug              bool //Debug模式
	MetaInfoDebugPrint bool
	MainThreadCheck    bool
	IsMetrics          bool
	IsMetricsPrint     bool
	CpuNum             int    //使用的最大cpu数量
	MaxPoolThread      uint32 //线程池最大线程数量
	MaxPoolJobQueue    uint32 //线程池最大任务队列长度
	HashCount          int    //容器桶数量
	CollectionVersion  int
	FrameInterval      time.Duration //帧间隔
	StopCallback       func(world *ecsWorld)
}

func NewDefaultWorldConfig() *WorldConfig {
	return &WorldConfig{
		Debug:              true,
		MetaInfoDebugPrint: true,
		MainThreadCheck:    true,
		IsMetrics:          true,
		IsMetricsPrint:     false,
		CpuNum:             runtime.NumCPU(),
		MaxPoolThread:      uint32(runtime.NumCPU() * 2),
		MaxPoolJobQueue:    10,
		HashCount:          runtime.NumCPU() * 4,
		FrameInterval:      time.Millisecond * 33,
	}
}

type IWorld interface {
	getStatus() WorldStatus
	getID() int64
	addFreeComponent(component IComponent)
	registerSystem(system ISystem)
	registerComponent(component IComponent)
	getMetrics() *Metrics
	getEntityInfo(id Entity) (*EntityInfo, bool)
	newEntity() *EntityInfo
	deleteEntity(entity Entity)
	getComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo
	optimize(t time.Duration, force bool)
	getSystem(sys reflect.Type) (ISystem, bool)
	addUtility(utility IUtility)
	getUtilityForT(typ reflect.Type) (unsafe.Pointer, bool)
	update()
	setStatus(status WorldStatus)
	addComponent(entity Entity, component IComponent)
	deleteComponent(entity Entity, component IComponent)
	deleteComponentByIntType(entity Entity, it uint16)
	getComponentSet(typ reflect.Type) IComponentSet
	getComponentSetByIntType(typ uint16) IComponentSet
	getComponentCollection() IComponentCollection
	getComponentMeta() *componentMeta
	getOrCreateComponentMetaInfo(component IComponent) *ComponentMetaInfo
	checkMainThread()
	base() *ecsWorld
}

type ecsWorld struct {
	id              int64
	status          WorldStatus
	config          *WorldConfig
	systemFlow      *systemFlow
	components      IComponentCollection
	entities        *EntitySet
	optimizer       *optimizer
	idGenerator     *EntityIDGenerator
	componentMeta   *componentMeta
	utilities       map[reflect.Type]IUtility
	workPool        *Pool
	metrics         *Metrics
	frame           uint64
	ts              time.Time
	delta           time.Duration
	pureUpdateDelta time.Duration
	mainThreadID    int64
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

	w.componentMeta = NewComponentMeta(w)
	w.utilities = make(map[reflect.Type]IUtility)

	w.metrics = NewMetrics(w.config.IsMetrics, w.config.IsMetricsPrint)

	w.components = NewComponentCollection(w, config.HashCount)
	w.optimizer = newOptimizer(w)

	if w.config.FrameInterval <= 0 {
		w.config.FrameInterval = time.Millisecond * 33
	}

	if w.config.HashCount == 0 {
		w.config.HashCount = config.CpuNum
	}

	sf := newSystemFlow(w)
	w.systemFlow = sf

	w.setStatus(WorldStatusInitialized)

	return w
}

func (w *ecsWorld) base() *ecsWorld {
	return w
}

func (w *ecsWorld) getID() int64 {
	return w.id
}

func (w *ecsWorld) switchMainThread() {
	atomic.StoreInt64(&w.mainThreadID, goroutineID())
}

func (w *ecsWorld) startup() {
	if w.getStatus() != WorldStatusInitialized {
		panic("world is not initialized or already running.")
	}

	if w.config.MetaInfoDebugPrint || w.config.Debug {
		w.systemFlow.SystemInfoPrint()
		w.componentMeta.ComponentMetaInfoPrint()
	}

	w.switchMainThread()
	w.workPool.Start()
	w.setStatus(WorldStatusRunning)
}

func (w *ecsWorld) update() {
	if w.config.MetaInfoDebugPrint {
		w.checkMainThread()
	}
	w.switchMainThread()
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

func (w *ecsWorld) stop() {
	w.workPool.Release()
}

func (w *ecsWorld) setStatus(status WorldStatus) {
	w.status = status
}

func (w *ecsWorld) getUtilityGetter() UtilityGetter {
	ug := UtilityGetter{}
	iw := IWorld(w)
	ug.world = &iw
	return ug
}

func (w *ecsWorld) addUtility(utility IUtility) {
	w.utilities[utility.Type()] = utility
}
func (w *ecsWorld) getUtilityForT(typ reflect.Type) (unsafe.Pointer, bool) {
	u, ok := w.utilities[typ]
	return u.getPointer(), ok
}

func (w *ecsWorld) getStatus() WorldStatus {
	return w.status
}

func (w *ecsWorld) getMetrics() *Metrics {
	return w.metrics
}

func (w *ecsWorld) registerSystem(system ISystem) {
	w.checkMainThread()
	w.systemFlow.register(system)
}

func (w *ecsWorld) registerComponent(component IComponent) {
	w.checkMainThread()
	w.componentMeta.GetOrCreateComponentMetaInfo(component)
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

func (w *ecsWorld) getEntityInfo(entity Entity) (*EntityInfo, bool) {
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

func (w *ecsWorld) newEntity() *EntityInfo {
	info := EntityInfo{entity: w.idGenerator.NewID(), compound: NewCompound(4)}
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

func (w *ecsWorld) addFreeComponent(component IComponent) {
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
	default:
		Log.Errorf("component not free type, %s", component.Type().String())
		return
	}
	w.addComponent(0, component)
}

func (w *ecsWorld) checkMainThread() {
	if !w.config.MainThreadCheck {
		return
	}
	if id := atomic.LoadInt64(&w.mainThreadID); id != goroutineID() && id > 0 {
		panic("not main thread")
	}
}
