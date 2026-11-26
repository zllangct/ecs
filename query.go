package ecs

import (
	"iter"
	"unsafe"
)

type QueryBuddies = Compound

func NewQueryBuddies(intType ...ComponentIntType) *QueryBuddies {
	b := &QueryBuddies{}
	for _, componentIntType := range intType {
		b.Add(componentIntType)
	}
	return b
}

type QueryConfig struct {
	queryBuddies QueryBuddies
	shapeKey     FixedCompound // buddies 的定长 key，NewQuery 时计算一次
}

func (q *QueryConfig) initDefault() {

}

type QueryOption func(q *QueryConfig)

func WithBuddies(buddies QueryBuddies) QueryOption {
	return func(q *QueryConfig) {
		q.queryBuddies.Merge(buddies)
	}
}

func WithComp[T any, TP ComponentPointer[T]]() QueryOption {
	it := GetIntType[T, TP]()
	return func(q *QueryConfig) {
		q.queryBuddies.Add(it)
	}
}

// querySlot 预解析的组件访问槽位：循环不变的查找（依赖校验、map、列号）
// 在 NewQuery 时完成，迭代中仅剩一次稀疏定位。
type querySlot func(EntityIndex) unsafe.Pointer

// queryCached 查询的预解析结果，缓存于 SystemContext（按 shapeKey），
// world.compVersion 变化（新组件池创建）时失效重建。
type queryCached struct {
	config *QueryConfig
	// group buddies 全部属于同一张组表时指向该表
	group *Archetype
	// groupFull buddies 精确覆盖 group 全集（无残段实体可能）
	groupFull bool
	// buddyGroups 与 config.queryBuddies 平行的所属组表
	buddyGroups []*Archetype
	// hasGrouped 任一 buddy 属于某张组表
	hasGrouped bool
	// minSet CSet 路径的最小组件池（主迭代源）
	minSet ComponentSet
	// pools 与 config.queryBuddies 平行的 buddy 池（CSet 路径/残段 Exist 过滤用）
	pools []ComponentSet
	// slots buddy → 预解析访问槽位（仅声明为非只读依赖的 buddy）
	slots map[ComponentIntType]querySlot
	// empty 解析时即知结果恒空（某非组 buddy 池不存在）
	empty bool
	// version 构建时的 world.compVersion
	version int
}

type Query struct {
	ctx    *SystemContext
	cached *queryCached
}

func (q *Query) Iter() iter.Seq2[EntityIndex, *EntityInfo] {
	empty := func(yield func(EntityIndex, *EntityInfo) bool) {}
	if !q.ctx.constraint.isValid() || q.cached == nil || q.cached.empty {
		return empty
	}
	c := q.cached
	stat := q.ctx.info.getOptReporter()
	stat.shapeUsageAddKey(c.config.shapeKey)
	if c.group != nil {
		g := c.group
		if c.groupFull {
			// 全集组表路径：行 ⟺ 拥有组全集 == buddies，无残段可能，跳过 IsSubSet 与 CSet
			if g.Len() == 0 {
				return empty
			}
			return func(yield func(EntityIndex, *EntityInfo) bool) {
				for row := 0; row < g.Len(); row++ {
					index := g.entities[row]
					info := q.ctx.world.entities.getByIndex(index)
					if info == nil {
						continue
					}
					if !yield(index, info) {
						return
					}
				}
			}
		}

		// 组表真子集路径：先扫组表行，再扫全部 buddy 在 CSet 残段中的交集
		resMinSet := pickMinNonEmptyPool(c.pools)
		if g.Len() == 0 && resMinSet == nil {
			return empty
		}
		return func(yield func(EntityIndex, *EntityInfo) bool) {
			for row := 0; row < g.Len(); row++ {
				index := g.entities[row]
				info := q.ctx.world.entities.getByIndex(index)
				if info == nil {
					continue
				}
				if !yield(index, info) {
					return
				}
			}
			if resMinSet == nil {
				return
			}
			for index, info := range intersectPools(q.ctx.world, resMinSet, c.pools) {
				if !yield(index, info) {
					return
				}
			}
		}
	}

	// 混合路径（含组内 + 组外组件，或跨多张组表）
	if c.hasGrouped {
		driver := pickMinMixedDriver(c.buddyGroups, c.pools)
		if driver < 0 {
			return empty
		}
		dg := c.buddyGroups[driver]
		dp := c.pools[driver]
		return func(yield func(EntityIndex, *EntityInfo) bool) {
			emit := func(index EntityIndex) bool {
				for j := range c.pools {
					if j == driver {
						continue
					}
					has := false
					if aj := c.buddyGroups[j]; aj != nil {
						_, has = aj.rowOf(index)
					}
					if !has && c.pools[j] != nil {
						has = c.pools[j].Exist(index)
					}
					if !has {
						return true
					}
				}
				info := q.ctx.world.entities.getByIndex(index)
				if info == nil {
					return true
				}
				return yield(index, info)
			}
			if dg != nil {
				for row := 0; row < dg.Len(); row++ {
					if !emit(dg.entities[row]) {
						return
					}
				}
			}
			if dp != nil {
				for _, index := range dp.EntityIndexes() {
					if !emit(index) {
						return
					}
				}
			}
		}
	}

	// 纯 CSet 路径：动态选当前最小池作驱动源（避免历史缓存 minSet 僵化）
	minSet := pickMinNonEmptyPool(c.pools)
	if minSet == nil {
		return empty
	}
	return intersectPools(q.ctx.world, minSet, c.pools)
}

func intersectPools(world *World, minSet ComponentSet, pools []ComponentSet) iter.Seq2[EntityIndex, *EntityInfo] {
	if len(pools) == 2 {
		other := pools[0]
		if other == minSet {
			other = pools[1]
		}
		return intersectPair(world, minSet, other)
	}
	if mergeable(pools) {
		return mergeJoinIter(world, minSet, pools)
	}
	return existFilterIter(world, minSet, pools)
}

// Both pools are nonempty. Keep the smaller pool's traversal order on every path.
func intersectPair(world *World, driver, other ComponentSet) iter.Seq2[EntityIndex, *EntityInfo] {
	keys := driver.EntityIndexes()
	if driver.IsKeyOrdered() && other.IsKeyOrdered() {
		peers := other.EntityIndexes()
		first, last := peers[0], peers[len(peers)-1]
		if keys[len(keys)-1] < first || keys[0] > last {
			return func(func(EntityIndex, *EntityInfo) bool) {}
		}
		// Distinct sorted peer keys filling their entire range prove membership
		// for every driver key inside it, with no per-entity peer access.
		if int64(last)-int64(first)+1 == int64(len(peers)) && keys[0] >= first && keys[len(keys)-1] <= last {
			return func(yield func(EntityIndex, *EntityInfo) bool) {
				for _, key := range keys {
					if info := world.entities.getByIndex(key); info != nil && !yield(key, info) {
						return
					}
				}
			}
		}
		// A merge stops at the final driver key. A small prefix can be cheap even
		// when the rest of the peer pool is large. Use 64-bit math for key spans.
		scanBound := min(int64(len(peers)), int64(keys[len(keys)-1])-int64(first)+1)
		const probeRatio = 8
		if scanBound < probeRatio*int64(len(keys)) {
			return func(yield func(EntityIndex, *EntityInfo) bool) {
				j := 0
				for _, key := range keys {
					for j < len(peers) && peers[j] < key {
						j++
					}
					if j == len(peers) {
						return
					}
					if peers[j] == key {
						j++
						if info := world.entities.getByIndex(key); info != nil && !yield(key, info) {
							return
						}
					}
				}
			}
		}
	}
	return func(yield func(EntityIndex, *EntityInfo) bool) {
		for _, key := range keys {
			if other.Exist(key) {
				if info := world.entities.getByIndex(key); info != nil && !yield(key, info) {
					return
				}
			}
		}
	}
}

// pickMinNonEmptyPool 存在任一 nil 或空池则返回 nil；否则返回元素最少的池
func pickMinNonEmptyPool(pools []ComponentSet) ComponentSet {
	var minSet ComponentSet
	for _, p := range pools {
		if p == nil || p.Len() == 0 {
			return nil
		}
		if minSet == nil || p.Len() < minSet.Len() {
			minSet = p
		}
	}
	return minSet
}

// pickMinMixedDriver 在混合查询中挑选总条目数（组表行 + CSet 残段）最小的 buddy 下标；任一 buddy 为空返回 -1
func pickMinMixedDriver(groups []*Archetype, pools []ComponentSet) int {
	best := -1
	bestLen := 0
	for i := range pools {
		total := 0
		if groups[i] != nil {
			total += groups[i].Len()
		}
		if pools[i] != nil {
			total += pools[i].Len()
		}
		if total == 0 {
			return -1
		}
		if best == -1 || total < bestLen {
			best = i
			bestLen = total
		}
	}
	return best
}

func existFilterIter(world *World, minSet ComponentSet, pools []ComponentSet) iter.Seq2[EntityIndex, *EntityInfo] {
	idx := minSet.EntityIndexes()
	return func(yield func(EntityIndex, *EntityInfo) bool) {
	outer:
		for _, index := range idx {
			// buddy 池稀疏表 Exist 过滤：一次数组索引，
			// 取代 EntityInfo→compound 指针链的 IsSubSet
			for _, pool := range pools {
				if pool == minSet {
					continue
				}
				if !pool.Exist(index) {
					continue outer
				}
			}
			info := world.entities.getByIndex(index)
			if info == nil {
				continue
			}
			if !yield(index, info) {
				return
			}
		}
	}
}

// mergeable 全部 buddy 池可按 key 有序归并（至少 2 个池且均有序）
func mergeable(pools []ComponentSet) bool {
	if len(pools) < 2 {
		return false
	}
	for _, p := range pools {
		if p == nil || !p.IsKeyOrdered() {
			return false
		}
	}
	return true
}

// mergeJoinIter 有序归并求交：以最小池 key 序列为基准，其余池游标同步推进。
// 归并命中即交集，全程顺序内存访问；游标在闭包内初始化，迭代器可重复消费。
func mergeJoinIter(world *World, minSet ComponentSet, pools []ComponentSet) iter.Seq2[EntityIndex, *EntityInfo] {
	keys := minSet.EntityIndexes()
	others := make([][]EntityIndex, 0, len(pools)-1)
	for _, p := range pools {
		if p != minSet {
			others = append(others, p.EntityIndexes())
		}
	}
	return func(yield func(EntityIndex, *EntityInfo) bool) {
		cursors := make([]int, len(others))
	outer:
		for _, k := range keys {
			for j, ok := range others {
				cur := cursors[j]
				for cur < len(ok) && ok[cur] < k {
					cur++
				}
				cursors[j] = cur
				if cur >= len(ok) {
					return // 某池耗尽，交集结束
				}
				if ok[cur] != k {
					continue outer
				}
			}
			info := world.entities.getByIndex(k)
			if info == nil {
				continue
			}
			if !yield(k, info) {
				return
			}
		}
	}
}

// QueryGet 经查询预解析槽位直取 buddy 组件。
// 约定（静默失败）：查询无效、类型不在 buddies、依赖未声明或为只读、
// 索引无组件时均返回 (nil, false)。
func QueryGet[T any, TP ComponentPointer[T]](q *Query, index EntityIndex) (*T, bool) {
	if q == nil || q.ctx == nil || !q.ctx.constraint.isValid() || q.cached == nil {
		return nil, false
	}
	slot, ok := q.cached.slots[GetIntType[T, TP]()]
	if !ok || slot == nil {
		return nil, false
	}
	p := slot(index)
	if p == nil {
		return nil, false
	}
	return (*T)(p), true
}
