package ecs

type EntityCollection struct {
	buckets []*Map[Entity, *EntityInfo]
	bucket  int64
}

func NewEntityCollection(k int) *EntityCollection {
	ec := &EntityCollection{}

	for i := 1; ; i++ {
		if c := int64(1 << i); int64(k) < c {
			ec.bucket = c - 1
			break
		}
	}

	ec.buckets = make([]*Map[Entity, *EntityInfo], ec.bucket+1)
	for index := range ec.buckets {
		ec.buckets[index] = &Map[Entity, *EntityInfo]{}
	}
	return ec
}

func (p *EntityCollection) getInfo(entity Entity) *EntityInfo {
	hash := int64(entity) & p.bucket

	info, ok := p.buckets[hash].Load(entity)
	if !ok {
		return nil
	}
	return info
}

func (p *EntityCollection) add(info *EntityInfo) {
	hash := info.hashKey() & p.bucket

	p.buckets[hash].Store(info.entity, info)
}

func (p *EntityCollection) delete(entity Entity) {
	hash := int64(entity) & p.bucket

	p.buckets[hash].Delete(entity)
}

func (p *EntityCollection) getBuckets() []*Map[Entity, *EntityInfo] {
	return p.buckets
}
