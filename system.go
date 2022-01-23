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
	Requirements() map[reflect.Type]struct{}
	Emit(event CustomEventName, args ...interface{})
	IsRequire(component IComponent) bool
	ID() int64
	Pause()
	Resume()
	Stop()

	isRequire(componentType reflect.Type) bool
	setOrder(order Order)
	setRequirements(rqs ...IComponent)
	setRequirementsByType(rqs ...reflect.Type)
	getState() SystemState
	setState(state SystemState)
	checkComponent(entity *EntityInfo, com IComponent) IComponent
	baseInit(world *ecsWorld, ins ISystem)
	eventDispatch()
}

type SystemObject interface {
	systemIdentification()
}

type SystemPointer[T SystemObject] interface {
	ISystem
	*T
}

type System[T SystemObject, TP SystemPointer[T]] struct {
	lock         sync.Mutex
	requirements map[reflect.Type]struct{}
	events       map[CustomEventName]CustomEventHandler
	eventQueue   *list.List
	order        Order
	world        *ecsWorld
	realType     reflect.Type
	state        SystemState
	id           int64
}

func (s System[T, TP]) systemIdentification() {}

func (s *System[T, TP]) Ins() (sys ISystem) {
	(*iface)(unsafe.Pointer(&sys)).data = unsafe.Pointer(s)
	return
}

func (s *System[T, TP]) rawInstance() *T {
	return (*T)(unsafe.Pointer(s))
}

func (s *System[T, TP]) ID() int64 {
	if s.id == 0 {
		s.id = UniqueID()
	}
	return s.id
}

func (s *System[T, TP]) EventRegister(event CustomEventName, fn CustomEventHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.events[event] = fn
}

func (s *System[T, TP]) Emit(event CustomEventName, args ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.eventQueue.PushBack(CustomEvent{
		Event: event,
		Args:  args,
	})
}

func (s *System[T, TP]) eventDispatch() {
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

func (s *System[T, TP]) SetRequirements(rqs ...IComponent) {
	s.setRequirements(rqs...)
}

func (s *System[T, TP]) isInitialized() bool {
	return s.state >= SystemStateInit
}

func (s *System[T, TP]) setRequirements(rqs ...IComponent) {
	if s.isInitialized() {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[value.Type()] = struct{}{}
	}
}

func (s *System[T, TP]) Pause() {
	s.Emit(SystemCustomEventPause, nil)
}

func (s *System[T, TP]) Resume() {
	s.Emit(SystemCustomEventResume, nil)
}

func (s *System[T, TP]) Stop() {
	s.Emit(SystemCustomEventStop, nil)
}

func (s *System[T, TP]) pause(e []interface{}) error {
	if s.getState() == SystemStateUpdate {
		s.setState(SystemStatePause)
	} else {
		return errors.New("system not running")
	}
	return nil
}

func (s *System[T, TP]) resume(e []interface{}) error {
	if s.getState() == SystemStatePause {
		s.setState(SystemStateUpdate)
	} else {
		return errors.New("system not pausing")
	}
	return nil
}

func (s *System[T, TP]) stop(e []interface{}) error {
	if s.getState() == SystemStatePause {
		s.setState(SystemStateUpdate)
	} else {
		return errors.New("system not pausing")
	}
	return nil
}

func (s *System[T, TP]) getState() SystemState {
	return s.state
}

func (s *System[T, TP]) setState(state SystemState) {
	s.state = state
}

func (s *System[T, TP]) setRequirementsByType(rqs ...reflect.Type) {
	if s.isInitialized() {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[value] = struct{}{}
	}
}

func (s *System[T, TP]) Requirements() map[reflect.Type]struct{} {
	return s.requirements
}

func (s *System[T, TP]) IsRequire(com IComponent) bool {
	return s.isRequire(com.Type())
}

func (s *System[T, TP]) isRequire(typ reflect.Type) bool {
	_, ok := s.requirements[typ]
	return ok
}

func (s *System[T, TP]) baseInit(world *ecsWorld, ins ISystem) {
	s.requirements = map[reflect.Type]struct{}{}
	s.events = make(map[CustomEventName]CustomEventHandler)
	s.eventQueue = list.New()

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

func (s *System[T, TP]) Type() reflect.Type {
	if s.realType == nil {
		s.realType = TypeOf[T]()
	}
	return s.realType
}

func (s *System[T, TP]) setOrder(order Order) {
	if s.isInitialized() {
		return
	}

	s.order = order
}

func (s *System[T, TP]) Order() Order {
	return s.order
}

func (s *System[T, TP]) World() IWorld {
	return s.world
}

func (s *System[T, TP]) GetInterested(typ reflect.Type) interface{} {
	if _, ok := s.requirements[typ]; !ok {
		return nil
	}

	return s.World().getComponents(typ)
}

func (s *System[T, TP]) GetInterestedNew() map[reflect.Type][]OperateInfo {
	ls := map[reflect.Type][]OperateInfo{}
	for typ, _ := range s.Requirements() {
		if n := s.World().getNewComponents(typ); n != nil {
			ls[typ] = n
		}
	}
	return ls
}

func (s *System[T, TP]) CheckComponent(info *EntityInfo, com IComponent) IComponent {
	return s.checkComponent(info, com)
}

func (s *System[T, TP]) checkComponent(entity *EntityInfo, com IComponent) IComponent {
	isRequire := s.IsRequire(com)
	if !isRequire {
		return nil
	}

	return entity.getComponent(com)
}

func (s *System[T, TP]) GetEntityInfo(entity Entity) *EntityInfo {
	return s.world.getEntityInfo(entity)
}
