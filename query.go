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
	// group buddies 全部属于同一张组表时，直接以组表行为主迭代源（零过滤连续扫描）
	group *Archetype
	// minSet CSet 路径的最小组件池（主迭代源）
	minSet ComponentSet
	// pools 与 config.queryBuddies 平行的 buddy 池（CSet 路径 Exist 过滤用）
	pools []ComponentSet
	// slots buddy → 预解析访问槽位（仅声明为非只读依赖的 buddy）
	slots map[ComponentIntType]querySlot
	// empty 解析时即知结果恒空（某 buddy 池不存在）
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
		// 组表路径：行 ⟺ 拥有组全集 ⊇ buddies，过滤恒真，跳过 IsSubSet。
		g := c.group
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
	if c.minSet == nil || c.minSet.Len() == 0 {
		return empty
	}
	// merge-join：全部 buddy 池按 key 有序时，有序归并求交（全程顺序访问，
	// 取代逐 key 稀疏查找与过滤）；任一池失序回退 Exist 过滤路径
	if mergeable(c.pools) {
		return mergeJoinIter(q.ctx.world, c.minSet, c.pools)
	}
	idx := c.minSet.EntityIndexes()
	pools := c.pools
	minSet := c.minSet
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
