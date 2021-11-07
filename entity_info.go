package ecs

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type EntityInfo struct {
	world      *ecsWorld
	mu  sync.RWMutex
	components map[reflect.Type]IComponent
	adding     map[reflect.Type]struct{}
	entity     Entity
}

func newEntityInfo(world *ecsWorld) *EntityInfo {
	entity := &EntityInfo{
		world:      world,
		components: make(map[reflect.Type]IComponent),
		adding	: make(map[reflect.Type]struct{}),
		entity:     newEntity(),
	}
	world.addEntity(entity)
	return entity
}

func (e *EntityInfo) Destroy() {
	e.mu.Lock()
	defer e.mu.Unlock()

	var components []IComponent
	for _, c := range e.components {
		components = append(components, c)
	}
	e.Remove(components...)
	e.world.deleteEntity(e)
}

func (e *EntityInfo) Entity() Entity {
	return e.entity
}

func (e *EntityInfo) hashKey() int64 {
	return int64(e.entity)
}

func (e *EntityInfo) HasByType(types ...reflect.Type) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.hasByType(types...)
}

func (e *EntityInfo) Has(components ...IComponent) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.has(components...)
}

func (e *EntityInfo) has(components ...IComponent) bool {
	for _, c := range components {
		_, ok := e.components[c.Type()]
		if !ok {
			_, ok = e.adding[c.Type()]
			if !ok {
				return false
			}
		}
	}
	return true
}

func (e *EntityInfo) hasByType(types ...reflect.Type) bool {
	for _, typ := range types {
		_, ok := e.components[typ]
		if !ok {
			return false
		}
	}
	return true
}

func (e *EntityInfo) Add(components ...IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, c := range components {
		if err := e.addComponent(c); err != nil {
			Log.Error("repeat component:", err)
		}
	}
}

func (e *EntityInfo) Remove(components ...IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, c := range components {
		typ := c.Type()
		if !e.has(c) {
			Log.Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		e.world.deleteComponent(e, c)
	}
}

func (e *EntityInfo) addComponent(com IComponent) error {
	switch com.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		return errors.New("this type of component can not add to entity")
	}
	com.setOwner(e)
	if e.has(com) {
		return fmt.Errorf("repeated component: %s", com.Type().Name())
	}
	e.adding[com.Type()] = Empty
	e.world.addComponent(e, com)
	return nil
}

func (e *EntityInfo) componentAdded(typ reflect.Type, com IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.components[typ] = com
}

func (e *EntityInfo) componentDeleted(typ reflect.Type, com IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.components, typ)
}

func (e *EntityInfo) getComponent(com IComponent) IComponent {
	return e.getComponentByType(com.Type())
}

func (e *EntityInfo) getComponentByType(typ reflect.Type) IComponent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.components[typ]
}
