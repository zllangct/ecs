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
	StopCallback       func(world *worldBase)
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

type iWorldBase interface {
	getStatus() WorldStatus
	getID() int64
	addFreeComponent(component IComponent)
	registerSystem(system ISystem)
	registerComponent(component IComponent)
	getMetrics() *Metrics
	getEntityInfo(id Entity) (*EntityInfo, bool)
	newEntity() *EntityInfo
	getComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo
	optimize(t time.Duration, force bool)
	getSystem(sys reflect.Type) (ISystem, bool)
	addUtility(utility IUtility)
	getUtilityForT(typ reflect.Type) (unsafe.Pointer, bool)
	update()
	setStatus(status WorldStatus)
	addComponent(entity Entity, component IComponent)
	getComponentSet(typ reflect.Type) IComponentSet
	getComponentSetByIntType(typ uint16) IComponentSet
	getComponentCollection() IComponentCollection
	getComponentMeta() *componentMeta
	getOrCreateComponentMetaInfo(component IComponent) *ComponentMetaInfo
	registerForT(system interface{}, order ...Order)
	base() *worldBase
}

type worldBase struct {
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

func (w *worldBase) init(config *WorldConfig) *worldBase {
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

func (w *worldBase) base() *worldBase {
	return w
}

func (w *worldBase) getID() int64 {
	return w.id
}

func (w *worldBase) SwitchMainThread() {
	atomic.StoreInt64(&w.mainThreadID, goroutineID())
}

func (w *worldBase) startup() {
	if w.getStatus() != WorldStatusInitialized {
		panic("world is not initialized or already running.")
	}

	if w.config.MetaInfoDebugPrint || w.config.Debug {
		w.systemFlow.SystemInfoPrint()
		w.componentMeta.ComponentMetaInfoPrint()
	}

	w.SwitchMainThread()
	w.workPool.Start()
	w.setStatus(WorldStatusRunning)
}

func (w *worldBase) update() {
	if w.config.MetaInfoDebugPrint {
		w.checkMainThread()
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

func (w *worldBase) optimize(t time.Duration, force bool) {
	w.optimizer.optimize(t, force)
}

func (w *worldBase) stop() {
	w.workPool.Release()
}

func (w *worldBase) setStatus(status WorldStatus) {
	w.status = status
}

func (w *worldBase) getUtilityGetter() UtilityGetter {
	ug := UtilityGetter{}
	iw := iWorldBase(w)
	ug.world = &iw
	return ug
}

func (w *worldBase) addUtility(utility IUtility) {
	w.utilities[utility.Type()] = utility
}
func (w *worldBase) getUtilityForT(typ reflect.Type) (unsafe.Pointer, bool) {
	u, ok := w.utilities[typ]
	return u.getPointer(), ok
}

func (w *worldBase) getStatus() WorldStatus {
	return w.status
}

func (w *worldBase) getMetrics() *Metrics {
	return w.metrics
}

func (w *worldBase) registerSystem(system ISystem) {
	if w.config.MainThreadCheck {
		w.checkMainThread()
	}
	w.systemFlow.register(system)
}

func (w *worldBase) registerComponent(component IComponent) {
	if w.config.MainThreadCheck {
		w.checkMainThread()
	}
	w.componentMeta.GetOrCreateComponentMetaInfo(component)
}

func (w *worldBase) registerForT(system interface{}, order ...Order) {
	sys := system.(ISystem)
	if len(order) > 0 {
		sys.setOrder(order[0])
	}
	w.registerSystem(system.(ISystem))
}

func (w *worldBase) getSystem(sys reflect.Type) (ISystem, bool) {
	s, ok := w.systemFlow.systems[sys]
	if ok {
		return s.(ISystem), ok
	}
	return nil, ok
}

func (w *worldBase) addJob(job func(), hashKey ...uint32) {
	w.workPool.Add(job, hashKey...)
}

func (w *worldBase) addEntity(info EntityInfo) *EntityInfo {
	return w.entities.Add(info)
}

func (w *worldBase) getEntityInfo(entity Entity) (*EntityInfo, bool) {
	return w.entities.GetEntityInfo(entity)
}

func (w *worldBase) deleteEntity(entity Entity) {
	w.entities.Remove(entity)
}

func (w *worldBase) getComponentSet(typ reflect.Type) IComponentSet {
	return w.components.getComponentSet(typ)
}

func (w *worldBase) getComponentSetByIntType(it uint16) IComponentSet {
	return w.components.getComponentSetByIntType(it)
}

func (w *worldBase) getComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo {
	return w.componentMeta.GetComponentMetaInfoByType(typ)
}

func (w *worldBase) getComponentCollection() IComponentCollection {
	return w.components
}

func (w *worldBase) getComponentMeta() *componentMeta {
	return w.componentMeta
}

func (w *worldBase) getOrCreateComponentMetaInfo(component IComponent) *ComponentMetaInfo {
	return w.componentMeta.GetOrCreateComponentMetaInfo(component)
}

func (w *worldBase) newEntity() *EntityInfo {
	info := EntityInfo{world: w, entity: w.idGenerator.NewID(), compound: NewCompound(4)}
	return w.addEntity(info)
}

func (w *worldBase) addComponent(entity Entity, component IComponent) {
	typ := component.Type()
	if !w.componentMeta.Exist(typ) {
		w.componentMeta.CreateComponentMetaInfo(component.Type(), component.getComponentType())
	}
	w.components.operate(CollectionOperateAdd, entity, component)
}

func (w *worldBase) deleteComponent(entity Entity, component IComponent) {
	w.components.operate(CollectionOperateDelete, entity, component)
}

func (w *worldBase) deleteComponentByIntType(entity Entity, it uint16) {
	w.components.deleteOperate(CollectionOperateDelete, entity, it)
}

func (w *worldBase) addFreeComponent(component IComponent) {
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
	default:
		Log.Errorf("component not free type, %s", component.Type().String())
		return
	}
	w.addComponent(0, component)
}

func (w *worldBase) checkMainThread() {
	if id := atomic.LoadInt64(&w.mainThreadID); id != goroutineID() && id > 0 {
		panic("not main thread")
	}
}
