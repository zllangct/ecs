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

type SystemInfo interface {
	isValid() bool
	getState() SystemState
	setState(state SystemState)
	getType() SystemType
	getDeps() []Dependency
	id() uint64
	name() string
	getOrder() Order
	getRaw() any
	impl(stage Stage) bool
	getContext() *SystemContext
	getDep(intType ComponentIntType) (Dependency, bool)
	getOptReporter() *optReporter
}

type SystemConfig struct {
	name         string
	dependencies []Dependency
	stage        Stage
	order        Order
}

func (s *SystemConfig) initDefault() {
	s.name = "unnamed"
	s.stage = StageUpdate
	s.order = OrderAppend
}

type SystemOption func(config *SystemConfig)

func WithDeps(comp ...Dependency) SystemOption {
	return func(c *SystemConfig) {
		c.dependencies = append(c.dependencies, comp...)
	}
}

func WithDep[T ComponentObject, TP ComponentPointer[T]](writable ...Writable) SystemOption {
	return func(c *SystemConfig) {
		c.dependencies = append(c.dependencies, Dep[T, TP](writable...))
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

type SystemConstraint struct {
	outdated bool
}

func (s *SystemConstraint) isValid() bool {
	return s.outdated
}

func (s *SystemConstraint) reset() {
	s.outdated = true
}

func (s *SystemConstraint) setOutdated() {
	s.outdated = false
}

type SystemContext struct {
	constraint SystemConstraint
	world      *world
	info       SystemInfo
}

type SystemInfoInstance struct {
	systemId   uint64
	systemName string
	impls      uint16
	config     *SystemConfig
	state      SystemState
	ctx        SystemContext
	typ        SystemType
	world      *world
	reporter   *optReporter
	raw        any
}

func newSystem(world *world, system any, typ SystemType) *SystemInfoInstance {
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

func (s *SystemInfoInstance) getDep(it ComponentIntType) (Dependency, bool) {
	for _, d := range s.config.dependencies {
		if d.intType() == it {
			return d, true
		}
	}
	return 0, false
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

func (s *SystemInfoInstance) getDeps() []Dependency {
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
