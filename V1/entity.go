package main

import (
	"errors"
	"reflect"
	"sync"
)

var(
	ErrComponentInvalid = errors.New("component invalid")
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
	defer p.Unlock()
	for _, c := range com {
		if c.GetOwner() != nil {
			panic(ErrComponentInvalid)
		}
		p.components = append(p.components, c)
		c.setOwner(p)
		p.runtime.ComponentAttach(c)
	}
}

func (p *Entity) RemoveComponent(com ... IComponent)  {
	p.Lock()
	defer p.Unlock()
	for _, c := range com {
		for i := 0; i< len(p.components);i++  {
			if reflect.TypeOf(p.components[i]) == reflect.TypeOf(c){
				p.components[i] = p.components[len(p.components)-1]
				p.components = p.components[:len(p.components)-1]
				p.runtime.ComponentRemove(c)
				break
			}
		}
	}
}

func (p *Entity) GetComponent(com IComponent) IComponent {
	typ:=reflect.TypeOf(com)
	p.RLock()
	for _, value := range p.components {
		if reflect.TypeOf(value) == typ {
			return value
		}
	}
	p.RUnlock()
	return nil
}