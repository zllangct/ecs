package ecs

type EntitySet struct {
	SparseArray[EntityIndex, EntityInfo]
}

func NewEntitySet() *EntitySet {
	return &EntitySet{
		SparseArray: *NewSparseArray[EntityIndex, EntityInfo](),
	}
}

func (es *EntitySet) Exist(entity Entity) bool {
	index := entity.Index()
	return es.SparseArray.Exist(index)
}

func (es *EntitySet) Get(entity Entity) (*EntityInfo, bool) {
	index := entity.Index()
	info := es.SparseArray.Get(index)
	if info == nil {
		return nil, false
	}
	return info, true
}

func (es *EntitySet) Add(entityInfo EntityInfo) *EntityInfo {
	index := entityInfo.entity.Index()
	return es.SparseArray.Add(index, &entityInfo)
}

func (es *EntitySet) Remove(entity Entity) *EntityInfo {
	index := entity.Index()
	return es.SparseArray.Remove(index)
}

func (es *EntitySet) getByIndex(index EntityIndex) *EntityInfo {
	return es.SparseArray.Get(index)
}
