package ecs

import (
	"reflect"
	"unsafe"
)

type IShapeGetter interface {
	optimize(collections map[reflect.Type]interface{})
	base() *getterBase
	getType() reflect.Type
}

type IShape interface {
	shapeSize() int
	setID(ids []int64)
	getElements() []unsafe.Pointer
}

type ShapePointer[T ShapeObject] interface {
	IShape
	*T
}

type ShapeObject interface {
	elementSize() []uintptr
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

type Shape2[T1, T2 ComponentObject] struct {
	C1 *T1
	C2 *T2
}

func (s Shape2[T1, T2]) elementSize() []uintptr {
	return []uintptr{TypeOf[T1]().Size(), TypeOf[T2]().Size()}
}

type ShapeGetter2[T1, T2 ComponentObject] struct{ getterBase }

func NewShapeGetter2[T1, T2 ComponentObject](sys ISystem) *ShapeGetter2[T1, T2] {
	getter := &ShapeGetter2[T1, T2]{getterBase{sys: sys}}
	opt := sys.getOptimizer()
	if _, ok := opt.shapeUsage[reflect.TypeOf(getter)]; !ok {
		opt.shapeUsage[reflect.TypeOf(getter)] = getter
	}
	return getter
}

func (s *ShapeGetter2[T1, T2]) Get() Shape2[T1, T2] {
	// TODO 需要检查是否是系统感兴趣的组件
	// 记录该类型的使用次数
	s.executeNum++

	return Shape2[T1, T2]{}
}

func (s *Shape2[T1, T2]) optimize(collections map[reflect.Type]interface{}) {
	// 执行内存布局优化
}
