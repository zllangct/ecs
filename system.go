package ecs

import (
	"reflect"
	"sync"
	"unsafe"
)

type SystemLifeCircleType int

const (
	SYSTEM_LIFE_CIRCLE_TYPE_NONE SystemLifeCircleType = iota
	SYSTEM_LIFE_CIRCLE_TYPE_Default
	SYSTEM_LIFE_CIRCLE_TYPE_ONCE
	SYSTEM_LIFE_CIRCLE_TYPE_REPEAT
)

type ISystem interface {
	//Init()
	Type() reflect.Type
	Order() Order
	World() *World
	Requirements() map[reflect.Type]struct{}
	Call(label int) interface{}

	baseInit(world *World, ins ISystem)
}

type System[T any] struct {
	lock sync.Mutex
	requirements map[reflect.Type]struct{}
	order Order
	world *World
	realType reflect.Type
	isInited  bool
}

func (s *System[T]) Ins() *T {
	return (*T)(unsafe.Pointer(s))
}

func (s *System[T]) Call(label int) interface{} {
	return nil
}

func (s *System[T]) SetRequirements(rqs ...IComponentTemplate) {
	//s.lock.Lock()
	//defer s.lock.Unlock()

	if s.isInited {
		return
	}
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[value.ComponentType()] = struct{}{}
	}
}

func (s *System[T]) Requirements() map[reflect.Type]struct{} {
	//s.lock.Lock()
	//defer s.lock.Unlock()

	return s.requirements
}

func (s *System[T]) baseInit(world *World, ins ISystem) {
	s.requirements = map[reflect.Type]struct{}{}
	s.SetOrder(ORDER_DEFAULT)
	s.world = world

	if i, ok := ins.(IEventInit); ok {
		i.Init()
	}

	s.isInited = true
}

func (s *System[T]) Type() reflect.Type {
	//s.lock.Lock()
	//defer s.lock.Unlock()

	if s.realType == nil {
		s.realType = reflect.TypeOf(*new(T))
	}
	return s.realType
}

func (s *System[T]) SetOrder(order Order) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.isInited {
		return
	}

	s.order = order
}

func (s *System[T]) Order() Order {
	//s.lock.Lock()
	//defer s.lock.Unlock()

	return s.order
}

func (s *System[T]) World() *World {
	return s.world
}

func (s *System[T]) GetInterested(typ reflect.Type) interface{}{
	if _, ok := s.requirements[typ]; !ok {
		return nil
	}

	return s.World().getComponents(typ)
}


func (s *System[T]) GetInterestedNew() map[reflect.Type][]ComponentOptResult {
	ls := map[reflect.Type][]ComponentOptResult{}
	for typ, _ := range s.Requirements() {
		if n :=s.World().getNewComponents(typ); n != nil {
			ls[typ] = n
		}
	}
	return ls
}

