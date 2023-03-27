package ecs

import (
	"reflect"
	"sync"
	"unsafe"
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

type SystemInitConstraint struct {
	sys *ISystem
}

func (s *SystemInitConstraint) getSystem() ISystem {
	if s.sys == nil {
		panic("out of initialization stage")
	}
	return *s.sys
}

func (s *SystemInitConstraint) SetBroken(reason string) {
	(*s.sys).setBroken()
	panic(reason)
}

func (s *SystemInitConstraint) isValid() bool {
	return *s.sys == nil
}

type ISystem interface {
	Type() reflect.Type
	Order() Order
	World() IWorld
	GetRequirements() map[reflect.Type]IRequirement
	IsRequire(component IComponent) bool
	ID() int64
	GetUtility() IUtility

	pause()
	resume()
	stop()
	getPointer() unsafe.Pointer
	isRequire(componentType reflect.Type) bool
	setOrder(order Order)
	setRequirements(initializer SystemInitConstraint, rqs ...IRequirement)
	getState() SystemState
	setState(state SystemState)
	setSecurity(isSafe bool)
	isThreadSafe() bool
	setExecuting(isExecuting bool)
	isExecuting() bool
	baseInit(world *ecsWorld, ins ISystem)
	getOptimizer() *OptimizerReporter
	getGetterCache() *GetterCache
	setBroken()
	isValid() bool
	setUtility(u IUtility)
}

type SystemObject interface {
	__SystemIdentification()
}

type SystemPointer[T SystemObject] interface {
	ISystem
	*T
}

type systemIdentification struct{}

func (s systemIdentification) __SystemIdentification() {}

type System[T SystemObject] struct {
	systemIdentification
	lock              sync.Mutex
	requirements      map[reflect.Type]IRequirement
	getterCache       *GetterCache
	order             Order
	optimizerReporter *OptimizerReporter
	world             *ecsWorld
	utility           IUtility
	realType          reflect.Type
	state             SystemState
	valid             bool
	isSafe            bool
	executing         bool
	id                int64
}

func (s *System[T]) instance() (sys ISystem) {
	(*iface)(unsafe.Pointer(&sys)).data = unsafe.Pointer(s)
	return
}

func (s *System[T]) rawInstance() *T {
	return (*T)(unsafe.Pointer(s))
}

func (s *System[T]) ID() int64 {
	if s.id == 0 {
		s.id = LocalUniqueID()
	}
	return s.id
}

func (s *System[T]) SetRequirements(initializer SystemInitConstraint, rqs ...IRequirement) {
	if initializer.isValid() {
		panic("out of initialization stage")
	}
	s.setRequirements(initializer, rqs...)
}

func (s *System[T]) setRequirementsInternal(rqs ...IRequirement) {
	if s.requirements == nil {
		s.requirements = map[reflect.Type]IRequirement{}
	}
	var typ reflect.Type
	for _, value := range rqs {
		typ = value.Type()
		s.requirements[typ] = value
	}
}

func (s *System[T]) isInitialized() bool {
	return s.state >= SystemStateInit
}

func (s *System[T]) setRequirements(initializer SystemInitConstraint, rqs ...IRequirement) {
	if s.requirements == nil {
		s.requirements = map[reflect.Type]IRequirement{}
	}
	var typ reflect.Type
	for _, value := range rqs {
		typ = value.Type()
		value.check(initializer)
		s.requirements[typ] = value
		s.World().getComponentMetaInfoByType(typ)
	}
}

func (s *System[T]) setUtility(u IUtility) {
	s.utility = u
}

func (s *System[T]) setSecurity(isSafe bool) {
	s.isSafe = isSafe
}
func (s *System[T]) isThreadSafe() bool {
	return s.isSafe
}

func (s *System[T]) GetUtility() IUtility {
	return s.utility
}

func (s *System[T]) pause() {
	if s.getState() == SystemStateUpdate {
		s.setState(SystemStatePause)
	}
}

func (s *System[T]) resume() {
	if s.getState() == SystemStatePause {
		s.setState(SystemStateUpdate)
	}
}

func (s *System[T]) stop() {
	if s.getState() < SystemStateDestroy {
		s.setState(SystemStateDestroy)
	}
}

func (s *System[T]) getState() SystemState {
	return s.state
}

func (s *System[T]) setState(state SystemState) {
	s.state = state
}

func (s *System[T]) setBroken() {
	s.valid = false
}

func (s *System[T]) isValid() bool {
	return s.valid
}

func (s *System[T]) setExecuting(isExecuting bool) {
	s.executing = isExecuting
}

func (s *System[T]) isExecuting() bool {
	return s.executing
}

func (s *System[T]) GetRequirements() map[reflect.Type]IRequirement {
	return s.requirements
}

func (s *System[T]) IsRequire(com IComponent) bool {
	return s.isRequire(com.Type())
}

func (s *System[T]) isRequire(typ reflect.Type) bool {
	_, ok := s.requirements[typ]
	return ok
}

func (s *System[T]) baseInit(world *ecsWorld, ins ISystem) {
	s.requirements = map[reflect.Type]IRequirement{}
	s.getterCache = NewGetterCache(len(s.requirements))

	if ins.Order() == OrderInvalid {
		s.setOrder(OrderDefault)
	}
	s.world = world

	s.valid = true

	initializer := SystemInitConstraint{}
	is := ISystem(s)
	initializer.sys = &is
	if i, ok := ins.(InitReceiver); ok {
		err := TryAndReport(func() error {
			return i.Init(initializer)
		})
		if err != nil {
			Log.Error(err)
		}
	}
	*initializer.sys = nil
	initializer.sys = nil

	s.state = SystemStateStart
}

func (s *System[T]) getPointer() unsafe.Pointer {
	return unsafe.Pointer(s)
}

func (s *System[T]) Type() reflect.Type {
	if s.realType == nil {
		s.realType = TypeOf[T]()
	}
	return s.realType
}

func (s *System[T]) setOrder(order Order) {
	if s.isInitialized() {
		return
	}

	s.order = order
}

func (s *System[T]) Order() Order {
	return s.order
}

func (s *System[T]) World() IWorld {
	return s.world
}

func (s *System[T]) GetEntityInfo(entity Entity) (*EntityInfo, bool) {
	return s.world.getEntityInfo(entity)
}

// get optimizer
func (s *System[T]) getOptimizer() *OptimizerReporter {
	if s.optimizerReporter == nil {
		s.optimizerReporter = &OptimizerReporter{}
		s.optimizerReporter.init()
	}
	return s.optimizerReporter
}

func (s *System[T]) getGetterCache() *GetterCache {
	return s.getterCache
}
