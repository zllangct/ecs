package ecs

import (
	runtime2 "runtime"
	"sync"
)

var _2n = []uint64{16, 32, 64, 128, 256, 512, 1024}
var base uint64

type EntityCollection struct {
	collection []map[uint64]*Entity
	cpuCount   uint64
	locks      []sync.RWMutex
}

func NewEntityCollection() *EntityCollection {
	ec := &EntityCollection{}
	numCpu := runtime2.NumCPU()

	ec.cpuCount = _2n[len(_2n)-1]
	for _, i := range _2n {
		if uint64(numCpu*2) < i {
			ec.cpuCount = i
			break
		}
	}
	base = ec.cpuCount - 1

	ec.collection = make([]map[uint64]*Entity, ec.cpuCount, ec.cpuCount)
	ec.locks = make([]sync.RWMutex, ec.cpuCount, ec.cpuCount)
	for index, _ := range ec.collection {
		ec.collection[index] = map[uint64]*Entity{}
		ec.locks[index] = sync.RWMutex{}
	}
	return ec
}

func (p *EntityCollection) get(id uint64) *Entity {
	hash := id & base
	p.locks[hash].RLock()
	defer p.locks[hash].RUnlock()
	return p.collection[hash][id]
}

func (p *EntityCollection) add(entity *Entity) {
	id := entity.ID
	hash := id & base
	p.locks[hash].Lock()
	p.collection[hash][id] = entity
	p.locks[hash].Unlock()
}

func (p *EntityCollection) delete(entity *Entity) {
	id := entity.ID
	p.deleteByID(id)
}

func (p *EntityCollection) deleteByID(id uint64) {
	hash := id & base
	p.locks[hash].Lock()
	delete(p.collection[hash], id)
	p.locks[hash].Unlock()
}
