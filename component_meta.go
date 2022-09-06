package ecs

import (
	"fmt"
	"reflect"
)

func GetComponentMeta[T ComponentObject](world iWorldBase) *ComponentMetaInfo {
	return world.getComponentMetaInfoByType(TypeOf[T]())
}

type ComponentMetaInfo struct {
	it            uint16
	componentType ComponentType
	o1            uint8
	typ           reflect.Type
}

type componentMeta struct {
	seq        uint16
	types      map[reflect.Type]uint16
	infos      *SparseArray[uint16, ComponentMetaInfo]
	disposable map[reflect.Type]uint16
	free       map[reflect.Type]uint16
}

func NewComponentMeta() *componentMeta {
	return &componentMeta{
		seq:        0,
		types:      make(map[reflect.Type]uint16),
		infos:      NewSparseArray[uint16, ComponentMetaInfo](),
		disposable: map[reflect.Type]uint16{},
		free:       map[reflect.Type]uint16{},
	}
}

func (c *componentMeta) CreateComponentMetaInfo(typ reflect.Type, ct ComponentType) *ComponentMetaInfo {
	if mainThreadDebug {
		checkMainThread()
	}
	c.seq++
	info := &ComponentMetaInfo{}
	info.it = c.seq
	info.componentType = ct
	info.typ = typ

	Log.Debugf("create component meta info: %v, %v", info.it, info.typ.String())

	c.types[typ] = info.it
	info = c.infos.Add(info.it, info)

	if ct&ComponentTypeDisposableMask > 0 {
		c.disposable[typ] = info.it
	}
	if ct&ComponentTypeFreeMask > 0 {
		c.free[typ] = info.it
	}

	return info
}

func (c *componentMeta) GetDisposableTypes() map[reflect.Type]uint16 {
	return c.disposable
}

func (c *componentMeta) GetFreeTypes() map[reflect.Type]uint16 {
	return c.free
}

func (c *componentMeta) Exist(typ reflect.Type) bool {
	_, ok := c.types[typ]
	return ok
}

func (c *componentMeta) GetOrCreateComponentMetaInfo(com IComponent) *ComponentMetaInfo {
	it, ok := c.types[com.Type()]
	if !ok {
		it = c.CreateComponentMetaInfo(com.Type(), com.getComponentType()).it
	}
	return c.infos.Get(it)
}

func (c *componentMeta) GetComponentMetaInfoByIntType(it uint16) *ComponentMetaInfo {
	info := c.infos.Get(it)
	if info == nil {
		panic("must register component first")
	}
	return info
}

func (c *componentMeta) GetComponentMetaInfoByType(typ reflect.Type) *ComponentMetaInfo {
	it, ok := c.types[typ]
	if !ok {
		panic(fmt.Sprintf("must register component %s first", typ.String()))
	}
	return c.infos.Get(it)
}
