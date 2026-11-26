package ecs

import (
	"cmp"
	"errors"
	"fmt"
	"runtime"
	"slices"
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
// 两阶段原子校验：全部组通过校验（已注册、非 disposable/nomadic、无 Go 指针、去重后 >=2 类型）后才统一写入 archetypes。
func (w *World) declareGroups(s *SystemInfoInstance) error {
	if w.config.DisableGroups {
		return nil
	}
	type pendingGroup struct {
		types []ComponentIntType
		sizes map[ComponentIntType]int
	}
	validated := make([]pendingGroup, 0, len(s.config.groups))
	for _, g := range s.config.groups {
		seen := make(map[ComponentIntType]struct{}, len(g))
		unique := make([]ComponentIntType, 0, len(g))
		sizes := make(map[ComponentIntType]int, len(g))
		for _, it := range g {
			if _, dup := seen[it]; dup {
				continue
			}
			seen[it] = struct{}{}
			info, ok := GetComponentRegistry().GetInfo(it)
			if !ok {
				return fmt.Errorf("ecs: WithGroup component type %d not registered", it)
			}
			proto := info.Proto()
			if proto.IsDisposable() || proto.IsNomadic() {
				return fmt.Errorf("ecs: WithGroup rejects disposable/nomadic component type %d", it)
			}
			if info.HasPointers {
				return fmt.Errorf("ecs: WithGroup rejects component type %d containing Go pointers (unsafe for raw byte columns)", it)
			}
			if info.Size == 0 {
				return fmt.Errorf("ecs: WithGroup rejects zero-sized tag component type %d", it)
			}
			unique = append(unique, it)
			sizes[it] = info.Size
		}
		if len(unique) < 2 {
			return fmt.Errorf("ecs: group needs >= 2 distinct component types, got %v", g)
		}
		validated = append(validated, pendingGroup{types: unique, sizes: sizes})
	}
	for _, pg := range validated {
		if err := w.archetypes.Declare(pg.types, pg.sizes); err != nil {
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
// 实体不存在/句柄代数已过期、组件集不存在或实体无此组件时返回 (nil, false)。
// 适用于 system 之外的世界装配/检查代码，不经过 WithDep 依赖校验；
// system 执行期内应使用 (*SystemContext).GetBuddy。
func (w *World) GetComponent[T any, TP ComponentPointer[T]](e Entity) (*T, bool) {
	if _, ok := w.entities.Get(e); !ok {
		return nil, false
	}
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

// DestroyEntity queues destruction at the next frame's first synchronization
// point. Requests in that batch run in entity-index order, with the complete
// handle as a tiebreaker. Pending component operations are applied afterwards.
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

	if len(list) > 1 {
		slices.SortFunc(list, func(a, b Entity) int {
			if a.Index() != b.Index() {
				return cmp.Compare(a.Index(), b.Index())
			}
			return cmp.Compare(a, b)
		})
	}

	for _, e := range list {
		// reuse 代数校验：忽略已被销毁/复用的旧句柄
		info, ok := w.entities.Get(e)
		if !ok {
			continue
		}
		for _, it := range info.compound {
			// Every complete row owns the group's smallest type. Testing
			// that type once also works for entities in several groups.
			if a, ok := w.archetypes.owner(it); ok && it == a.types[0] {
				a.removeRow(e.Index())
			}
			if set := w.components[it]; set != nil {
				set.Remove(e)
			}
		}
		w.entities.SparseArray.Remove(e.Index())
		w.idGenerator.FreeID(e)
	}
}

func (w *World) Update() error {
	now := time.Now()
	if w.lastUpdate.IsZero() {
		w.lastUpdate = now
	}
	delta := now.Sub(w.lastUpdate)
	w.lastUpdate = now
	if err := w.Step(delta); err != nil {
		return err
	}
	if w.config.AutoOptimize {
		budget := time.Second/time.Duration(w.config.Rate) - time.Since(now)
		return w.Optimize(budget, false)
	}
	return nil
}

// Step 使用调用方显式传入的 delta 推进一帧，不按墙钟余量重排存储。
// WithWorldAutoOptimize 的帧预算仅用于 Update。
func (w *World) Step(delta time.Duration) error {
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

	w.frame++
	w.delta = delta
	e := Event{
		Frame: w.frame,
		Delta: w.delta,
	}
	return w.systems.Execute(e)
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
	// 释放组件集合（含 arena）、实体集合、组表与调度流资源
	for _, set := range w.components {
		set.Free()
	}
	w.components = nil
	if w.entities != nil {
		w.entities.Free()
	}
	if w.archetypes != nil {
		w.archetypes.clear()
	}
	w.destroyMu.Lock()
	w.pendingDestroy = nil
	w.destroyMu.Unlock()
	if w.systems != nil {
		w.systems.stop()
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
	if w.status != WorldStatusInitialized {
		panic("system register only in world init")
	}
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
	if w.status != WorldStatusInitialized {
		panic("system register only in world init")
	}
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
