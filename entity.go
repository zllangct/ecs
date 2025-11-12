package ecs

import (
	"sort"
	"unsafe"
)

type Entity int64

func (e Entity) ToInt64() int64 {
	return int64(e)
}

func (e Entity) toReuseID() ReuseID {
	return *(*ReuseID)(unsafe.Pointer(&e))
}

func (e Entity) Index() EntityIndex {
	return (*(*ReuseID)(unsafe.Pointer(&e))).index
}

type EntityIndex int32

type ReuseID struct {
	index EntityIndex
	reuse int32
}

func (r *ReuseID) ToInt64() int64 {
	return *(*int64)(unsafe.Pointer(r))
}

func (r *ReuseID) ToEntity() Entity {
	return *(*Entity)(unsafe.Pointer(r))
}

type EntityIDGenerator struct {
	ids     []ReuseID
	free    EntityIndex
	pending EntityIndex
	len     int32

	removeDelay []ReuseID
	delayFree   int32
	delayCap    int32
}

func NewEntityIDGenerator(initSize int, delayCap int) *EntityIDGenerator {
	g := &EntityIDGenerator{}
	g.ids = make([]ReuseID, initSize)
	for i := 0; i < len(g.ids); i++ {
		g.ids[i].index = EntityIndex(i + 1)
	}
	g.free = 1
	g.pending = EntityIndex(initSize)
	g.len = 0
	g.removeDelay = make([]ReuseID, delayCap)
	g.delayCap = int32(delayCap)
	g.delayFree = 0
	return g
}

func (e *EntityIDGenerator) NewID() Entity {
	id := ReuseID{}
	if e.free == e.pending {
		e.ids = append(e.ids, ReuseID{index: e.free, reuse: 0})
		id = e.ids[e.pending]
		e.free++
		e.pending++
	} else {
		next := e.ids[e.free].index
		e.ids[e.free].index = e.free
		id = e.ids[e.free]
		e.free = next
	}
	e.len++
	return id.ToEntity()
}

func (e *EntityIDGenerator) FreeID(entity Entity) {
	e.len--

	realID := entity.toReuseID()
	e.ids[realID.index].index = -1
	e.ids[realID.index].reuse++

	e.removeDelay[e.delayFree] = realID
	e.delayFree++
	if e.delayFree >= e.delayCap {
		e.delayFlush()
	}
	if e.pending > 1024 && e.pending < EntityIndex(len(e.ids))/2 {
		e.ids = e.ids[:e.len*5/8]
	}
}

func (e *EntityIDGenerator) delayFlush() {
	sort.Slice(e.removeDelay, func(i, j int) bool {
		return e.removeDelay[i].index < e.removeDelay[j].index
	})
	lastFree := e.free
	nextFree := e.free
	if e.free < e.pending {
		nextFree = e.ids[e.free].index
	}
	for i := int32(0); i < e.delayFree; i++ {
		tempID := e.removeDelay[i]
		if tempID.index < lastFree {
			e.ids[tempID.index].index = lastFree
			e.free = tempID.index
			nextFree = lastFree
			lastFree = tempID.index
			continue
		} else {
			for {
				if tempID.index < nextFree {
					e.ids[lastFree].index = tempID.index
					e.ids[tempID.index].index = nextFree
					lastFree = tempID.index
					break
				} else {
					lastFree = nextFree
					nextFree = e.ids[nextFree].index
				}
			}
		}
	}
	e.delayFree = 0
}
