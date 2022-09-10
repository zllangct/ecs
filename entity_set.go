package ecs

type EntitySet struct {
	SparseArray[int32, EntityInfo]
}

func NewEntityCollection() *EntitySet {
	return &EntitySet{
		SparseArray: *NewSparseArray[int32, EntityInfo](),
	}
}

func (c *EntitySet) Exist(entity Entity) bool {
	index := entity.ToRealID().index
	return c.SparseArray.Exist(index)
}

func (c *EntitySet) GetEntityInfo(entity Entity) (*EntityInfo, bool) {
	index := entity.ToRealID().index
	info := c.Get(index)
	if info == nil {
		return nil, false
	}
	return info, true
}

func (c *EntitySet) Add(entityInfo EntityInfo) *EntityInfo {
	index := entityInfo.entity.ToRealID().index
	return c.SparseArray.Add(index, &entityInfo)
}

func (c *EntitySet) Remove(entity Entity) *EntityInfo {
	index := entity.ToRealID().index
	return c.SparseArray.Remove(index)
}
