package ecs

type EntityInfo struct {
	entity   Entity
	compound Compound
}

func (e *EntityInfo) Destroy(world IWorld) {
	for i := 0; i < len(e.compound); i++ {
		world.deleteComponentByIntType(e.entity, e.compound[i])
	}
	// must be last
	world.deleteEntity(e.entity)
}

func (e *EntityInfo) Entity() Entity {
	return e.entity
}

func (e *EntityInfo) Add(world IWorld, components ...IComponent) {
	for _, c := range components {
		if !e.compound.Exist(world.getComponentMetaInfoByType(c.Type()).it) {
			world.addComponent(e.entity, c)
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

func (e *EntityInfo) HasType(world *ecsWorld, components ...IComponent) bool {
	for i := 0; i < len(components); i++ {
		if !e.compound.Exist(world.getComponentMetaInfoByType(components[i].Type()).it) {
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

func (e *EntityInfo) Remove(world IWorld, components ...IComponent) {
	for _, c := range components {
		if e.compound.Exist(world.getComponentMetaInfoByType(c.Type()).it) {
			world.deleteComponent(e.entity, c)
		}
	}
}
