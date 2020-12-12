package ecs

import "unsafe"

type ContainerWithId struct {
	Container
	id2idx map[uint64]int
	idx2id map[int]uint64
}

func NewContainerWithId(size uintptr) *ContainerWithId {
	return &ContainerWithId{
		Container: Container{
			buf:  make([]byte, 0, size),
			len:  -1,
			unit: size,
		},
		idx2id: map[int]uint64{},
		id2idx: map[uint64]int{},
	}
}

func (p *ContainerWithId) Add(pointer unsafe.Pointer, id ...uint64) (int, unsafe.Pointer) {
	if len(id) > 0 {
		_, ok := p.id2idx[id[0]]
		if ok {
			return 0, nil
		}
	}
	idx, ptr := p.Container.Add(pointer)
	if len(id) > 0 {
		p.id2idx[id[0]] = p.len
		p.idx2id[p.len] = id[0]
	}
	return idx, ptr
}

func (p *ContainerWithId) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	p.id2idx[p.idx2id[p.len]] = idx
	delete(p.id2idx, p.idx2id[idx])
	p.idx2id[idx] = p.idx2id[p.len]
	delete(p.idx2id, p.len)

	p.Container.Remove(idx)
}

func (p *ContainerWithId) RemoveById(id uint64) {
	idx, ok := p.id2idx[id]
	if !ok {
		return
	}
	p.Remove(idx)
}

func (p *ContainerWithId) Get(idx int) unsafe.Pointer {
	return p.Container.Get(idx)
}

func (p *ContainerWithId) GetById(id uint64) unsafe.Pointer {
	idx, ok := p.id2idx[id]
	if !ok {
		return nil
	}
	return p.Container.Get(idx)
}

func (p *ContainerWithId) GetId(idx int) uint64 {
	if id, ok := p.idx2id[idx]; ok {
		return id
	}
	return 0
}
