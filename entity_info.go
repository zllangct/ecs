package ecs

type EntityInfo struct {
	world    *ecsWorld
	entity   Entity
	compound Compound
}

func (e *EntityInfo) Destroy() {
	e.world.deleteEntity(e.entity)
	for i := 0; i < len(e.compound); i++ {
		e.world.deleteComponentByIntType(e.entity, e.compound[i])
	}
}

func (e *EntityInfo) Entity() Entity {
	return e.entity
}

func (e *EntityInfo) Add(components ...IComponent) {
	if mainThreadDebug {
		checkMainThread()
	}
	for _, c := range components {
		if !e.compound.Exist(e.world.getComponentMetaInfoByType(c.Type()).it) {
			e.world.addComponent(e.entity, c)
		}
	}
}

func (e *EntityInfo) Has(its ...uint16) bool {
	for i := 0; i < len(its); i++ {
		if !e.compound.Exist(its[i]) {
			return false
		}
	}
	return true
}

func (e *EntityInfo) HasType(components ...IComponent) bool {
	for i := 0; i < len(components); i++ {
		if !e.compound.Exist(e.world.getComponentMetaInfoByType(components[i].Type()).it) {
			return false
		}
	}
	return true
}

func (e *EntityInfo) addToCompound(it uint16) {
	e.compound.Add(it)
}

func (e *EntityInfo) removeFromCompound(it uint16) {
	e.compound.Remove(it)
}

func (e *EntityInfo) Remove(components ...IComponent) {
	if mainThreadDebug {
		checkMainThread()
	}
	for _, c := range components {
		if e.compound.Exist(e.world.getComponentMetaInfoByType(c.Type()).it) {
			e.world.deleteComponent(e.entity, c)
		}
	}
}
