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
	world      *worldBase
	seq        uint16
	types      map[reflect.Type]uint16
	infos      *SparseArray[uint16, ComponentMetaInfo]
	disposable map[reflect.Type]uint16
	free       map[reflect.Type]uint16
}

func NewComponentMeta(world *worldBase) *componentMeta {
	return &componentMeta{
		world:      world,
		seq:        0,
		types:      make(map[reflect.Type]uint16),
		infos:      NewSparseArray[uint16, ComponentMetaInfo](),
		disposable: map[reflect.Type]uint16{},
		free:       map[reflect.Type]uint16{},
	}
}

func (c *componentMeta) CreateComponentMetaInfo(typ reflect.Type, ct ComponentType) *ComponentMetaInfo {
	if c.world.config.MainThreadCheck {
		c.world.checkMainThread()
	}
	c.seq++
	info := &ComponentMetaInfo{}
	info.it = c.seq
	info.componentType = ct
	info.typ = typ

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

func (c *componentMeta) ComponentMetaInfoPrint() {
	fn := func(m map[reflect.Type]uint16) {
		total := len(m)
		count := 0
		prefix := "│  ├─"
		prefix2 := "│  └─"
		str := ""
		for typ, _ := range m {
			str += " " + typ.Name()
			count++
			if count%5 == 0 {
				if count == total {
					Log.Infof("%s%s", prefix2, str)
				} else {
					Log.Infof("%s%s", prefix, str)
				}
				str = ""
			}
		}
		if str != "" {
			Log.Infof("%s%s", prefix2, str)
			str = ""
		}
	}

	Log.Infof("┌──────────────── # Component Info # ─────────────────")
	Log.Infof("├─ Total: %d", len(c.types))
	fn(c.types)
	Log.Infof("├─ Disposable: %d", len(c.disposable))
	fn(c.disposable)
	Log.Infof("├─ Free: %d", len(c.free))
	fn(c.free)
	Log.Infof("└────────────── # Component Info End # ───────────────")
}
