package ecs

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func TestEntityIDGenerator_NewID(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		e := NewEntityIDGenerator(10, 3)
		id1 := e.NewID()
		id2 := e.NewID()
		id3 := e.NewID()

		e.FreeID(id2)

		id4 := e.NewID()

		e.FreeID(id1)
		e.FreeID(id4)
		e.FreeID(id3)

		var m []int64
		for i := 0; i < 11; i++ {
			newID := e.NewID()
			m = append(m, newID)
		}

		for _, id := range m {
			e.FreeID(id)
		}
	})
}

type _EntityTest struct {
	seq      int32
	freeList map[RealID]struct{}
}

func (e *_EntityTest) NewID() int64 {
	id := RealID{}
	if len(e.freeList) > 0 {
		for i, _ := range e.freeList {
			id = i
			delete(e.freeList, i)
		}
	} else {
		id = RealID{index: e.seq, reuse: 0}
		e.seq++
	}
	return id.ToInt64()
}

func (e *_EntityTest) FreeID(id int64) {
	real := *(*RealID)(unsafe.Pointer(&id))
	e.freeList[real] = struct{}{}
}

func BenchmarkEntityIDGenerator_New(b *testing.B) {
	e := NewEntityIDGenerator(1024, 10)
	idmap := map[int64]struct{}{}
	e2 := &_EntityTest{freeList: map[RealID]struct{}{}}
	idmap2 := map[int64]struct{}{}

	for i := 0; i < 10000; i++ {
		id := e.NewID()
		idmap[id] = struct{}{}

		id = e2.NewID()
		idmap2[id] = struct{}{}
	}

	b.Run("map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e.NewID()
		}
	})
	b.Run("gen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			e2.NewID()
		}
	})
}

func TestEntityIDGenerator_Free(t *testing.T) {
	size := 100000
	e := NewEntityIDGenerator(1024, 100)
	idmap := map[int64]struct{}{}
	for i := 0; i < size; i++ {
		id := e.NewID()
		idmap[id] = struct{}{}
	}

	e2 := &_EntityTest{freeList: map[RealID]struct{}{}}
	idmap2 := map[int64]struct{}{}
	for i := 0; i < size; i++ {
		id := e2.NewID()
		idmap2[id] = struct{}{}
	}

	t.Run("gen", func(t *testing.T) {
		start := time.Now()
		for i, _ := range idmap {
			e.FreeID(i)
		}
		println(time.Since(start).Nanoseconds() / int64(size))
	})
	t.Run("map", func(t *testing.T) {
		start := time.Now()
		for i, _ := range idmap2 {
			e2.FreeID(i)
		}
		println(time.Since(start).Nanoseconds() / int64(size))
	})
}

func BenchmarkEntityIDGenerator_NewFreeRandom(b *testing.B) {

	e := NewEntityIDGenerator(1024, 100)
	idmap := map[int64]struct{}{}
	e2 := &_EntityTest{freeList: map[RealID]struct{}{}}
	idmap2 := map[int64]struct{}{}

	for i := 0; i < 1000000; i++ {
		id := e.NewID()
		idmap[id] = struct{}{}

		id = e2.NewID()
		idmap2[id] = struct{}{}
	}

	b.Run("map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := rand.Intn(100)
			if r%2 == 0 {
				id := e.NewID()
				idmap[id] = struct{}{}
			} else {
				for i2, _ := range idmap {
					e.FreeID(i2)
					delete(idmap, i2)
					break
				}
			}
		}
	})
	b.Run("gen", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := rand.Intn(100)
			if r%2 == 0 {
				id := e2.NewID()
				idmap[id] = struct{}{}
			} else {
				for i2, _ := range idmap {
					e2.FreeID(i2)
					delete(idmap2, i2)
					break
				}
			}
		}
	})
}
