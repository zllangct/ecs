package ecs

type Entity int64

func newEntity() Entity {
	return Entity(UniqueID())
}

func (e Entity) ToInt64() int64 {
	return int64(e)
}
