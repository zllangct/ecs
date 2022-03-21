package ecs

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type TypeList []reflect.Type

func newTypeList(cap int) TypeList {
	return make(TypeList, 0, cap)
}

func (tl TypeList) Contains(t reflect.Type) bool {
	for _, t2 := range tl {
		if t2 == t {
			return true
		}
	}
	return false
}

func (tl TypeList) Find(t reflect.Type) (int, bool) {
	for i, t2 := range tl {
		if t2 == t {
			return i, true
		}
	}
	return 0, false
}

func (tl *TypeList) Remove(t reflect.Type) {
	i, ok := tl.Find(t)
	if !ok {
		return
	}
	//*tl = append((*tl)[:i], (*tl)[i+1:]...)
	(*tl)[i], (*tl)[len(*tl)-1] = (*tl)[len(*tl)-1], (*tl)[i]
	*tl = (*tl)[:len(*tl)-1]
}

func (tl *TypeList) Append(t ...reflect.Type) {
	/* todo
	大量分配对象, TypeList 频繁修改, 考虑链表, map的插入效率,find、delete效率高
	*/
	*tl = append(*tl, t...)
}

type EntityInfo struct {
	world      *ecsWorld
	mu         sync.RWMutex
	components map[reflect.Type]IComponent
	adding     TypeList
	removing   TypeList
	entity     Entity
}

func newEntityInfo(world *ecsWorld) *EntityInfo {
	entity := &EntityInfo{
		world:      world,
		components: make(map[reflect.Type]IComponent),
		adding:     newTypeList(3),
		removing:   newTypeList(3),
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
	for _, c := range components {
		e.deleteComponent(c)
	}
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
	has := true
	for _, c := range components {
		_, has = e.components[c.Type()]
		if !has {
			return false
		}
	}
	return has
}

func (e *EntityInfo) hasByType(types ...reflect.Type) bool {
	has := true
	for _, typ := range types {
		_, has := e.components[typ]
		if !has {
			return false
		}
	}
	return has
}

func (e *EntityInfo) Add(components ...IComponent) []error {
	e.mu.Lock()
	defer e.mu.Unlock()

	var errs []error
	for _, c := range components {
		if err := e.addComponent(c); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EntityInfo) Remove(components ...IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, c := range components {
		if !e.has(c) {
			continue
		}
		e.deleteComponent(c)
	}
}

func (e *EntityInfo) addComponent(com IComponent) error {
	ct := com.getComponentType()
	switch ct {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		return errors.New("this type of component can not add to entity")
	}
	canAdd := true
	typ := com.Type()
	if e.has(com) {
		if _, ok := e.removing.Find(typ); !ok {
			canAdd = false
		}
	}
	if !canAdd {
		return errors.New(fmt.Sprintf("repeated component: %s", typ.Name()))
	}
	_, ok := e.adding.Find(typ)
	if ok {
		return errors.New(fmt.Sprintf("repeated component: %s", typ.Name()))
	}

	//Log.Info("add component: ", com.Type().Name())

	e.adding.Append(typ)
	com.setOwner(e)

	if ct == ComponentTypeDisposable {
		e.deleteComponent(com)
	}
	e.world.addComponent(e, com)

	return nil
}

func (e *EntityInfo) deleteComponent(com IComponent) {
	//Log.Info("delete component: ", com.Type().Name())
	e.removing.Append(com.Type())
	e.world.deleteComponent(e, com)
}

func (e *EntityInfo) componentAdded(typ reflect.Type, com IComponent) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.adding.Remove(typ)

	e.components[typ] = com
	if com.getComponentType() == ComponentTypeDisposable {
		e.removing.Append(typ)
	}

	//Log.Info("component added: ", typ.Name())
}

func (e *EntityInfo) componentDeleted(typ reflect.Type, comType ComponentType) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.removing.Remove(typ)

	delete(e.components, typ)
	//Log.Info("component removed: ", typ.Name())
}

func (e *EntityInfo) getComponent(com IComponent) IComponent {
	return e.getComponentByType(com.Type())
}

func (e *EntityInfo) getComponentByType(typ reflect.Type) IComponent {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return e.components[typ]
}
