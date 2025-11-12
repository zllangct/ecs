package ecs

type EntitySet struct {
	SparseArray[EntityIndex, EntityInfo]
}

func NewEntitySet() *EntitySet {
	return &EntitySet{
		SparseArray: *NewSparseArray[EntityIndex, EntityInfo](),
	}
}

func (c *EntitySet) Exist(entity Entity) bool {
	index := entity.Index()
	return c.SparseArray.Exist(index)
}

func (c *EntitySet) Get(entity Entity) (*EntityInfo, bool) {
	index := entity.Index()
	info := c.SparseArray.Get(index)
	if info == nil {
		return nil, false
	}
	return info, true
}

func (c *EntitySet) Add(entityInfo EntityInfo) *EntityInfo {
	index := entityInfo.entity.Index()
	return c.SparseArray.Add(index, &entityInfo)
}

func (c *EntitySet) Remove(entity Entity) *EntityInfo {
	index := entity.Index()
	return c.SparseArray.Remove(index)
}

func (c *EntitySet) getByIndex(index EntityIndex) *EntityInfo {
	return c.SparseArray.Get(index)
}
