package ecs

import (
	"reflect"
	"unsafe"
)

type IShapeGetter interface {
	base() *getterBase
	getType() reflect.Type
}

type getterBase struct {
	sys        ISystem
	executeNum int64
	typ        reflect.Type
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
	readOnly   []bool
}

type ShapeGetter[T any] struct {
	getterBase
	initializer  SystemInitializer
	mainKeyIndex int
	subTypes     []uint16
	subOffset    []uintptr
	containers   []IComponentSet
	readOnly     []bool
	cur          *T
	valid        bool
}

func NewShapeGetter[T any](initializer SystemInitializer) *ShapeGetter[T] {
	if initializer.isValid() {
		panic("out of initialization stage")
	}
	sys := initializer.getSystem()
	getter := &ShapeGetter[T]{
		getterBase:  getterBase{sys: sys},
		initializer: initializer,
	}

	typ := reflect.TypeOf(getter)
	getter.init(typ, getter)

	sysReq := sys.GetRequirements()
	if sysReq == nil {
		return nil
	}

	getter.cur = new(T)
	typIns := reflect.TypeOf(*getter.cur)
	for i := 0; i < typIns.NumField(); i++ {
		field := typIns.Field(i)
		if !field.Type.Implements(reflect.TypeOf((*IComponent)(nil)).Elem()) || !sys.isRequire(field.Type.Elem()) {
			continue
		}
		if r, ok := sysReq[field.Type.Elem()]; ok {
			if r.getPermission() == ComponentReadOnly {
				getter.readOnly = append(getter.readOnly, true)
			} else {
				getter.readOnly = append(getter.readOnly, false)
			}
		}
		meta := sys.World().getComponentMetaInfoByType(field.Type.Elem())
		getter.subTypes = append(getter.subTypes, meta.it)
		getter.subOffset = append(getter.subOffset, field.Offset)
	}

	getter.containers = make([]IComponentSet, len(getter.subTypes))

	if len(getter.subTypes) == 0 {
		return nil
	}

	getter.valid = true

	return getter
}

func (s *ShapeGetter[T]) IsValid() bool {
	return s.valid
}

func (s *ShapeGetter[T]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[ShapeGetter[T]]()
	}
	return s.typ
}

func (s *ShapeGetter[T]) Get() IShapeIterator[T] {
	s.executeNum++

	if !s.valid {
		return EmptyShapeIter[T]()
	}

	var mainComponent ICollection
	var mainKeyIndex int
	for i := 0; i < len(s.subTypes); i++ {
		c := s.sys.World().getComponentSetByIntType(s.subTypes[i])
		if c == nil || c.Len() == 0 {
			return EmptyShapeIter[T]()
		}
		if mainComponent == nil || mainComponent.Len() > c.Len() {
			mainComponent = c
			mainKeyIndex = i
		}
		s.containers[i] = c
	}

	if s.mainKeyIndex == 0 {
		mainKeyIndex = s.mainKeyIndex
		mainComponent = s.containers[mainKeyIndex]
	}

	return NewShapeIterator[T](
		ShapeIndices{
			subTypes:   s.subTypes,
			subOffset:  s.subOffset,
			containers: s.containers,
			readOnly:   s.readOnly,
		},
		mainKeyIndex)
}

func (s *ShapeGetter[T]) GetSpecific(entity Entity) (*T, bool) {
	if !s.valid {
		return s.cur, false
	}
	for i := 0; i < len(s.subTypes); i++ {
		subPointer := s.containers[i].getPointerByEntity(entity)
		if subPointer == nil {
			return s.cur, false
		}
		if s.readOnly[i] {
			*(**byte)(unsafe.Pointer((uintptr)(unsafe.Pointer(s.cur)) + s.subOffset[i])) = &(*(*byte)(subPointer))
		} else {
			*(**byte)(unsafe.Pointer((uintptr)(unsafe.Pointer(s.cur)) + s.subOffset[i])) = (*byte)(subPointer)
		}
	}
	return s.cur, true
}

func (s *ShapeGetter[T]) SetGuide(component IComponent) *ShapeGetter[T] {
	meta := s.initializer.getSystem().World().getComponentMetaInfoByType(component.Type())
	for i, r := range s.subTypes {
		if r == meta.it {
			s.mainKeyIndex = i
			return s
		}
	}
	return s
}
