package ecs

import (
	"fmt"
)

type SystemType uint8

const (
	SystemTypeInvalid SystemType = iota
	SystemTypeLight
	SystemTypeStandard
)

type SystemState uint8

const (
	SystemStateInvalid SystemState = iota
	SystemStateInit
	SystemStateStart
	SystemStatePause
	SystemStateUpdate
	SystemStateDestroy
	SystemStateDestroyed
)

type System interface{}

type SystemPointer[T any] interface {
	UpdateReceiver
	*T
}

type SystemInfo interface {
	isValid() bool
	getState() SystemState
	setState(state SystemState)
	getType() SystemType
	getDeps() []ComponentDependency
	id() uint64
	name() string
	getOrder() Order
	getRaw() any
	impl(stage Stage) bool
	getContext() *SystemContext
	getDep(intType ComponentIntType) (ComponentDependency, bool)
	getOptReporter() *optReporter
}

type SystemConfig struct {
	name         string
	dependencies []ComponentDependency
	stage        Stage
	order        Order
	// groups 本 system 声明的固定访问组（布局优化用，World.Register 时收集）
	groups [][]ComponentIntType
}

func (s *SystemConfig) initDefault() {
	s.name = "unnamed"
	s.stage = StageUpdate
	s.order = OrderAppend
}

type SystemOption func(config *SystemConfig)

func WithDeps(comp ...ComponentDependency) SystemOption {
	return func(c *SystemConfig) {
		c.dependencies = append(c.dependencies, comp...)
	}
}

func WithDep[T any, TP ComponentPointer[T]](writable ...Writable) SystemOption {
	return func(c *SystemConfig) {
		c.dependencies = append(c.dependencies, Dep[T, TP](writable...))
	}
}

func WithDepReadOnly[T any, TP ComponentPointer[T]]() SystemOption {
	return func(c *SystemConfig) {
		c.dependencies = append(c.dependencies, Dep[T, TP](ReadOnly))
	}
}

// WithGroup 声明本 system 固定联合访问的组件组（布局优化，不取代 WithDeps 权限声明）。
// 组内组件必须已 RegisterComponent 且非 disposable/nomadic（World.Register 期校验）；
// 至少 2 个组件类型；多 system 声明的重叠组会合并为并集组。
func WithGroup(deps ...ComponentDependency) SystemOption {
	return func(c *SystemConfig) {
		g := make([]ComponentIntType, 0, len(deps))
		for _, d := range deps {
			g = append(g, d.intType())
		}
		c.groups = append(c.groups, g)
	}
}

func WithStage(stage Stage) SystemOption {
	return func(c *SystemConfig) {
		c.stage = stage
	}
}

func WithOrder(order Order) SystemOption {
	return func(c *SystemConfig) {
		c.order = order
	}
}

func WithName(name string) SystemOption {
	return func(c *SystemConfig) {
		c.name = name
	}
}

// SystemConstraint 约束 SystemContext 只在 system 执行期间（含 Init）有效
type SystemConstraint struct {
	active bool
}

func (s *SystemConstraint) isValid() bool {
	return s.active
}

// activate 在 system 回调执行前激活约束
func (s *SystemConstraint) activate() {
	s.active = true
}

// deactivate 在 system 回调返回后关闭约束
func (s *SystemConstraint) deactivate() {
	s.active = false
}

type SystemContext struct {
	constraint SystemConstraint
	world      *World
	info       SystemInfo
	// queryCache 查询预解析缓存（shapeKey → queryCached），跨帧复用，
	// world.compVersion 变化或组表未定型时失效重建
	queryCache map[FixedCompound]*queryCached
}

// AddComponents 在 system 执行期间提交组件添加操作，写入本 system 的独立队列
// （单写无锁），帧同步点统一生效。语义与 EntityInfo.Add 一致：已存在/游牧组件跳过。
func (ctx *SystemContext) AddComponents(entity Entity, comps ...Component) {
	if !ctx.constraint.isValid() {
		return
	}
	info, ok := ctx.world.getEntityInfo(entity)
	if !ok {
		return
	}
	for _, comp := range comps {
		if info.compound.Exist(GetIntTypeByComp(comp)) {
			continue
		}
		if comp.IsNomadic() {
			continue
		}
		ctx.world.opLog.operateForSystem(ctx.info.id(), Operate{
			Entity: entity,
			Op:     ComponentOperateAdd,
			Comp:   comp,
		})
	}
}

type SystemInfoInstance struct {
	systemId   uint64
	systemName string
	impls      uint16
	config     *SystemConfig
	state      SystemState
	ctx        SystemContext
	typ        SystemType
	world      *World
	reporter   *optReporter
	raw        any
}

func newSystem(world *World, system any, typ SystemType) *SystemInfoInstance {
	c := &SystemConfig{}
	c.initDefault()
	impls := implsCheck(system)
	info := &SystemInfoInstance{
		config:   c,
		systemId: LocalUniqueID(),
		typ:      typ,
		raw:      system,
		impls:    impls,
		world:    world,
		reporter: newOptReporter(),
	}
	return info
}

func (s *SystemInfoInstance) init(opts ...SystemOption) {
	for _, opt := range opts {
		opt(s.config)
	}
}

func (s *SystemInfoInstance) getDep(it ComponentIntType) (ComponentDependency, bool) {
	for _, d := range s.config.dependencies {
		if d.intType() == it {
			return d, true
		}
	}
	return nil, false
}

func (s *SystemInfoInstance) getContext() *SystemContext {
	if s.ctx.world == nil {
		s.ctx.world = s.world
	}
	if s.ctx.info == nil {
		s.ctx.info = s
	}
	return &s.ctx
}

func (s *SystemInfoInstance) isValid() bool {
	return true
}

func (s *SystemInfoInstance) getState() SystemState {
	return s.state
}

func (s *SystemInfoInstance) setState(state SystemState) {
	s.state = state
}

func (s *SystemInfoInstance) getType() SystemType {
	return s.typ
}

func (s *SystemInfoInstance) getRaw() any {
	return s.raw
}

func (s *SystemInfoInstance) getDeps() []ComponentDependency {
	return s.config.dependencies
}

func (s *SystemInfoInstance) id() uint64 {
	return s.systemId
}

func (s *SystemInfoInstance) impl(stage Stage) bool {
	return s.impls>>stage&1 == 1
}

func (s *SystemInfoInstance) name() string {
	if len(s.systemName) == 0 {
		if s.config.name == "unnamed" {
			s.systemName = fmt.Sprintf("[%d]unnamed", s.systemId)
		} else {
			s.systemName = s.config.name
		}
	}
	return s.systemName
}

func (s *SystemInfoInstance) getOrder() Order {
	return s.config.order
}

func (s *SystemInfoInstance) getOptReporter() *optReporter {
	return s.reporter
}
