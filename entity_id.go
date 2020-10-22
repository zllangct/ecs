package ecs

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var seq uint32
var timestamp uint64

type EntityID = uint64

func init() {
	rand.Seed(time.Now().UnixNano())
}

func UniqueID() EntityID {
	tNow := uint64(time.Now().UnixNano()) << 32
	tTemp := atomic.LoadUint64(&timestamp)
	if tTemp != tNow {
		atomic.StoreUint32(&seq, 0)
		for {
			if atomic.CompareAndSwapUint64(&timestamp, tTemp, tNow) {
				break
			} else {
				tTemp = atomic.LoadUint64(&timestamp)
				tNow = uint64(time.Now().UnixNano()) << 32
			}
		}
	}
	s := atomic.AddUint32(&seq, 1)
	return tNow + uint64(rand.Uint32()&0xFFFF0000+s)
}
