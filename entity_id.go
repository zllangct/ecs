package ecs

import (
	"math/rand"
	"sync/atomic"
	"time"
)

var seq uint32
var timestamp int64

type EntityID = int64

func init() {
	rand.Seed(time.Now().UnixNano())
}

func UniqueID() EntityID {
	tNow := int64(time.Now().UnixNano()) << 32
	tTemp := atomic.LoadInt64(&timestamp)
	if tTemp != tNow {
		atomic.StoreUint32(&seq, 0)
		for {
			if atomic.CompareAndSwapInt64(&timestamp, tTemp, tNow) {
				break
			} else {
				tTemp = atomic.LoadInt64(&timestamp)
				tNow = int64(time.Now().UnixNano()) << 32
			}
		}
	}
	s := atomic.AddUint32(&seq, 1)
	return tNow + int64((s<<16)&0xFFFF0000+rand.Uint32()&0x0000FFFF)
}
