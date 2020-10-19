package ecs

import (
	"reflect"
	"sync"
)

type Entity struct {
	sync.RWMutex
	//private
	runtime    *Runtime
	components map[reflect.Type]IComponent
	//public
	id uint64
}

func NewEntity(runtime *Runtime) *Entity {
	entity := &Entity{
		runtime:    runtime,
		components: make(map[reflect.Type]IComponent),
		id:         UniqueID(),
	}
	runtime.AddEntity(entity)
	return entity
}

func (p *Entity) Destroy() {
	p.runtime.DeleteEntity(p)
	for _, c := range p.components {
		p.runtime.ComponentRemove(c)
	}
}

func (p *Entity) ID() uint64 {
	return p.id
}

func (p *Entity) Has(typ reflect.Type) bool {
	p.RLock()
	defer p.RUnlock()
	_, ok := p.components[typ]
	return ok
}

func (p *Entity) AddComponent(com ...IComponent) {
	p.Lock()
	defer p.Unlock()
	for _, c := range com {
		if c.GetOwner() != nil {
			continue
		}
		p.components[reflect.TypeOf(c)] = c
		c.setOwner(p)
		p.runtime.ComponentAttach(c)
	}
}

func (p *Entity) RemoveComponent(com ...IComponent) {
	p.Lock()
	defer p.Unlock()
	for _, c := range com {
		delete(p.components, reflect.TypeOf(c))
		p.runtime.ComponentRemove(c)
	}
}

func (p *Entity) GetComponent(com IComponent) IComponent {
	p.RLock()
	defer p.RUnlock()
	return p.components[reflect.TypeOf(com)]
}
