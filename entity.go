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
	runtime    *World
	components map[reflect.Type]IComponent
	//public
	id uint64
}

func NewEntity(world *World) *Entity {
	entity := &Entity{
		runtime:    world,
		components: make(map[reflect.Type]IComponent),
		id:         UniqueID(),
	}
	world.AddEntity(entity)
	return entity
}

func (e *Entity) Destroy() {
	for _, c := range e.components {
		e.runtime.ComponentRemove(c.GetOwner(), c)
	}
	e.runtime.DeleteEntity(e)
}

func (e *Entity) ID() uint64 {
	return e.id
}

func (e *Entity) Has(types ...reflect.Type) bool {
	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.has(types...)
}

func (e *Entity) has(types ...reflect.Type) bool {
	for _, typ := range types {
		_, ok := e.components[typ]
		if !ok {
			return false
		}
	}
	return true
}

func (e *Entity) AddComponent(components ...IComponent) {
	for _, c := range components {
		if err := e.addComponent(c); err != nil{
			e.runtime.Error("repeat component:", err)
		}
	}
}

func (e *Entity) addComponent(com IComponent) error {
	if com.GetOwner() != nil {
		return errors.New("the owner of component is nil")
	}
	typ := com.GetType()
	if e.has(typ) {
		return fmt.Errorf("repeated component: %s", typ.Name())
	}
	e.runtime.ComponentAttach(e, com)
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
		if !e.Has(typ) {
			e.runtime.Error(errors.New("repeat component:" + typ.Name()))
			continue
		}
		delete(e.components, typ)
		e.runtime.ComponentRemove(c.GetOwner(), c)
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

func AttachTo[T IComponent](e *Entity, com ... *T) {
	var ins T
	if len(com) == 0 {
		ins = *new(T)
		e.AddComponent(ins)
	} else {
		e.AddComponent(com...)
	}
}

func GetComponentFrom[T IComponent](e *Entity) IComponent{
	return e.GetComponent(reflect.TypeOf(*new(T)))
}


