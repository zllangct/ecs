package ecs

import "iter"

type QueryBuddies = Compound

func NewQueryBuddies(intType ...ComponentIntType) *QueryBuddies {
	b := &QueryBuddies{}
	for _, componentIntType := range intType {
		b.Add(componentIntType)
	}
	return b
}

type QueryConfig struct {
	queryBuddies QueryBuddies
}

func (q *QueryConfig) initDefault() {

}

type QueryOption func(q *QueryConfig)

func WithBuddies(buddies QueryBuddies) QueryOption {
	return func(q *QueryConfig) {
		q.queryBuddies.Merge(buddies)
	}
}

func WithComp[T ComponentObject, TP ComponentPointer[T]]() QueryOption {
	it := GetIntType[T, TP]()
	return func(q *QueryConfig) {
		q.queryBuddies.Add(it)
	}
}

type QueryIterator struct {
	ctx    *SystemContext
	config QueryConfig
	minSet ComponentSet
}

func (q *QueryIterator) Iter() iter.Seq2[EntityIndex, *EntityInfo] {
	if !q.ctx.constraint.isValid() ||
		len(q.config.queryBuddies) == 0 ||
		q.minSet == nil ||
		q.minSet.Len() == 0 {
		return func(yield func(EntityIndex, *EntityInfo) bool) {}
	}
	stat := q.ctx.info.getOptReporter()
	stat.shapeUsageAdd(q.config.queryBuddies)
	idx := q.minSet.EntityIndexes()
	return func(yield func(EntityIndex, *EntityInfo) bool) {
		var info *EntityInfo
		for _, index := range idx {
			info = q.ctx.world.entities.getByIndex(index)
			if info.compound.IsSubSet(q.config.queryBuddies) {
				yield(index, info)
			}
		}
	}
}
