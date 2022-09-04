package ecs

import (
	"sort"
	"unsafe"
)

type Entity int64

func (e Entity) ToInt64() int64 {
	return int64(e)
}

func (e Entity) ToRealID() RealID {
	return *(*RealID)(unsafe.Pointer(&e))
}

type RealID struct {
	index int32
	reuse int32
}

func (r *RealID) ToInt64() int64 {
	return *(*int64)(unsafe.Pointer(r))
}

func (r *RealID) ToEntity() Entity {
	return *(*Entity)(unsafe.Pointer(r))
}

type EntityIDGenerator struct {
	ids     []RealID
	free    int32
	pending int32
	len     int32

	removeDelay []RealID
	delayFree   int32
	delayCap    int32
}

func NewEntityIDGenerator(initSize int, delayCap int) *EntityIDGenerator {
	g := &EntityIDGenerator{}
	g.ids = make([]RealID, initSize)
	for i := 0; i < len(g.ids); i++ {
		g.ids[i].index = int32(i + 1)
	}
	g.free = 1
	g.pending = int32(initSize)
	g.len = 0
	g.removeDelay = make([]RealID, delayCap)
	g.delayCap = int32(delayCap)
	g.delayFree = 0
	return g
}

func (e *EntityIDGenerator) NewID() Entity {
	id := RealID{}
	if e.free == e.pending {
		e.ids = append(e.ids, RealID{index: e.free, reuse: 0})
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

	real := entity.ToRealID()
	e.ids[real.index].index = -1

	e.removeDelay[e.delayFree] = real
	e.delayFree++
	if e.delayFree >= e.delayCap {
		e.delayFlush()
	}
	if e.pending > 1024 && e.pending < int32(len(e.ids))/2 {
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
