package ecs

type EntityTemplate struct {
	Components []Component
}

func (e *EntityTemplate) Instance(world World) *EntityInfo {
	return world.NewEntity(WithComponents(e.Components...))
}

func (e *EntityTemplate) InstanceN(world World, num int) []*EntityInfo {
	entities := make([]*EntityInfo, 0, num)
	for i := 0; i < num; i++ {
		entities = append(entities, e.Instance(world))
	}
	return entities
}
