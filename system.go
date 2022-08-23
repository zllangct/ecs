package ecs

import (
	"container/list"
	"errors"
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

const (
	SystemCustomEventInvalid CustomEventName = ""
	SystemCustomEventPause                   = "__internal__Pause"
	SystemCustomEventResume                  = "__internal__Resume"
	SystemCustomEventStop                    = "__internal__Stop"
)

type ISystem interface {
	Type() reflect.Type
	Order() Order
	World() IWorld
	Requirements() map[reflect.Type]IRequirement
	Emit(event CustomEventName, args ...interface{})
	IsRequire(component IComponent) (IRequirement, bool)
	ID() int64
	Pause()
	Resume()
	Stop()

	getPointer() unsafe.Pointer
	isRequire(componentType reflect.Type) (IRequirement, bool)
	setOrder(order Order)
	setRequirements(rqs ...IRequirement)
	getState() SystemState
	setState(state SystemState)
	setSecurity(isSafe bool)
	isThreadSafe() bool
	setExecuting(isExecuting bool)
	isExecuting() bool
	baseInit(world *ecsWorld, ins ISystem)
	eventDispatch()
	getOptimizer() *OptimizerReporter
	getGetterCache() map[reflect.Type]interface{}
}

type SystemObject interface {
	systemIdentification()
}

type SystemPointer[T SystemObject] interface {
	ISystem
	*T
}

type System[T SystemObject] struct {
	lock              sync.Mutex
	requirements      map[reflect.Type]IRequirement
	events            map[CustomEventName]CustomEventHandler
	getterCache       map[reflect.Type]interface{}
	eventQueue        *list.List
	order             Order
	optimizerReporter *OptimizerReporter
	world             *ecsWorld
	utility           *Utility[T]
	realType          reflect.Type
	state             SystemState
	isSafe            bool
	executing         bool
	id                int64
}

func (s System[T]) systemIdentification() {}

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

func (s *System[T]) EventRegister(event CustomEventName, fn CustomEventHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.events[event] = fn
}

func (s *System[T]) Emit(event CustomEventName, args ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.eventQueue.PushBack(CustomEvent{
		Event: event,
		Args:  args,
	})
}

func (s *System[T]) eventDispatch() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i := s.eventQueue.Front(); i != nil; i = i.Next() {
		e := i.Value.(CustomEvent)
		if fn, ok := s.events[e.Event]; ok {
			err := TryAndReport(func() {
				fn(e.Args)
			})
			if err != nil {
				Log.Error(err)
			}
		} else {
			Log.Errorf("event not found: %s", e.Event)
		}
	}

	s.eventQueue.Init()
}

func (s *System[T]) SetRequirements(rqs ...IRequirement) {
	s.setRequirements(rqs...)
}

func (s *System[T]) isInitialized() bool {
	return s.state >= SystemStateInit
}

func (s *System[T]) setRequirements(rqs ...IRequirement) {
	if s.isInitialized() {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]IRequirement{}
	}
	var typ reflect.Type
	for _, value := range rqs {
		typ = value.Type()
		s.requirements[typ] = value
		ComponentMeta.GetComponentMetaInfo(typ)
	}
}

func (s *System[T]) setSecurity(isSafe bool) {
	s.isSafe = isSafe
}
func (s *System[T]) isThreadSafe() bool {
	return s.isSafe
}

func (s *System[T]) SetUtility(utility *Utility[T]) {
	s.utility = utility
}

func (s *System[T]) GetUtility() *Utility[T] {
	return s.utility
}

func (s *System[T]) Pause() {
	s.Emit(SystemCustomEventPause, nil)
}

func (s *System[T]) Resume() {
	s.Emit(SystemCustomEventResume, nil)
}

func (s *System[T]) Stop() {
	s.Emit(SystemCustomEventStop, nil)
}

func (s *System[T]) pause(e []interface{}) error {
	if s.getState() == SystemStateUpdate {
		s.setState(SystemStatePause)
	} else {
		return errors.New("system not running")
	}
	return nil
}

func (s *System[T]) resume(e []interface{}) error {
	if s.getState() == SystemStatePause {
		s.setState(SystemStateUpdate)
	} else {
		return errors.New("system not pausing")
	}
	return nil
}

func (s *System[T]) stop(e []interface{}) error {
	if s.getState() == SystemStatePause {
		s.setState(SystemStateUpdate)
	} else {
		return errors.New("system not pausing")
	}
	return nil
}

func (s *System[T]) getState() SystemState {
	return s.state
}

func (s *System[T]) setState(state SystemState) {
	s.state = state
}

func (s *System[T]) setExecuting(isExecuting bool) {
	s.executing = isExecuting
}

func (s *System[T]) isExecuting() bool {
	return s.executing
}

func (s *System[T]) Requirements() map[reflect.Type]IRequirement {
	return s.requirements
}

func (s *System[T]) IsRequire(com IComponent) (IRequirement, bool) {
	return s.isRequire(com.Type())
}

func (s *System[T]) isRequire(typ reflect.Type) (IRequirement, bool) {
	r, ok := s.requirements[typ]
	return r, ok
}

func (s *System[T]) baseInit(world *ecsWorld, ins ISystem) {
	s.requirements = map[reflect.Type]IRequirement{}
	s.events = make(map[CustomEventName]CustomEventHandler)
	s.eventQueue = list.New()
	s.getterCache = map[reflect.Type]interface{}{}

	if ins.Order() == OrderInvalid {
		s.setOrder(OrderDefault)
	}
	s.world = world

	s.EventRegister(SystemCustomEventPause, func(i []interface{}) {
		err := s.pause(i)
		if err != nil {
			Log.Error(err)
		}
	})
	s.EventRegister(SystemCustomEventResume, func(i []interface{}) {
		err := s.resume(i)
		if err != nil {
			Log.Error(err)
		}
	})
	s.EventRegister(SystemCustomEventStop, func(i []interface{}) {
		err := s.stop(i)
		if err != nil {
			Log.Error(err)
		}
	})

	if i, ok := ins.(InitReceiver); ok {
		err := TryAndReport(func() {
			i.Init()
		})
		if err != nil {
			Log.Error(err)
		}
	}

	s.state = SystemStateInit
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

func (s *System[T]) GetEntityInfo(entity Entity) EntityInfo {
	return s.world.GetEntityInfo(entity)
}

// get optimizer
func (s *System[T]) getOptimizer() *OptimizerReporter {
	if s.optimizerReporter == nil {
		s.optimizerReporter = &OptimizerReporter{}
		s.optimizerReporter.init()
	}
	return s.optimizerReporter
}

func (s *System[T]) getGetterCache() map[reflect.Type]interface{} {
	return s.getterCache
}
