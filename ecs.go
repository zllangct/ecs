package ecs

import (
	"iter"
)

var ECS string = "ecs"

func Query(ctx *SystemContext, opt ...QueryOption) QueryIterator {
	qi := QueryIterator{
		ctx: ctx,
	}

	if !ctx.constraint.isValid() || len(opt) == 0 {
		return qi
	}

	c := QueryConfig{}
	c.initDefault()
	for _, option := range opt {
		option(&c)
	}
	qi.config = c

	var minSet ComponentSet
	for _, buddy := range c.queryBuddies {
		s, ok := ctx.world.getComponentSet(buddy)
		if ok {
			if minSet == nil {
				minSet = s
				continue
			}
			if s.Len() < minSet.Len() {
				minSet = s
			}
		}
	}

	if minSet == nil {
		return qi
	}
	qi.minSet = minSet

	return qi
}

func GetComponents[T ComponentObject, TP ComponentPointer[T]](ctx *SystemContext) iter.Seq2[EntityIndex, *T] {
	empty := func(yield func(EntityIndex, *T) bool) {
		return
	}
	if !ctx.constraint.isValid() {
		return empty
	}
	it := GetIntType[T, TP]()

	dep, ok := ctx.info.getDep(it)
	if !ok {
		return empty
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return empty
	}
	set, ok := s.(*CSet[T])
	if !ok {
		return empty
	}

	if dep.readonly() {
		return set.IterReadOnly()
	} else {
		return set.Iter()
	}
}

func GetBuddy[T ComponentObject, TP ComponentPointer[T]](ctx *SystemContext, index EntityIndex) (*T, bool) {
	if !ctx.constraint.isValid() {
		return nil, false
	}
	it := GetIntType[T, TP]()
	dep, ok := ctx.info.getDep(it)
	if !ok {
		return nil, false
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return nil, false
	}
	b := s.get(index)
	if b == nil {
		return nil, false
	}
	if dep.readonly() {
		return &*(*T)(b), true
	} else {
		return (*T)(b), true
	}
}
