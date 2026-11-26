package ecs

import (
	"cmp"
	"fmt"
	"slices"
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
	// Last merged free slot. Allocation and restore invalidate the hint,
	// so it always points into this generator's current sorted freelist.
	mergeHint EntityIndex
	// For scattered frees, keep the largest known free slot per page.
	// Zero means that the page has no known anchor; the list is authoritative.
	freePages []EntityIndex
}

const freePageSize = 128

func NewEntityIDGenerator(initSize int, delayCap int) *EntityIDGenerator {
	if initSize < 2 {
		initSize = 2
	}
	if delayCap <= 0 {
		delayCap = 10
	}
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
	e.mergeHint = 0
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
	if e.freePages != nil {
		page := int(id.index) / freePageSize
		if page < len(e.freePages) && e.freePages[page] == id.index {
			// Allocation takes the smallest free ID. Once the largest known
			// slot is taken, every smaller known slot in that page is gone.
			e.freePages[page] = 0
		}
	}
	return id.ToEntity()
}

// FreeID 释放实体 ID。index 不会立即进入 freelist，而是攒入 removeDelay 缓冲，
// 凑满 delayCap 后批量归并（delayFlush），避免高频释放时的链表维护开销。
// 重复释放、越界、代数不匹配的句柄属于逻辑错误，直接 panic。
func (e *EntityIDGenerator) FreeID(entity Entity) {
	realID := entity.toReuseID()
	if realID.index <= 0 || realID.index >= EntityIndex(len(e.ids)) {
		panic(fmt.Sprintf("ecs: FreeID out of range, index=%d", realID.index))
	}
	slot := &e.ids[realID.index]
	if slot.index == -1 {
		panic(fmt.Sprintf("ecs: double free, index=%d", realID.index))
	}
	if slot.index != realID.index {
		panic(fmt.Sprintf("ecs: free unallocated entity slot, index=%d", realID.index))
	}
	if slot.reuse != realID.reuse {
		panic(fmt.Sprintf("ecs: stale entity handle, index=%d want reuse=%d got %d",
			realID.index, slot.reuse, realID.reuse))
	}

	e.len--
	slot.index = -1 // 标记为待回收
	slot.reuse++

	e.removeDelay[e.delayFree] = realID
	e.delayFree++
	if e.delayFree >= e.delayCap {
		e.delayFlush()
	}
}

// delayFlush 将排序后的 removeDelay 归并进 ids 内嵌的 freelist 链。
// freelist 链不变量：free 为链头；空闲槽 x 的 ids[x].index 为下一空闲 index；
// 链尾指向 pending（哨兵，表示链结束）；链严格升序，保证小号 index 优先复用。
// 只遍历 freelist 链与缓冲批，代价与链长/批量成正比，不做 ids 全表扫描。
func (e *EntityIDGenerator) delayFlush() {
	delayed := e.removeDelay[:e.delayFree]
	slices.SortFunc(delayed, func(a, b ReuseID) int {
		return cmp.Compare(a.index, b.index)
	})

	first, last := delayed[0].index, delayed[len(delayed)-1].index
	pages := int(last)/freePageSize + 1
	// Sequential frees need only mergeHint. Build the optional index when
	// a later window overlaps an earlier one and spans multiple pages.
	if e.freePages != nil || (e.mergeHint > first && last-first >= freePageSize) {
		if pages > len(e.freePages) {
			e.freePages = append(e.freePages, make([]EntityIndex, pages-len(e.freePages))...)
		}
	}

	link := &e.free
	var previous EntityIndex
	if e.mergeHint != 0 && e.mergeHint < first {
		link = &e.ids[e.mergeHint].index
		previous = e.mergeHint
	}
	for _, id := range delayed {
		page := int(id.index) / freePageSize
		// Use at most one lookup; a sparse world must not scan empty pages.
		if e.freePages != nil && page > 0 {
			if anchor := e.freePages[page-1]; anchor > previous {
				link = &e.ids[anchor].index
			}
		}
		for *link != e.pending && *link < id.index {
			link = &e.ids[*link].index
		}
		next := *link
		*link = id.index
		link = &e.ids[id.index].index
		*link = next
		if e.freePages != nil && id.index > e.freePages[page] {
			e.freePages[page] = id.index
		}
		previous = id.index
	}
	e.delayFree = 0
	e.mergeHint = last
}
