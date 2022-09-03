package ecs

import (
	"reflect"
	"sync"
)

var ComponentMeta = &componentMeta{
	seq:     0,
	types:   map[reflect.Type]uint16{},
	it2Type: map[uint16]reflect.Type{},
}

func GetComponentMeta[T ComponentObject]() ComponentMetaInfo {
	return ComponentMeta.GetComponentMetaInfo(TypeOf[T]())
}

type ComponentMetaInfo struct {
	it uint16
}

type componentMeta struct {
	mu      sync.RWMutex
	seq     uint16
	types   map[reflect.Type]uint16
	it2Type map[uint16]reflect.Type
}

func (c *componentMeta) ConvertToType(it uint16) reflect.Type {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.it2Type[it]
}

func (c *componentMeta) GetComponentMetaInfo(typ reflect.Type) ComponentMetaInfo {
	it := uint16(0)
	c.mu.RLock()
	if v, ok := c.types[typ]; ok {
		it = v
		c.mu.RUnlock()
	} else {
		c.mu.RUnlock()
		c.mu.Lock()
		if v, ok = c.types[typ]; ok {
			it = v
		} else {
			c.seq++
			it = c.seq
			c.types[typ] = it
			c.it2Type[it] = typ
		}
		c.mu.Unlock()
	}

	return ComponentMetaInfo{
		it: it,
	}
}
