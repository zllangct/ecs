package ecs

import (
	"cmp"
	"fmt"
	"math/bits"
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
	// Optional tree of free slots for scattered releases. Its second half
	// contains one maximum per page; every parent holds its children's max.
	// Zero means no known free slot. The linked list remains authoritative.
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
	id := ReuseID{index: e.free}
	if e.free == e.pending {
		e.ids = append(e.ids, id)
		id = e.ids[e.pending]
		e.free++
		e.pending++
	} else {
		slot := &e.ids[id.index]
		id.reuse = slot.reuse
		e.free = slot.index
		slot.index = id.index
	}
	e.len++
	if e.freePages != nil {
		e.forgetPageAnchor(id.index)
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
	// Sequential frees only need mergeHint. A singleton window may also
	// need a long search from the head when releases arrive out of order.
	if e.freePages != nil || (e.mergeHint > first && last-e.free >= freePageSize) {
		e.growFreePages(last)
	}

	link := &e.free
	var previous EntityIndex
	if e.mergeHint != 0 && e.mergeHint < first {
		link = &e.ids[e.mergeHint].index
		previous = e.mergeHint
	}
	for _, id := range delayed {
		if e.freePages != nil {
			if anchor := e.previousPageAnchor(int(id.index) / freePageSize); anchor > previous {
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
		if e.freePages != nil {
			e.recordPageAnchor(id.index)
		}
		previous = id.index
	}
	e.delayFree = 0
	e.mergeHint = last
}

func (e *EntityIDGenerator) growFreePages(last EntityIndex) {
	oldLeaves := len(e.freePages) / 2
	page := int(last) / freePageSize
	if page < oldLeaves {
		return
	}
	leaves := 1 << bits.Len(uint(page))
	tree := make([]EntityIndex, 2*leaves)
	copy(tree[leaves:], e.freePages[oldLeaves:])
	// Index rows that were already free before this range was tracked.
	// The never-allocated tail of the initial freelist has generation zero.
	start := e.free
	if oldLeaves > 0 && e.freePages[1] != 0 {
		start = e.ids[e.freePages[1]].index
	}
	for id := start; id != e.pending && int(id)/freePageSize < leaves; id = e.ids[id].index {
		if e.ids[id].reuse == 0 {
			break
		}
		tree[leaves+int(id)/freePageSize] = id
	}
	for i := leaves - 1; i > 0; i-- {
		tree[i] = max(tree[2*i], tree[2*i+1])
	}
	e.freePages = tree
}

// previousPageAnchor returns the nearest known predecessor in a lower page.
// Dense releases usually hit the adjacent page. A sparse release skips empty
// ranges by ascending the tree until a nonempty left sibling is found.
func (e *EntityIDGenerator) previousPageAnchor(page int) EntityIndex {
	if page == 0 {
		return 0
	}
	node := len(e.freePages)/2 + page
	if adjacent := e.freePages[node-1]; adjacent != 0 {
		return adjacent
	}
	for ; node > 1; node /= 2 {
		if node%2 != 0 && e.freePages[node-1] != 0 {
			return e.freePages[node-1]
		}
	}
	return 0
}

func (e *EntityIDGenerator) recordPageAnchor(index EntityIndex) {
	leaf := len(e.freePages)/2 + int(index)/freePageSize
	for node := leaf; node > 0 && e.freePages[node] < index; node /= 2 {
		e.freePages[node] = index
	}
}

func (e *EntityIDGenerator) forgetPageAnchor(index EntityIndex) {
	node := len(e.freePages)/2 + int(index)/freePageSize
	if node >= len(e.freePages) || e.freePages[node] != index {
		return
	}
	// Allocation takes the smallest free ID. When the page maximum is
	// taken, all smaller known free IDs in this page have already gone.
	e.freePages[node] = 0
	for node /= 2; node > 0 && e.freePages[node] == index; node /= 2 {
		e.freePages[node] = max(e.freePages[2*node], e.freePages[2*node+1])
	}
}
