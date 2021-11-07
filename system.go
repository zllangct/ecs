package ecs

import (
	"container/list"
	"reflect"
	"sync"
	"unsafe"
)

type SystemLifeCircleType int

type ISystem interface {
	Type() reflect.Type
	Order() Order
	World() IWorld
	Requirements() map[reflect.Type]struct{}
	Emit(event string, args ...interface{})

	IsRequire(component IComponent) bool

	isRequire(componentType reflect.Type) bool
	setOrder(order Order)
	setRequirements(rqs ...IComponent)
	setRequirementsByType(rqs ...reflect.Type)
	checkComponent(entity *EntityInfo, com IComponent) IComponent
	baseInit(world *ecsWorld, ins ISystem)
	eventDispatch()
}

type ISystemTemplate interface {
	d074634084a1556083fcd17c0254b557()
}

type System[T any] struct {
	lock         sync.Mutex
	requirements map[reflect.Type]struct{}
	events       map[string]SysEventHandler
	eventQueue   *list.List
	order        Order
	world        *ecsWorld
	realType     reflect.Type
	isInited     bool
}

type ITest interface {
	Test()
}

type TestA[T any] struct {
	base testA[T]
}

type testA[T any] struct {

}

func (t *testA[T]) Test() {}

func (s System[T]) d074634084a1556083fcd17c0254b557() {}

func (s *System[T]) Ins() (sys ISystem) {
	(*iface)(unsafe.Pointer(&sys)).data = unsafe.Pointer(s)
	return
}

func (s *System[T]) RawIns() *T {
	return (*T)(unsafe.Pointer(s))
}

func (s *System[T]) EventRegister(event string, fn SysEventHandler) {
	s.events[event] = fn
}

func (s *System[T]) Emit(event string, args ...interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.eventQueue.PushBack(SystemCustomEvent{
		Event: event,
		Args:  args,
	})
}

func (s *System[T]) eventDispatch() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i := s.eventQueue.Front(); i != nil; i = i.Next() {
		e := i.Value.(SystemCustomEvent)
		if fn, ok := s.events[e.Event]; ok {
			err := TryAndReport(func() {
				fn(e.Args...)
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

func (s *System[T]) SetRequirements(rqs ...IComponent) {
	s.setRequirements(rqs...)
}

func (s *System[T]) setRequirements(rqs ...IComponent) {
	if s.isInited {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[value.Type()] = struct{}{}
	}
}

func (s *System[T]) setRequirementsByType(rqs ...reflect.Type) {
	if s.isInited {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[value] = struct{}{}
	}
}

func (s *System[T]) Requirements() map[reflect.Type]struct{} {
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
	s.requirements = map[reflect.Type]struct{}{}
	s.eventQueue = list.New()

	if ins.Order() == OrderInvalid {
		s.setOrder(OrderDefault)
	}
	s.world = world

	if i, ok := ins.(IEventInit); ok {
		i.Init()
	}

	s.isInited = true
}

func (s *System[T]) Type() reflect.Type {
	if s.realType == nil {
		s.realType = TypeOf[T]()
	}
	return s.realType
}

func (s *System[T]) setOrder(order Order) {
	if s.isInited {
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

func (s *System[T]) GetInterested(typ reflect.Type) interface{} {
	if _, ok := s.requirements[typ]; !ok {
		return nil
	}

	return s.World().getComponents(typ)
}

func (s *System[T]) GetInterestedNew() map[reflect.Type][]OperateInfo {
	ls := map[reflect.Type][]OperateInfo{}
	for typ, _ := range s.Requirements() {
		if n := s.World().getNewComponents(typ); n != nil {
			ls[typ] = n
		}
	}
	return ls
}

func (s *System[T]) CheckComponent(info *EntityInfo, com IComponent) IComponent {
	return s.checkComponent(info, com)
}

func (s *System[T]) checkComponent(entity *EntityInfo, com IComponent) IComponent {
	isRequire := s.IsRequire(com)
	if !isRequire {
		return nil
	}

	return entity.getComponent(com)
}

func (s *System[T]) GetEntityInfo(entity Entity) *EntityInfo {
	return s.world.getEntityInfo(entity)
}