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
	Type() reflect.Type
	Order() Order
	Requirements() map[reflect.Type]struct{}
	Call(label int) interface{}
}

type ISystemBaseInit interface{
	BaseInit(world *World)
}

type System[T any] struct {
	sync.Mutex
	requirements map[reflect.Type]struct{}
	order Order
	world *World
	realType reflect.Type
	isPreFilter  bool
}

func (s *System[T]) Ins() *T {
	return (*T)(unsafe.Pointer(s))
}

func (s *System[T]) Call(label int) interface{} {
	return nil
}

func (s *System[T]) SetRequirements(rqs ...IComponent) {
	if s.requirements == nil {
		s.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		s.requirements[reflect.TypeOf(value)] = struct{}{}
	}
}

func (s *System[T]) Requirements() map[reflect.Type]struct{} {
	return s.requirements
}

func (s *System[T]) BaseInit(world *World) {
	s.requirements = map[reflect.Type]struct{}{}
	s.SetOrder(ORDER_DEFAULT)
	s.world = world
}

func (s *System[T]) Type() reflect.Type {
	s.Lock()
	defer s.Unlock()

	if s.realType == nil {
		s.realType = reflect.TypeOf(*new(T))
	}
	return s.realType
}

func (s *System[T]) SetOrder(order Order) {
	s.Lock()
	defer s.Unlock()

	s.order = order
}

func (s *System[T]) Order() Order {
	s.Lock()
	defer s.Unlock()

	return s.order
}

func (s *System[T]) IsConcerned(com IComponent) bool {
	cType := com.Type()
	if _, concerned := s.requirements[cType]; concerned {
		for r := range s.requirements {
			if r != cType {
				if !com.Owner().Has(r) {
					concerned = false
					break
				}
			}
		}
		return concerned
	}
	return false
}

func (s *System[T]) GetWorld() *World {
	return s.world
}

