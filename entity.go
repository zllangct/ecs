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
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.hasByType(types...)
}

func (e *Entity) Has(cs ...IComponent) bool {
	e.lock.RLock()
	defer e.lock.RUnlock()

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

func (e *Entity) AddByTemplate(templates ...IComponentTemplate) {
	for _, c := range templates {
		if err := e.addByTemplate(c); err != nil{
			e.world.Logger().Error("repeat component:", err)
		}
	}
}

func (e *Entity) addByTemplate(com IComponentTemplate) error {
	com = com.SetOwner(e)
	typ := com.ComponentType()
	if e.hasByType(typ) {
		return fmt.Errorf("repeated component: %s", typ.Name())
	}
	e.world.ComponentTemplateAttach(e, com)
	return nil
}

func (e *Entity) AddComponent(components ...IComponent) {
	for _, c := range components {
		if err := e.addComponent(c); err != nil{
			e.world.Logger().Error("repeat component:", err)
		}
	}
}

func (e *Entity) addComponent(com IComponent) error {
	if com.Owner() != nil {
		return errors.New("the owner of component is nil")
	}
	if e.has(com) {
		return fmt.Errorf("repeated component: %s", com.Type().Name())
	}
	e.world.ComponentAttach(e, com)
	return nil
}

func (e *Entity) addComponentNoLock(typ reflect.Type, com IComponent){
	e.components[typ] = com
}

func (e *Entity) componentAdded(typ reflect.Type, com IComponent) {
	e.lock.Lock()
	defer e.lock.Unlock()

	e.addComponentNoLock(typ, com)
}

func (e *Entity) RemoveComponent(com ...IComponent) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, c := range com {
		typ := reflect.TypeOf(c)
		if !e.hasByType(typ) {
			e.world.Logger().Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		delete(e.components, typ)
		e.world.ComponentRemove(c.Owner(), c)
	}
}

func (e *Entity) GetComponent(com IComponent) IComponent {
	return e.getComponent(reflect.TypeOf(com).Elem())
}

func (e *Entity) getComponent(typ reflect.Type) IComponent {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.components[typ]
}


