package ecs

import (
	"sync"
)

type EntityCollection struct {
	collection []map[uint64]*Entity
	base       uint64
	locks      []sync.RWMutex
}

func NewEntityCollection(k int) *EntityCollection {
	ec := &EntityCollection{}

	for i := 1; ; i++ {
		if c := uint64(1 << i); uint64(k) < c {
			ec.base = c - 1
			break
		}
	}

	ec.collection = make([]map[uint64]*Entity, ec.base+1)
	ec.locks = make([]sync.RWMutex, ec.base+1)
	for index := range ec.collection {
		ec.collection[index] = map[uint64]*Entity{}
		ec.locks[index] = sync.RWMutex{}
	}
	return ec
}

func (p *EntityCollection) get(id uint64) *Entity {
	hash := id & p.base
	p.locks[hash].RLock()
	defer p.locks[hash].RUnlock()
	return p.collection[hash][id]
}

func (p *EntityCollection) add(entity *Entity) {
	hash := entity.id & p.base
	println(entity.id, entity.id%(p.base+1), entity.id&p.base)
	p.locks[hash].Lock()
	p.collection[hash][entity.id] = entity
	p.locks[hash].Unlock()
}

func (p *EntityCollection) delete(entity *Entity) {
	p.deleteByID(entity.id)
}

func (p *EntityCollection) deleteByID(id uint64) {
	hash := id & p.base
	p.locks[hash].Lock()
	delete(p.collection[hash], id)
	p.locks[hash].Unlock()
}
