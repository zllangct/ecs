package ecs

import (
	"iter"
	"unsafe"
)

var ECS string = "ecs"

// NewQuery 创建查询。
//
// Deprecated: 使用 (*SystemContext).NewQuery。
func NewQuery(ctx *SystemContext, opt ...QueryOption) Query {
	return ctx.NewQuery(opt...)
}

// NewQuery 创建查询。
// 查询的预解析结果（主迭代源、buddy 池、访问槽位）按 shapeKey 缓存于
// SystemContext，跨帧复用；组件池版本变化或组表未定型时自动重建。
func (ctx *SystemContext) NewQuery(opt ...QueryOption) Query {
	qi := Query{ctx: ctx}
	if !ctx.constraint.isValid() || len(opt) == 0 {
		return qi
	}

	c := &QueryConfig{}
	c.initDefault()
	for _, option := range opt {
		option(c)
	}
	if len(c.queryBuddies) == 0 {
		return qi
	}
	c.shapeKey = NewFixedCompound(c.queryBuddies)

	// 组表未定型（注册/初始化期）不做缓存，每次新建
	if !ctx.world.archetypes.sealed {
		qi.cached = ctx.resolveQuery(c)
		return qi
	}
	if ctx.queryCache == nil {
		ctx.queryCache = map[FixedCompound]*queryCached{}
	}
	cached, ok := ctx.queryCache[c.shapeKey]
	if !ok || cached.version != ctx.world.compVersion {
		cached = ctx.resolveQuery(c)
		ctx.queryCache[c.shapeKey] = cached
	}
	qi.cached = cached
	return qi
}

// resolveQuery 一次性完成查询的全部循环不变解析：
// 主迭代源（组表行 或 最小 buddy 池）、buddy 池列表、访问槽位。
func (ctx *SystemContext) resolveQuery(c *QueryConfig) *queryCached {
	cached := &queryCached{config: c, version: ctx.world.compVersion}

	// 组表分派：buddies 全部属于同一张组表时，直接以组表行为主迭代源
	var group *Archetype
	for i, buddy := range c.queryBuddies {
		a, ok := ctx.world.archetypes.owner(buddy)
		if !ok {
			group = nil
			break
		}
		if i == 0 {
			group = a
		} else if a != group {
			group = nil
			break
		}
	}
	cached.group = group

	// buddy 池 + 最小池
	cached.pools = make([]ComponentSet, len(c.queryBuddies))
	var minSet ComponentSet
	for i, buddy := range c.queryBuddies {
		s, ok := ctx.world.getComponentSet(buddy)
		if !ok {
			// 组路径下 buddy 数据可全在表内（无残段池）；CSet 路径缺池即恒空
			if group == nil {
				cached.empty = true
				return cached
			}
			continue
		}
		cached.pools[i] = s
		if minSet == nil || s.Len() < minSet.Len() {
			minSet = s
		}
	}
	cached.minSet = minSet
	if group == nil && minSet == nil {
		cached.empty = true
		return cached
	}

	// 访问槽位：仅声明为非只读依赖的 buddy（权限语义与 GetBuddy 一致）
	cached.slots = map[ComponentIntType]querySlot{}
	for _, buddy := range c.queryBuddies {
		dep, ok := ctx.info.getDep(buddy)
		if !ok || dep.readonly() {
			continue
		}
		if a, ok := ctx.world.archetypes.owner(buddy); ok {
			col, _ := a.colOf(buddy)
			cached.slots[buddy] = func(index EntityIndex) unsafe.Pointer {
				row, ok := a.rowOf(index)
				if !ok {
					return nil
				}
				return a.cellPtr(col, row)
			}
			continue
		}
		s := cached.pools[buddyIndex(c, buddy)]
		if s == nil {
			continue
		}
		cached.slots[buddy] = func(index EntityIndex) unsafe.Pointer {
			return s.get(index)
		}
	}

	return cached
}

// buddyIndex buddy 在 queryBuddies 中的下标
func buddyIndex(c *QueryConfig, buddy ComponentIntType) int {
	for i, b := range c.queryBuddies {
		if b == buddy {
			return i
		}
	}
	return -1
}

// GetComponents 遍历本 system 声明过依赖的组件集合。
//
// Deprecated: 使用 (*SystemContext).GetComponents。
func GetComponents[T any, TP ComponentPointer[T]](ctx *SystemContext) iter.Seq2[EntityIndex, *T] {
	return ctx.GetComponents[T, TP]()
}

// GetComponents 遍历本 system 声明过依赖的组件集合。
// 约定（静默失败）：约束失效、依赖未声明（WithDep）、组件集不存在、类型不匹配时
// 均返回空迭代器。依赖未声明是配置错误，请在开发期通过 system 行为缺失发现。
func (ctx *SystemContext) GetComponents[T any, TP ComponentPointer[T]]() iter.Seq2[EntityIndex, *T] {
	empty := func(yield func(EntityIndex, *T) bool) {
	}
	if !ctx.constraint.isValid() {
		return empty
	}
	it := GetIntType[T, TP]()

	dep, ok := ctx.info.getDep(it)
	if !ok {
		return empty
	}

	// 只读依赖不得通过可写 API 访问（约定：静默返回空），
	// 只读访问请使用 GetComponentsReadOnly（codegen 零拷贝视图）。
	if dep.readonly() {
		return empty
	}

	// 组表分派：T 属于某组时，数据 = 组表 T 列 + CSet 残段，两段顺序迭代
	if a, ok := ctx.world.archetypes.owner(it); ok {
		colSeq := colIterSeq[T](a, it)
		s, ok := ctx.world.getComponentSet(it)
		if !ok {
			return colSeq
		}
		set, ok := s.(*CSet[T, TP])
		if !ok {
			return colSeq
		}
		return func(yield func(EntityIndex, *T) bool) {
			for idx, p := range colSeq {
				if !yield(idx, p) {
					return
				}
			}
			for idx, p := range set.Iter() {
				if !yield(idx, p) {
					return
				}
			}
		}
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return empty
	}
	set, ok := s.(*CSet[T, TP])
	if !ok {
		return empty
	}

	return set.Iter()
}

// colIterSeq 组表列迭代（定长 stride 步进，零查找）
func colIterSeq[T any](a *Archetype, it ComponentIntType) iter.Seq2[EntityIndex, *T] {
	col, _ := a.colOf(it)
	return func(yield func(EntityIndex, *T) bool) {
		size := a.sizes[col]
		data := a.columns[col]
		for row := 0; row < a.Len(); row++ {
			if !yield(a.entities[row], (*T)(unsafe.Pointer(&data[row*size]))) {
				return
			}
		}
	}
}

// GetComponentsReadOnly 以只读视图遍历本 system 声明为只读依赖的组件集合。
//
// Deprecated: 使用 (*SystemContext).GetComponentsReadOnly。
func GetComponentsReadOnly[TR ReadOnlyView[TR]](ctx *SystemContext) iter.Seq2[EntityIndex, TR] {
	return ctx.GetComponentsReadOnly[TR]()
}

// GetComponentsReadOnly 以只读视图遍历本 system 声明为只读依赖的组件集合。
// TR 为 codegen 生成的零拷贝只读视图（仅 getter，编译期防写），
// 视图自描述组件类型信息，无需再书写组件类型参数。
// 约定（静默失败）：约束失效、依赖未声明、依赖为可写（ReadWrite）、
// 组件集不存在时均返回空迭代器——只读视图仅限只读依赖使用。
func (ctx *SystemContext) GetComponentsReadOnly[TR ReadOnlyView[TR]]() iter.Seq2[EntityIndex, TR] {
	empty := func(yield func(EntityIndex, TR) bool) {
	}
	if !ctx.constraint.isValid() {
		return empty
	}
	var zero TR
	it := ComponentIntType(zero.ComponentPacketIdentifier())

	dep, ok := ctx.info.getDep(it)
	if !ok {
		return empty
	}
	if !dep.readonly() {
		return empty
	}

	// 组表分派：先扫组表列，再扫 CSet 残段
	if a, ok := ctx.world.archetypes.owner(it); ok {
		col, _ := a.colOf(it)
		return func(yield func(EntityIndex, TR) bool) {
			size := a.sizes[col]
			data := a.columns[col]
			for row := 0; row < a.Len(); row++ {
				if !yield(a.entities[row], zero.FromPtr(unsafe.Pointer(&data[row*size]))) {
					return
				}
			}
			if s, ok := ctx.world.getComponentSet(it); ok {
				for index, p := range s.iterPtr() {
					if !yield(index, zero.FromPtr(p)) {
						return
					}
				}
			}
		}
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return empty
	}

	return func(yield func(EntityIndex, TR) bool) {
		for index, p := range s.iterPtr() {
			if !yield(index, zero.FromPtr(p)) {
				return
			}
		}
	}
}

// GetBuddy 按实体索引读取组件。
//
// Deprecated: 使用 (*SystemContext).GetBuddy。
func GetBuddy[T any, TP ComponentPointer[T]](ctx *SystemContext, index EntityIndex) (*T, bool) {
	return ctx.GetBuddy[T, TP](index)
}

// GetBuddy 按实体索引读取组件。
func (ctx *SystemContext) GetBuddy[T any, TP ComponentPointer[T]](index EntityIndex) (*T, bool) {
	if !ctx.constraint.isValid() {
		return nil, false
	}
	it := GetIntType[T, TP]()
	// 依赖未声明不得访问（约定：静默返回空）
	dep, ok := ctx.info.getDep(it)
	if !ok {
		return nil, false
	}
	// 只读依赖不得通过可写 API 访问（约定：静默返回空），
	// 只读访问请使用 GetBuddyReadOnly（codegen 零拷贝视图）。
	if dep.readonly() {
		return nil, false
	}

	// 组表分派：命中行则直取列内存
	if a, ok := ctx.world.archetypes.owner(it); ok {
		if row, ok := a.rowOf(index); ok {
			col, _ := a.colOf(it)
			return (*T)(a.cellPtr(col, row)), true
		}
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return nil, false
	}
	b := s.get(index)
	if b == nil {
		return nil, false
	}
	return (*T)(b), true
}

// GetBuddyReadOnly 按实体索引以只读视图读取组件。
//
// Deprecated: 使用 (*SystemContext).GetBuddyReadOnly。
func GetBuddyReadOnly[TR ReadOnlyView[TR]](ctx *SystemContext, index EntityIndex) (TR, bool) {
	return ctx.GetBuddyReadOnly[TR](index)
}

// GetBuddyReadOnly 按实体索引以只读视图读取组件。
// 约定（静默失败）：约束失效、依赖未声明、依赖为可写（ReadWrite）、
// 组件集不存在、索引无组件时均返回 (零值, false)——只读视图仅限只读依赖使用。
func (ctx *SystemContext) GetBuddyReadOnly[TR ReadOnlyView[TR]](index EntityIndex) (TR, bool) {
	var zero TR
	if !ctx.constraint.isValid() {
		return zero, false
	}
	it := ComponentIntType(zero.ComponentPacketIdentifier())
	dep, ok := ctx.info.getDep(it)
	if !ok {
		return zero, false
	}
	if !dep.readonly() {
		return zero, false
	}

	// 组表分派：命中行则直取列内存
	if a, ok := ctx.world.archetypes.owner(it); ok {
		if row, ok := a.rowOf(index); ok {
			col, _ := a.colOf(it)
			return zero.FromPtr(a.cellPtr(col, row)), true
		}
	}

	s, ok := ctx.world.getComponentSet(it)
	if !ok {
		return zero, false
	}
	b := s.get(index)
	if b == nil {
		return zero, false
	}
	return zero.FromPtr(b), true
}
