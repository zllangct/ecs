package ecs

import (
	"sync"
)

type EntityCollection struct {
	collection []map[Entity]*EntityInfo
	base       int64
	locks      []sync.RWMutex
}

func NewEntityCollection(k int) *EntityCollection {
	ec := &EntityCollection{}

	for i := 1; ; i++ {
		if c := int64(1 << i); int64(k) < c {
			ec.base = c - 1
			break
		}
	}

	ec.collection = make([]map[Entity]*EntityInfo, ec.base+1)
	ec.locks = make([]sync.RWMutex, ec.base+1)
	for index := range ec.collection {
		ec.collection[index] = map[Entity]*EntityInfo{}
		ec.locks[index] = sync.RWMutex{}
	}
	return ec
}

func (p *EntityCollection) getInfo(entity Entity) *EntityInfo {
	hash := int64(entity) & p.base

	p.locks[hash].RLock()
	defer p.locks[hash].RUnlock()

	return p.collection[hash][entity]
}

func (p *EntityCollection) add(entity *EntityInfo) {
	hash := entity.hashKey() & p.base

	p.locks[hash].Lock()
	defer p.locks[hash].Unlock()

	p.collection[hash][entity.entity] = entity
}

func (p *EntityCollection) delete(entity Entity) {
	hash := int64(entity) & p.base

	p.locks[hash].Lock()
	defer p.locks[hash].Unlock()

	delete(p.collection[hash], entity)
}
