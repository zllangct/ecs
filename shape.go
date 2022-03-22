package ecs

import (
	"reflect"
	"unsafe"
)

type IShapeGetter interface {
	base() *getterBase
	getType() reflect.Type
}

type IShape interface {
	GetEntity() Entity
	setEntity(e Entity)
	parse(head unsafe.Pointer, eleSize []uintptr)
	getElements() []unsafe.Pointer
}

type ShapePointer[T ShapeObject] interface {
	IShape
	*T
}

type ShapeObject interface {
	//elementSize() []uintptr
}

type getterBase struct {
	sys        ISystem
	executeNum int64
	typ        reflect.Type
	eleSize    []int
}

func (s *getterBase) base() *getterBase {
	return s
}

func (s *getterBase) getType() reflect.Type {
	return s.typ
}

type ShapeBase struct {
	head   unsafe.Pointer
	entity Entity
}

func (s *ShapeBase) GetEntity() Entity {
	return s.entity
}

// set entity
func (s *ShapeBase) setEntity(e Entity) {
	s.entity = e
}

type Shape2[T1, T2 ComponentObject] struct {
	ShapeBase
	C1 *T1
	C2 *T2
}

func (s Shape2[T1, T2]) getElements() []unsafe.Pointer {
	return []unsafe.Pointer{unsafe.Pointer(s.C1), unsafe.Pointer(s.C2)}
}

func (s *Shape2[T1, T2]) parse(head unsafe.Pointer, eleSize []uintptr) {
	s.head = head
	offset := uintptr(0)
	s.setEntity(getComponentOwnerEntity(unsafe.Pointer(uintptr(head))))
	s.C1 = (*T1)(unsafe.Pointer(uintptr(head) + offset))
	offset += eleSize[0]
	s.C2 = (*T2)(unsafe.Pointer(uintptr(head) + offset))
}

type ShapeGetter2[T1, T2 ComponentObject] struct{ getterBase }

func NewShapeGetter2[T1, T2 ComponentObject](sys ISystem) *ShapeGetter2[T1, T2] {
	getter := &ShapeGetter2[T1, T2]{getterBase{sys: sys}}
	opt := sys.getOptimizer()
	if _, ok := opt.shapeUsage[getter.getType()]; !ok {
		opt.shapeUsage[getter.getType()] = getter
	}
	return getter
}

func (s *ShapeGetter2[T1, T2]) Get() Shape2[T1, T2] {
	// TODO 需要检查是否是系统感兴趣的组件
	// 记录该类型的使用次数
	s.executeNum++

	return Shape2[T1, T2]{}
}
