package ecs

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Entity struct {
	lock sync.RWMutex
	//private
	world      *World
	components map[reflect.Type]IComponent
	//public
	id int64
}

func NewEntity(world *World) *Entity {
	entity := &Entity{
		world:      world,
		components: make(map[reflect.Type]IComponent),
		id:         UniqueID(),
	}
	world.AddEntity(entity)
	return entity
}

func (e *Entity) Destroy() {
	for _, c := range e.components {
		e.world.ComponentRemove(c.Owner(), c)
	}
	e.world.DeleteEntity(e)
}

func (e *Entity) ID() int64 {
	return e.id
}

func (e *Entity) HasByType(types ...reflect.Type) bool {
	return e.hasByType(types...)
}

func (e *Entity) Has(cs ...IComponent) bool {
	return e.has(cs...)
}

func (e *Entity) has(cs ...IComponent) bool {
	for _, c := range cs {
		_, ok := e.components[c.Type()]
		if !ok {
			return false
		}
	}
	return true
}

func (e *Entity) hasByType(types ...reflect.Type) bool {
	for _, typ := range types {
		_, ok := e.components[typ]
		if !ok {
			return false
		}
	}
	return true
}

func (e *Entity) Add(components ...IComponent) {
	for _, c := range components {
		if err := e.addComponent(c); err != nil{
			e.world.Logger().Error("repeat component:", err)
		}
	}
}

func (e *Entity) addComponent(com IComponent) error {
	com.setOwner(e)
	if e.has(com) {
		return fmt.Errorf("repeated component: %s", com.Type().Name())
	}
	e.world.ComponentAttach(e, com)
	return nil
}

func (e *Entity) componentAdded(typ reflect.Type, com IComponent) {
	e.components[typ] = com
}

func (e *Entity) Remove(com ...IComponent) {
	for _, c := range com {
		typ := c.Type()
		if !e.has(c) {
			e.world.Logger().Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		e.world.ComponentRemove(c.Owner(), c)
	}
}

func (e *Entity) componentDeleted(typ reflect.Type, com IComponent) {
	delete(e.components, typ)
}

func (e *Entity) getComponent(com IComponent) IComponent {
	return e.getComponentByType(com.Type())
}

func (e *Entity) getComponentByType(typ reflect.Type) IComponent {
	return e.components[typ]
}


