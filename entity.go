package ecs

import (
	"errors"
	"reflect"
	"sync"
)

type Entity struct {
	lock sync.RWMutex
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
	for _, c := range p.components {
		p.runtime.ComponentRemove(c.GetOwner(), c)
	}
	p.runtime.DeleteEntity(p)
}

func (p *Entity) ID() uint64 {
	return p.id
}

func (p *Entity) Has(types ...reflect.Type) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.has(types...)
}

func (p *Entity) has(types ...reflect.Type) bool {
	for _, typ := range types {
		_, ok := p.components[typ]
		if !ok {
			return false
		}
	}
	return true
}

func (p *Entity) AddComponent(com ...IComponent) {
	for _, c := range com {
		p.lock.Lock()
		p.addComponent(c)
		p.lock.Unlock()
	}
}

func (p *Entity) addComponent(com IComponent) {
	if com.GetOwner() != nil {
		return
	}
	typ := com.GetRealType()
	if p.has(typ) {
		p.runtime.Error("repeat component:", typ.Name())
		return
	}
	p.runtime.ComponentAttach(p, com)
}

func (p *Entity) componentAdded(typ reflect.Type, com IComponent) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.components[typ] = com
}

func (p *Entity) RemoveComponent(com ...IComponent) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, c := range com {
		typ := reflect.TypeOf(c)
		if !p.Has(typ) {
			p.runtime.Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		delete(p.components, typ)
		p.runtime.ComponentRemove(c.GetOwner(), c)
	}
}

func (p *Entity) GetComponent(com IComponentType) IComponent {
	p.lock.RLock()
	defer p.lock.RUnlock()

	return p.components[reflect.TypeOf(com)]
}
