package ecs

import (
	"errors"
	"reflect"
)

type IShapeGetter interface {
	base() *getterBase
	getType() reflect.Type
}

type getterBase struct {
	sys        ISystem
	executeNum int64
	typ        reflect.Type
	req        []IRequirement
}

func (s *getterBase) base() *getterBase {
	return s
}

func (s *getterBase) init(typ reflect.Type, getter IShapeGetter) {
	opt := s.sys.getOptimizer()
	if _, ok := opt.shapeUsage[typ]; !ok {
		opt.shapeUsage[typ] = getter
	}
}

type ShapeIndices struct {
	subTypes   []uint16
	subOffset  []uintptr
	containers []IComponentSet
}

type ShapeGetter[T any] struct {
	getterBase

	subTypes   []uint16
	subOffset  []uintptr
	containers []IComponentSet
}

func NewShapeGetter[T any](initializer *SystemInitializer) (*ShapeGetter[T], error) {
	if initializer.sys == nil {
		return nil, errors.New("invalid system initializer")
	}
	sys := initializer.sys
	getter := &ShapeGetter[T]{
		getterBase: getterBase{sys: sys},
	}

	typ := reflect.TypeOf(getter)
	getter.init(typ, getter)

	sysReq := sys.GetRequirements()
	if sysReq == nil {
		return nil, errors.New("system requirement should be set before use shape")
	}

	typIns := reflect.TypeOf(*new(T))
	for i := 0; i < typIns.NumField(); i++ {
		field := typIns.Field(i)
		if !field.Type.Implements(reflect.TypeOf((*IComponent)(nil)).Elem()) || !sys.isRequire(field.Type.Elem()) {
			continue
		}
		intType := GetIntType(initializer.sys.World(), field.Type.Elem())
		getter.subTypes = append(getter.subTypes, intType)
		getter.subOffset = append(getter.subOffset, field.Offset)
		getter.containers = append(getter.containers, sys.World().getComponentSetByIntType(intType))
	}

	if len(getter.subTypes) == 0 {
		return nil, errors.New("no valid component found in shape")
	}

	return getter, nil
}

func (s *ShapeGetter[T]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[ShapeGetter[T]]()
	}
	return s.typ
}

func (s *ShapeGetter[T]) Get() IShapeIterator[T] {
	s.executeNum++
	var mainComponent ICollection
	var mainKeyIndex int
	for i, r := range s.containers {
		if mainComponent == nil || mainComponent.Len() > r.Len() {
			mainComponent = r
			mainKeyIndex = i
		}
	}

	if mainComponent.Len() == 0 {
		return EmptyShapeIter[T]()
	}

	return NewShapeIterator[T](ShapeIndices{
		subTypes:   s.subTypes,
		subOffset:  s.subOffset,
		containers: s.containers}, mainKeyIndex)
}
