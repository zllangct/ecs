package ecs

type EntityCollection struct {
	entities *UnorderedCollectionWithID[Entity]
	refs     *UnorderedCollectionWithID[EntityInfo]
}

func NewEntityCollection() *EntityCollection {
	c := NewUnorderedCollectionWithID[Entity]()
	refs := NewUnorderedCollectionWithID[EntityInfo]()
	return &EntityCollection{entities: c, refs: refs}
}

func (c *EntityCollection) Add(entity Entity) {
	c.entities.Add(&entity, int64(entity))
}

func (c *EntityCollection) Remove(entity Entity) {
	c.entities.Remove(int64(entity))
}
