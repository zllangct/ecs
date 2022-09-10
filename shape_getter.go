package ecs

import (
	"reflect"
	"unsafe"
)

type IShape interface {
	base() *shapeBase
	getType() reflect.Type
}

type shapeBase struct {
	sys        ISystem
	executeNum int64
	typ        reflect.Type
}

func (s *shapeBase) base() *shapeBase {
	return s
}

func (s *shapeBase) init(typ reflect.Type, getter IShape) {
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

type Shape[T any] struct {
	shapeBase
	initializer  SystemInitConstraint
	mainKeyIndex int
	subTypes     []uint16
	subOffset    []uintptr
	containers   []IComponentSet
	readOnly     []bool
	cur          *T
	valid        bool
}

func NewShape[T any](initializer SystemInitConstraint) *Shape[T] {
	if initializer.isValid() {
		panic("out of initialization stage")
	}
	sys := initializer.getSystem()
	getter := &Shape[T]{
		shapeBase:   shapeBase{sys: sys},
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

func (s *Shape[T]) IsValid() bool {
	return s.valid
}

func (s *Shape[T]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[Shape[T]]()
	}
	return s.typ
}

func (s *Shape[T]) Get() IShapeIterator[T] {
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

func (s *Shape[T]) GetSpecific(entity Entity) (*T, bool) {
	if !s.valid {
		return s.cur, false
	}
	for i := 0; i < len(s.subTypes); i++ {
		subPointer := s.containers[i].getPointerByEntity(entity)
		if subPointer == nil {
			return s.cur, false
		}
		if s.readOnly[i] {
			*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.subOffset[i])) = &(*(*byte)(subPointer))
		} else {
			*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.subOffset[i])) = (*byte)(subPointer)
		}
	}
	return s.cur, true
}

func (s *Shape[T]) SetGuide(component IComponent) *Shape[T] {
	meta := s.initializer.getSystem().World().getComponentMetaInfoByType(component.Type())
	for i, r := range s.subTypes {
		if r == meta.it {
			s.mainKeyIndex = i
			return s
		}
	}
	return s
}
