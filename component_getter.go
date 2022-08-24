package ecs

type ComponentGetter[T ComponentObject] struct {
	permission ComponentPermission
	set        *ComponentSet[T]
}

func NewComponentGetter[T ComponentObject](sys ISystem) *ComponentGetter[T] {
	typ := TypeOf[T]()
	r, isRequire := sys.isRequire(typ)
	if !isRequire {
		return nil
	}
	getter := &ComponentGetter[T]{}
	getter.set = sys.World().getComponents(typ).(*ComponentSet[T])
	getter.permission = r.getPermission()
	if sys.getGetterCache()[typ] == nil {
		sys.getGetterCache()[typ] = getter
	}
	return getter
}

func (c *ComponentGetter[T]) Get(entity Entity) *T {
	var ret *T
	if c.permission == ComponentReadOnly {
		temp := *c.set.getByEntity(entity)
		ret = &temp
	} else {
		ret = c.set.getByEntity(entity)
	}
	return ret
}
