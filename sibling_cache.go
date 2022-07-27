package ecs

import (
	"sync"
	"sync/atomic"
)

var cacheOpPool = &sync.Pool{
	New: func() interface{} {
		return &CacheOp{}
	},
}

type siblingInfo struct {
	compound Compound
}

type CacheOp struct {
	Entity Entity
	It     uint16
	Op     uint8
	O      uint8
}

func (s *CacheOp) Set(entity Entity, it uint16, op uint8) *CacheOp {
	s.Entity = entity
	s.It = it
	s.Op = op
	return s
}

type siblingCache struct {
	world        *ecsWorld
	cache        map[Entity]*siblingInfo
	ch           chan *CacheOp
	pause        chan chan struct{}
	isCollecting atomic.Bool
}

func newSiblingCache(world *ecsWorld, chanSize int) *siblingCache {
	return &siblingCache{
		cache: make(map[Entity]*siblingInfo),
		ch:    make(chan *CacheOp, chanSize),
		pause: make(chan chan struct{}),
	}
}

func (s *siblingCache) StartCollector() {
	if !s.isCollecting.CompareAndSwap(false, true) {
		return
	}
	go func() {
		var p chan struct{}
		var op *CacheOp
		for {
			select {
			case p = <-s.pause:
				for !s.isCollecting.CompareAndSwap(true, false) {
				}
				p <- struct{}{}
				return
			case op = <-s.ch:
				switch op.Op {
				case 1:
					if cache, ok := s.cache[op.Entity]; ok {
						_ = cache.compound.Add(op.It)
					} else {
						s.cache[op.Entity] = &siblingInfo{
							compound: Compound{op.It},
						}
					}
				case 2:
					if cache, ok := s.cache[op.Entity]; ok {
						cache.compound.Remove(op.It)
					}
				}
			}
		}
	}()
}

func (s *siblingCache) PauseCollector() {
	r := make(chan struct{})
	s.pause <- r
	<-r
	return
}
