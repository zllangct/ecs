package ecs

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type EntityInfo struct {
	world      *ecsWorld
	mu         sync.RWMutex
	components map[reflect.Type]IComponent
	adding     map[reflect.Type]struct{}
	once       map[reflect.Type]IComponent
	entity     Entity
}

func newEntityInfo(world *ecsWorld) *EntityInfo {
	entity := &EntityInfo{
		world:      world,
		components: make(map[reflect.Type]IComponent),
		adding:     make(map[reflect.Type]struct{}),
		once:       make(map[reflect.Type]IComponent),
		entity:     newEntity(),
	}
	world.addEntity(entity)
	return entity
}

func (e *EntityInfo) Destroy() {
	e.mu.Lock()
	defer e.mu.Unlock()

	var components []IComponent = make([]IComponent, 0, len(e.components))
	for _, c := range e.components {
		components = append(components, c)
	}
	e.remove(components...)
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
	ok := false
	for _, c := range components {
		switch c.getComponentType() {
		case ComponentTypeDisposable:
			_, ok = e.once[c.Type()]
		case ComponentTypeNormal:
			_, ok = e.components[c.Type()]
		}
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
			_, ok = e.once[typ]
			if !ok {
				_, ok = e.adding[typ]
				if !ok {
					return false
				}
			}
		}
	}
	return true
}

func (e *EntityInfo) Add(components ...IComponent) []error {
	e.mu.Lock()
	defer e.mu.Unlock()
	var errors []error
	for _, c := range components {
		if err := e.addComponent(c); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

func (e *EntityInfo) remove(components ...IComponent) {
	for _, c := range components {
		//typ := c.Type()
		if !e.has(c) {
			continue
		}
		if c.getComponentType() == ComponentTypeNormal {
			e.world.deleteComponent(e, c)
		}
	}
}

func (e *EntityInfo) Remove(components ...IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.remove(components...)
}

func (e *EntityInfo) addComponent(com IComponent) error {
	ct := com.getComponentType()
	switch ct {
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

	delete(e.adding, typ)

	if com.getComponentType() == ComponentTypeDisposable {
		e.once[typ] = com
	} else {
		e.components[typ] = com
	}
}

func (e *EntityInfo) componentDeleted(typ reflect.Type, comType ComponentType) {
	e.mu.Lock()
	defer e.mu.Unlock()

	switch comType {
	case ComponentTypeNormal:
		delete(e.components, typ)
	case ComponentTypeDisposable:
		delete(e.once, typ)
	}
}

func (e *EntityInfo) getComponent(com IComponent) IComponent {
	return e.getComponentByType(com.Type())
}

func (e *EntityInfo) getComponentByType(typ reflect.Type) IComponent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.components[typ]
}

func (e *EntityInfo) clearDisposable() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.once) > 0 {
		e.once = make(map[reflect.Type]IComponent)
	}
}
