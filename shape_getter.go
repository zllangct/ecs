package ecs

import (
	"errors"
	"reflect"
	"sync"
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

var shapeCaches = sync.Map{}

type ShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]] struct{ getterBase }

func NewShapeGetter[T ShapeObject, TP ShapeObjectPointer[T]](sys ISystem) (*ShapeGetter[T, TP], error) {
	getter := &ShapeGetter[T, TP]{getterBase{sys: sys}}
	typ := reflect.TypeOf(getter)
	getter.init(typ, getter)
	var temp T
	var req []IRequirement
	sysReq := sys.Requirements()
	for _, t := range TP(&temp).eleTypes() {
		if r, ok := sysReq[t]; ok {
			req = append(req, r)
		} else {
			return nil, errors.New("component not interested")
		}
	}
	getter.req = req
	return getter, nil
}

func (s *ShapeGetter[T, TP]) getType() reflect.Type {
	if s.typ == nil {
		s.typ = TypeOf[ShapeGetter[T, TP]]()
	}
	return s.typ
}

func (s *ShapeGetter[T, TP]) Get() IShapeIterator[T, TP] {
	s.executeNum++
	var min ICollection
	for _, r := range s.req {
		c := s.sys.World().getComponents(r.Type())
		if c == nil || c.Len() == 0 {
			return EmptyShapeIter[T, TP]()
		}
		if min == nil || min.Len() > c.Len() {
			min = c
		}
	}

	if min.Len() == 0 {
		return EmptyShapeIter[T, TP]()
	}

	var cache []T
	obj, ok := shapeCaches.Load(TypeOf[T]())
	if ok {
		cache = obj.([]T)
	} else {
		cache, ok = s.cache(min)
		shapeCaches.LoadOrStore(TypeOf[T](), cache)
	}

	return NewShapeIterator[T, TP](cache)
}

func (s *ShapeGetter[T, TP]) cache(guide ICollection) ([]T, bool) {
	var cache []T = make([]T, guide.Len())
	var err error = nil
	index := 0
	guide.Range(func(ele any) bool {
		c, ok := ele.(IComponent)
		if !ok {
			err = errors.New("element not component")
			return false
		}
		ok = TP(&cache[index]).parse(c.Owner(), s.req)
		if ok {
			index++
		}
		return true
	})
	if err != nil {
		return nil, false
	}
	return cache[:index], true
}
