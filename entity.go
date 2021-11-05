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

func newEntity(world *World) *Entity {
	entity := &Entity{
		world:      world,
		components: make(map[reflect.Type]IComponent),
		id:         UniqueEntityID(),
	}
	world.addEntity(entity)
	return entity
}

func (e *Entity) Destroy() {
	var components []IComponent
	for _, c := range e.components {
		components = append(components, c)
	}
	e.Remove(components...)
	e.world.deleteEntity(e)
}

func (e *Entity) GetID() int64 {
	return e.id
}

func (e *Entity) HasByType(types ...reflect.Type) bool {
	return e.hasByType(types...)
}

func (e *Entity) Has(components ...IComponent) bool {
	return e.has(components...)
}

func (e *Entity) has(components ...IComponent) bool {
	for _, c := range components {
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
		if err := e.addComponent(c); err != nil {
			Log.Error("repeat component:", err)
		}
	}
}

func (e *Entity) Remove(components ...IComponent) {
	for _, c := range components {
		typ := c.Type()
		if !e.has(c) {
			Log.Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		e.world.components.TempTemplateOperate(e, c.Template(), CollectionOperateDelete)
	}
}

func (e *Entity) addComponent(com IComponent) error {
	com.setOwner(e)
	if e.has(com) {
		return fmt.Errorf("repeated component: %s", com.Type().Name())
	}
	e.world.components.TempTemplateOperate(e, com.Template(), CollectionOperateAdd)
	return nil
}

func (e *Entity) componentAdded(typ reflect.Type, com IComponent) {
	e.components[typ] = com
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
