package ecs

import "reflect"

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

type ShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]] struct{ getterBase }

func NewShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]](sys ISystem) *ShapeGetter[T, TP] {
	getter := &ShapeGetter[T, TP]{getterBase{sys: sys}}
	typ := reflect.TypeOf(getter)
	getter.init(typ, getter)
	var temp T
	var req []IRequirement
	sysReq := sys.Requirements()
	for _, t := range TP(&temp).eleTypes() {
		if r, ok := sysReq[t]; ok {
			req = append(req, r)
		}
	}
	getter.req = req
	return getter
}

func (s *ShapeGetter[T, TP]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[ShapeGetter[T, TP]]()
	}
	return s.typ
}

func (s *ShapeGetter[T, TP]) Iter() IShapeIterator[T, TP] {
	s.executeNum++
	var min ICollection
	for _, r := range s.req {
		c := s.sys.World().getComponents(r.Type())
		if min == nil || min.Len() > c.Len() {
			min = c
		}
	}
	return NewShapeIterator[T, TP](min, s.req)
}
