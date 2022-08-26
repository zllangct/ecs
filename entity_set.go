package ecs

type EntitySet struct {
	indices  []int64
	entities *UnorderedCollection[Entity]
	refs     *UnorderedCollection[EntityInfo]
}

func NewEntityCollection() *EntitySet {
	indices := make([]int64, 0)
	c := NewUnorderedCollection[Entity]()
	refs := NewUnorderedCollection[EntityInfo]()
	return &EntitySet{indices: indices, entities: c, refs: refs}
}

func (c *EntitySet) Exist(entity Entity) bool {
	realID := entity.ToRealID()
	if c.indices[realID.index] <= 0 {
		return false
	}
	return true
}

func (c *EntitySet) GetEntityInfo(entity Entity) (*EntityInfo, bool) {
	realID := entity.ToRealID()
	if c.indices[realID.index] <= 0 {
		return nil, false
	}
	info := c.refs.Get(int64(realID.index))
	if info == nil {
		panic("index mismatch")
	}
	return info, true
}

func (c *EntitySet) Len() int {
	return c.entities.Len()
}

func (c *EntitySet) Add(entity Entity, entityInfo *EntityInfo) bool {
	realID := entity.ToRealID()
	if c.indices[realID.index] != -1 {
		return false
	}

	_, idx := c.entities.Add(&entity)
	_, refIdx := c.refs.Add(entityInfo)
	if idx != refIdx {
		panic("index mismatch")
	}
	c.indices[realID.index] = idx
	return true
}

func (c *EntitySet) Remove(entity Entity) bool {
	realID := entity.ToRealID()
	if c.indices[realID.index] <= 0 {
		return false
	}

	_, oldIdx, newIdx := c.entities.Remove(int64(realID.index))
	_, oldRefIdx, newRefIdx := c.refs.Remove(int64(realID.index))
	if oldIdx != oldRefIdx || newIdx != newRefIdx {
		panic("index mismatch")
	}
	changed := c.entities.Get(newIdx).ToRealID()
	c.indices[changed.index] = newIdx

	return true
}
