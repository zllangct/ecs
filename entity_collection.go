package ecs

import (
	runtime2 "runtime"
	"sync"
)

type EntityCollection struct {
	collection []map[uint64]*Entity
	cpuCount int
	lock []sync.RWMutex
}

func NewEntityCollection() *EntityCollection {
	ec:=&EntityCollection{}
	ec.cpuCount = runtime2.NumCPU()
	ec.collection = make([]map[uint64]*Entity,ec.cpuCount,ec.cpuCount)
	ec.lock = make([]sync.RWMutex,ec.cpuCount,ec.cpuCount)
	for index, _ := range ec.collection {
		ec.collection[index] = map[uint64]*Entity{}
		ec.lock[index] = sync.RWMutex{}
	}
	return ec
}

func (p *EntityCollection)Get(id uint64)*Entity {
	hash := id % uint64(p.cpuCount)
	p.lock[hash].RLock()
	defer p.lock[hash].RUnlock()
	return p.collection[hash][id]
}

func (p *EntityCollection)Add(entity *Entity) {
	id:=entity.ID
	hash := id % uint64(p.cpuCount)
	p.lock[hash].Lock()
	p.collection[hash][id]=entity
	p.lock[hash].Unlock()
}

func (p *EntityCollection) Delete(entity *Entity) {
	id:=entity.ID
	p.DeleteByID(id)
}

func (p *EntityCollection) DeleteByID(id uint64) {
	hash := id % uint64(p.cpuCount)
	p.lock[hash].Lock()
	 delete(p.collection[hash],id)
	p.lock[hash].Unlock()
}



