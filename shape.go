package ecs

import (
	"reflect"
)

type IShapeGetter interface {
	base() *getterBase
	getType() reflect.Type
}

type ShapeObject interface {
	shapeBase() ShapeBase
	parse(info *EntityInfo, typs []reflect.Type) bool
}

type ShapeObjectPointer[T ShapeObject] interface {
	*T
}

type getterBase struct {
	sys        ISystem
	executeNum int64
	typ        reflect.Type
	eleType    []reflect.Type
}

func (s *getterBase) base() *getterBase {
	return s
}

func (s *getterBase) getIter() []interface{} {
	var cs []interface{}
	for _, t := range s.eleType {
		c := s.sys.World().getComponents(t)
		cs = append(cs, c)
	}
	return cs
}

func (s *getterBase) init(typ reflect.Type, getter IShapeGetter) {
	opt := s.sys.getOptimizer()
	if _, ok := opt.shapeUsage[typ]; !ok {
		opt.shapeUsage[typ] = getter
	}
}

type ShapeBase struct {
	entity Entity
}

func (s ShapeBase) shapeBase() ShapeBase {
	return s
}

type Shape2[T1, T2 ComponentObject] struct {
	ShapeBase
	C1 *T1
	C2 *T2
}

func (s *Shape2[T1, T2]) baseIndex(sys ISystem, typs []reflect.Type) func() *EntityInfo {
	c := sys.World().getComponents(typs[0])
	_ = c
	return nil
}

func (s *Shape2[T1, T2]) parse(info *EntityInfo, typs []reflect.Type) bool {
	if len(typs) != 2 {
		return false
	}
	c1 := info.getComponentByTypeInSystem(typs[0])
	if c1 == nil {
		return false
	}
	c2 := info.getComponentByTypeInSystem(typs[1])
	if c2 == nil {
		return false
	}
	//s.C1 = c1.(*T1)
	//s.C2 = c2.(*T2)
	return true
}

type ShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]] struct{ getterBase }

func NewShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]](sys ISystem) *ShapeGetter[T, TP] {
	getter := &ShapeGetter[T, TP]{getterBase{sys: sys}}
	typ := reflect.TypeOf(getter)
	getter.init(typ, getter)
	return getter
}

func (s *ShapeGetter[T, TP]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[ShapeGetter[T, TP]]()
	}
	return s.typ
}

func (s *ShapeGetter[T, TP]) Get() *T {
	// TODO 需要检查是否是系统感兴趣的组件
	// 记录该类型的使用次数
	s.executeNum++

	return nil
}

func GetRelated[T ShapeObject]() *T {

	return nil
}

func (s *ShapeGetter[T, TP]) getRelated(sys ISystem, typ []reflect.Type) IShapeIterator {
	var min ICollection
	for _, t := range typ {
		c := sys.World().getComponents(t)
		if min.Len() > c.Len() {
			min = c
		}
	}
	return NewShapeIterator(min, typ, false)
}
