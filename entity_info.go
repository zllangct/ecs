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

func (e *EntityInfo) Add(components ...IComponent) []error {
	for _, c := range components {
		e.world.addComponent(e.entity, c)
		e.compound.Add(e.world.GetComponentMetaInfo(c.Type()).it)
	}
	return nil
}

func (e *EntityInfo) addComponentInternal(it uint16) {
	e.compound.Add(it)
}

func (e *EntityInfo) removeComponentInternal(it uint16) {
	e.compound.Remove(it)
}

func (e *EntityInfo) Remove(components ...IComponent) {
	for _, c := range components {
		e.world.deleteComponent(e.entity, c)
		e.compound.Remove(e.world.GetComponentMetaInfo(c.Type()).it)
	}
}
