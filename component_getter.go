package ecs

import (
	"reflect"
	"unsafe"
)

// reflect.Type is not comparable, so we need to implement our own
// to replace by SparseArray[K, V] when reflect.Type is comparable
type GetterCache struct {
	indices []reflect.Type
	values  []unsafe.Pointer
}

func NewGetterCache(initCap ...int) *GetterCache {
	cap := 0
	if len(initCap) > 0 {
		cap = initCap[0]
	}
	return &GetterCache{
		indices: make([]reflect.Type, 0, cap),
		values:  make([]unsafe.Pointer, 0, cap),
	}
}

func (g *GetterCache) Add(key reflect.Type, value unsafe.Pointer) {
	g.indices = append(g.indices, key)
	g.values = append(g.values, value)
}

func (g *GetterCache) Remove(key reflect.Type) {
	idx := -1
	for i, t := range g.indices {
		if t == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return
	}
	// remove from indices
	g.indices = append(g.indices[:idx], g.indices[idx+1:]...)
	// remove from values
	g.values = append(g.values[:idx], g.values[idx+1:]...)
}

func (g *GetterCache) Get(key reflect.Type) unsafe.Pointer {
	for i, t := range g.indices {
		if t == key {
			return g.values[i]
		}
	}
	return nil
}

type ComponentGetter[T ComponentObject] struct {
	permission ComponentPermission
	set        *ComponentSet[T]
}

func NewComponentGetter[T ComponentObject](sys ISystem) *ComponentGetter[T] {
	typ := TypeOf[T]()

	r, isRequire := sys.GetRequirements()[typ]
	if !isRequire {
		return nil
	}
	getter := &ComponentGetter[T]{}
	seti := sys.World().getComponentSet(typ)
	if seti == nil {
		return nil
	}
	getter.set = seti.(*ComponentSet[T])
	getter.permission = r.getPermission()
	return getter
}

func (c *ComponentGetter[T]) Get(entity Entity) *T {
	if c.permission == ComponentReadOnly {
		return &(*c.set.getByEntity(entity))
	} else {
		return c.set.getByEntity(entity)
	}
}
