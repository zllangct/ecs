package ecs

import (
	"sync"
)

type EntityCollection struct {
	collection []map[int64]*Entity
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

	ec.collection = make([]map[int64]*Entity, ec.base+1)
	ec.locks = make([]sync.RWMutex, ec.base+1)
	for index := range ec.collection {
		ec.collection[index] = map[int64]*Entity{}
		ec.locks[index] = sync.RWMutex{}
	}
	return ec
}

func (p *EntityCollection) get(id int64) *Entity {
	hash := id & p.base

	p.locks[hash].RLock()
	defer p.locks[hash].RUnlock()

	return p.collection[hash][id]
}

func (p *EntityCollection) add(entity *Entity) {
	hash := entity.id & p.base

	p.locks[hash].Lock()
	defer p.locks[hash].Unlock()

	p.collection[hash][entity.id] = entity
}

func (p *EntityCollection) delete(entity *Entity) {
	p.deleteByID(entity.id)
}

func (p *EntityCollection) deleteByID(id int64) {
	hash := id & p.base

	p.locks[hash].Lock()
	defer p.locks[hash].Unlock()

	delete(p.collection[hash], id)
}
