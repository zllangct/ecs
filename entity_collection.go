package ecs

type EntityCollection struct {
	UnorderedCollectionWithID[Entity]
}

func NewEntityCollection() *EntityCollection {
	c := NewUnorderedCollectionWithID[Entity]()
	return &EntityCollection{*c}
}

func (c *EntityCollection) Add(entity Entity) {
	c.UnorderedCollectionWithID.Add(&entity, int64(entity))
}
