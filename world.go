package ecs

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"
)

type WorldStatus int

const (
	WorldStatusInitializing WorldStatus = iota
	WorldStatusInitialized
	WorldStatusRunning
	WorldStatusStop
)

type WorldExecuteMode uint8

const (
	ExecuteModeInvalid WorldExecuteMode = iota
	ExecuteModeLinear
	ExecuteModeParallel
)

type WorldOption func(config *WorldConfig)

type WorldConfig struct {
	ExecuteMode  WorldExecuteMode
	AutoOptimize bool
	Rate         int
	// DisableGroups 关闭组表存储（布局优化层），关闭后语义完全等价（回退安全网/A-B 对比）
	DisableGroups bool
}

func (s *WorldConfig) initDefault() {
	s.ExecuteMode = ExecuteModeLinear
	s.AutoOptimize = false
	s.Rate = 30
}

func WithWorldSyncMode() WorldOption {
	return func(config *WorldConfig) {
		config.ExecuteMode = ExecuteModeLinear
	}
}

func WithWorldASyncMode() WorldOption {
	return func(config *WorldConfig) {
		config.ExecuteMode = ExecuteModeParallel
	}
}

func WithWorldAutoOptimize() WorldOption {
	return func(config *WorldConfig) {
		config.AutoOptimize = true
	}
}

func WithWorldDefaultUpdateRate(rate int) WorldOption {
	return func(config *WorldConfig) {
		config.Rate = rate
	}
}

// WithoutGroups 关闭组表存储：WithGroup 声明被忽略，所有组件留在 CSet。
func WithoutGroups() WorldOption {
	return func(config *WorldConfig) {
		config.DisableGroups = true
	}
}

type serializableWorld struct {
	idGenerator     *EntityIDGenerator
	entities        *EntitySet
	components      map[ComponentIntType]ComponentSet
	disposableTypes []ComponentIntType
	nomadicTypes    []ComponentIntType
	frame           uint64
}

type World struct {
	serializableWorld
	status     WorldStatus
	config     *WorldConfig
	optimizer  *optimizer
	opLog      *OpLog
	archetypes *ArchetypeRegistry
	// compVersion 组件池版本：components 新增池时递增，查询缓存据此失效
	compVersion int

	systems *flow

	pendingDestroy []Entity
	destroyMu      sync.Mutex

	lastUpdate time.Time
	delta      time.Duration
}

func NewWorld(opts ...WorldOption) *World {
	c := &WorldConfig{}
	c.initDefault()
	for _, opt := range opts {
		opt(c)
	}

	w := &World{}
	w.status = WorldStatusInitializing
	w.config = c
	w.idGenerator = NewEntityIDGenerator(1024, 10)
	w.entities = NewEntitySet()
	w.components = map[ComponentIntType]ComponentSet{}
	flowExecOpt := WithFlowSyncMode()
	if c.ExecuteMode == ExecuteModeParallel {
		flowExecOpt = WithFlowASyncMode()
	}
	w.systems = newSystemFlow(w, flowExecOpt)
	w.opLog = NewOpLog(w, runtime.NumCPU())
	w.optimizer = newOptimizer(w)
	w.archetypes = NewArchetypeRegistry()
	w.status = WorldStatusInitialized
	return w
}

// declareGroups 收集 system 声明的固定访问组（注册期调用）。
// 校验：组内组件必须已 RegisterComponent 且非 disposable/nomadic。
func (w *World) declareGroups(s *SystemInfoInstance) error {
	if w.config.DisableGroups {
		return nil
	}
	for _, g := range s.config.groups {
		sizes := make(map[ComponentIntType]int, len(g))
		for _, it := range g {
			info, ok := GetComponentRegistry().GetInfo(it)
			if !ok {
				return fmt.Errorf("ecs: WithGroup component type %d not registered", it)
			}
			proto := info.Proto()
			if proto.IsDisposable() || proto.IsNomadic() {
				return fmt.Errorf("ecs: WithGroup rejects disposable/nomadic component type %d", it)
			}
			sizes[it] = info.Size
		}
		if err := w.archetypes.Declare(g, sizes); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) NewEntity(opts ...EntityOption) Entity {
	return w.newEntity(opts...).entity
}

func (w *World) GetEntityInfo(entity Entity) (*EntityInfo, bool) {
	return w.entities.Get(entity)
}

func (w *World) newEntity(opts ...EntityOption) *EntityInfo {
	config := EntityConfig{}
	config.initDefault()
	for _, opt := range opts {
		opt(&config)
	}
	ins := w.entities.Add(EntityInfo{
		world:    w,
		entity:   w.idGenerator.NewID(),
		compound: NewCompound(4),
	})
	if len(config.comps) > 0 {
		ins.Add(config.comps...)
	}
	return ins
}

func (w *World) getComponentSet(it ComponentIntType) (ComponentSet, bool) {
	s, ok := w.components[it]
	return s, ok
}

// GetComponent 类型安全地直读实体组件，无需接口断言。
// 组件集不存在或实体无此组件时返回 (nil, false)。
// 适用于 system 之外的世界装配/检查代码，不经过 WithDep 依赖校验；
// system 执行期内应使用 (*SystemContext).GetBuddy。
func (w *World) GetComponent[T any, TP ComponentPointer[T]](e Entity) (*T, bool) {
	it := GetIntType[T, TP]()
	if a, ok := w.archetypes.owner(it); ok {
		if row, ok := a.rowOf(e.Index()); ok {
			col, _ := a.colOf(it)
			return (*T)(a.cellPtr(col, row)), true
		}
	}
	s, ok := w.getComponentSet(it)
	if !ok {
		return nil, false
	}
	b := s.get(e.Index())
	if b == nil {
		return nil, false
	}
	return (*T)(b), true
}

func (w *World) componentOp(op Operate) {
	w.opLog.operate(op)
}

func (w *World) AddNomadic(comps ...Component) {
	for _, comp := range comps {
		if !comp.IsNomadic() {
			continue
		}
		op := Operate{
			Op:   ComponentOperateAdd,
			Comp: comp,
		}
		w.componentOp(op)
	}
}

// DestroyEntity 销毁实体：延迟到帧同步点，回收 ID、清理所有组件与 compound
func (w *World) DestroyEntity(entity Entity) {
	w.destroyMu.Lock()
	w.pendingDestroy = append(w.pendingDestroy, entity)
	w.destroyMu.Unlock()
}

// flushPendingDestroy 帧同步点执行实体销毁
func (w *World) flushPendingDestroy() {
	w.destroyMu.Lock()
	list := w.pendingDestroy
	w.pendingDestroy = nil
	w.destroyMu.Unlock()

	for _, e := range list {
		// reuse 代数校验：忽略已被销毁/复用的旧句柄
		if _, ok := w.entities.Get(e); !ok {
			continue
		}
		for _, set := range w.components {
			set.Remove(e)
		}
		for _, a := range w.archetypes.groups {
			a.removeRow(e.Index())
		}
		w.entities.Remove(e)
		w.idGenerator.FreeID(e)
	}
}

func (w *World) Update() error {
	if w.status == WorldStatusStop {
		return errors.New("ecs: update on stopped world")
	}
	if w.status == WorldStatusInitialized {
		if err := w.archetypes.Finalize(); err != nil {
			return err
		}
		// Unmarshal 恢复的世界数据在 CSet 中，按 ownerOf 收拢入行
		w.absorbAll()
		w.status = WorldStatusRunning
	}
	now := time.Now()
	if w.lastUpdate.IsZero() {
		w.lastUpdate = now
	}

	w.frame++
	w.delta = now.Sub(w.lastUpdate)
	w.lastUpdate = now
	e := Event{
		Frame: w.frame,
		Delta: w.delta,
	}
	err := w.systems.Execute(e)
	if err != nil {
		return err
	}

	if w.config.AutoOptimize {
		elapsed := time.Since(now) // 本帧系统执行耗时
		t := (time.Second / time.Duration(w.config.Rate)) - elapsed
		err = w.Optimize(t, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *World) Destroy() error {
	if w.status == WorldStatusStop {
		return nil
	}
	// 所有系统进入 Destroy 状态，执行 Destroy stage 链
	for _, info := range w.systems.systems {
		info.setState(SystemStateDestroy)
	}
	e := Event{Frame: w.frame, Delta: w.delta}
	if err := w.systems.executeLinear(e); err != nil {
		return err
	}
	w.setStatus(WorldStatusStop)
	return nil
}

func (w *World) Optimize(t time.Duration, force bool) error {
	if w.optimizer == nil || t < 0 {
		return nil
	}
	w.optimizer.optimize(t, force)
	return nil
}

func (w *World) Register[T System, TP SystemPointer[T]](opts ...SystemOption) error {
	system := TP(new(T))
	s := newSystem(w, system, SystemTypeStandard)
	s.setState(SystemStateInit)

	ctx := &SystemInitContext{
		b: *s.getContext(),
	}
	ctx.b.constraint.activate()
	ctx.SetOption(opts...)
	if i, ok := any(system).(InitReceiver); ok {
		err := TryAndReport(func() error {
			return i.Init(ctx)
		})
		if err != nil {
			ctx.b.constraint.deactivate()
			return err
		}
	}
	ctx.b.constraint.deactivate()

	s.init(ctx.opts...)
	if err := w.declareGroups(s); err != nil {
		return err
	}
	w.systems.register(s)
	return nil
}

func (w *World) RegisterLight(system LightSystem, opts ...SystemOption) error {
	// LightSystem 的 Update 定义在指针接收者上，需传指针以实现 UpdateReceiver
	s := newSystem(w, &system, SystemTypeLight)
	s.setState(SystemStateInit)
	s.init(opts...)
	if err := w.declareGroups(s); err != nil {
		return err
	}
	w.systems.register(s)
	return nil
}

func (w *World) flushPendingOperate() ([]func(), []func(), func()) {
	return w.opLog.getOpTasks()
}

func (w *World) clearDisposable() error {
	for _, it := range w.disposableTypes {
		set, ok := w.getComponentSet(it)
		if !ok {
			continue
		}
		for _, index := range set.EntityIndexes() {
			e := w.entities.getByIndex(index)
			if e == nil {
				continue
			}
			e.compound.Remove(it)
		}

		set.Reset()
	}
	return nil
}

func (w *World) clearNomadic() error {
	for _, it := range w.nomadicTypes {
		set, ok := w.getComponentSet(it)
		if !ok {
			continue
		}
		set.Reset()
	}
	return nil
}

func (w *World) setStatus(status WorldStatus) {
	w.status = status
}

func (w *World) getStatus() WorldStatus {
	return w.status
}

func (w *World) getEntityInfo(entity Entity) (*EntityInfo, bool) {
	return w.entities.Get(entity)
}
