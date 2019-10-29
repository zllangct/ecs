package main

import (
	"reflect"
	"sync"
)

type Entity struct {
	sync.RWMutex
	//private
	runtime *Runtime
	components []IComponent
	//public
	ID uint64
}

func NewEntity(runtime *Runtime)*Entity  {
	entity:= &Entity{
		runtime:    runtime,
		components: make([]IComponent,0),
		ID:         UniqueID(),
	}
	runtime.AddEntity(entity)
	return entity
}

func (p *Entity) Destroy()  {
	p.runtime.DeleteEntity(p)
}

func (p *Entity) Has(typ reflect.Type) bool {
	p.RLock()
	for _, value := range p.components {
		if reflect.TypeOf(value) == typ {
			return true
		}
	}
	p.RUnlock()

	return false
}

func (p *Entity) AddComponent(com ... IComponent)  {
	p.Lock()
	for _, c := range com {
		p.components = append(p.components, c)
		c.setOwner(p)
	}
	p.Unlock()
}

func (p *Entity) GetComponent(typ reflect.Type) interface{} {
	p.RLock()
	for _, value := range p.components {
		if reflect.TypeOf(value) == typ {
			return value
		}
	}
	p.RUnlock()
	return nil
}