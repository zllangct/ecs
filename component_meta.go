package ecs

import (
	"reflect"
	"sync"
)

func GetComponentMeta[T ComponentObject](world IWorld) *ComponentMetaInfo {
	return world.GetComponentMetaInfo(TypeOf[T]())
}

type ComponentMetaInfo struct {
	it            uint16
	componentType ComponentType
	o1            uint8
	typ           reflect.Type
}

type componentMeta struct {
	mu    sync.RWMutex
	seq   uint16
	types map[reflect.Type]uint16
	infos *SparseArray[uint16, ComponentMetaInfo]
}

func NewComponentMeta() *componentMeta {
	return &componentMeta{
		seq:   0,
		types: make(map[reflect.Type]uint16),
		infos: NewSparseArray[uint16, ComponentMetaInfo](),
	}
}

func (c *componentMeta) CreateComponentMetaInfo(typ reflect.Type, ct ComponentType) ComponentMetaInfo {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.seq++
	info := ComponentMetaInfo{}
	info.it = c.seq
	info.componentType = ct
	info.typ = typ

	c.types[typ] = info.it
	c.infos.Add(info.it, &info)
	return info
}

func (c *componentMeta) Exist(typ reflect.Type) bool {
	_, ok := c.types[typ]
	return ok
}

func (c *componentMeta) GetComponentMetaInfo(it uint16) *ComponentMetaInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.infos.Get(it)
}

func (c *componentMeta) GetComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.infos.Get(c.types[typ])
}
