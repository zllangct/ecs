package ecs

import (
	"iter"
	"unsafe"
)

// View2 组表双列 lockstep 迭代：行内零查找，双列指针同步递增。
// 要求 A、B 均为本 system 声明的非只读依赖；两者同属一张组表时走 lockstep，
// 否则回退 Query{A,B} 通用路径（结果集恒正确，仅性能差异）。
// 约定（静默失败）：约束失效、依赖未声明、任一依赖为只读时返回空迭代器。
func View2[A, B any, PA ComponentPointer[A], PB ComponentPointer[B]](ctx *SystemContext) iter.Seq2[*A, *B] {
	empty := func(yield func(*A, *B) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	itB := GetIntType[B, PB]()
	depA, okA := ctx.info.getDep(itA)
	depB, okB := ctx.info.getDep(itB)
	if !okA || !okB || depA.readonly() || depB.readonly() {
		return empty
	}

	// lockstep 路径：同属一张组表
	aA, okA := ctx.world.archetypes.owner(itA)
	aB, okB := ctx.world.archetypes.owner(itB)
	if okA && okB && aA == aB {
		colA, _ := aA.colOf(itA)
		colB, _ := aA.colOf(itB)
		sizeA, sizeB := aA.sizes[colA], aA.sizes[colB]
		dataA, dataB := aA.columns[colA], aA.columns[colB]
		rows := aA.Len()
		isSubset := distinctComponentCount(itA, itB) < len(aA.types)
		return func(yield func(*A, *B) bool) {
			for row := 0; row < rows; row++ {
				if !yield(
					(*A)(unsafe.Pointer(&dataA[row*sizeA])),
					(*B)(unsafe.Pointer(&dataB[row*sizeB])),
				) {
					return
				}
			}
			if isSubset {
				for ptrs := range iterResidualPtrs(ctx, itA, itB) {
					if !yield((*A)(ptrs[0]), (*B)(ptrs[1])) {
						return
					}
				}
			}
		}
	}

	// 回退路径：Query{A,B} + 预解析槽位
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB]())
	return func(yield func(*A, *B) bool) {
		for index := range q.Iter() {
			a, okA := QueryGet[A, PA](&q, index)
			b, okB := QueryGet[B, PB](&q, index)
			if !okA || !okB {
				continue
			}
			if !yield(a, b) {
				return
			}
		}
	}
}

// View2RO 组表双列 lockstep 迭代（A 可写 + B 只读视图）。
// 要求 A 为非只读依赖、B 为只读依赖（视图自描述类型信息）；
// 同属一张组表时走 lockstep，否则回退通用路径。
// 约定（静默失败）：约束失效、依赖未声明/权限不匹配时返回空迭代器。
func View2RO[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq2[*A, BR] {
	return buildView2RO(ctx, prepareView2RO[A, BR, PA])
}

// 构造迭代器时即解析权限与列。将解析函数作为参数，使入口足够小，
// Go 编译器可以把完整组循环与调用方的 range 循环体一起内联。
func buildView2RO[A any, BR ReadOnlyView[BR]](ctx *SystemContext, resolve func(*SystemContext) view2ROPlan[A, BR]) iter.Seq2[*A, BR] {
	plan := resolve(ctx)
	return func(yield func(*A, BR) bool) {
		if plan.fallback != nil {
			plan.fallback(yield)
			return
		}
		var view BR
		offsetA, offsetB := 0, 0
		// 每列长度为 rows * stride。只在访问当前行时构造指针，末行之后仅推进整数偏移。
		for range plan.rows {
			a := (*A)(unsafe.Add(plan.baseA, offsetA))
			b := view.FromPtr(unsafe.Add(plan.baseB, offsetB))
			if !yield(a, b) {
				return
			}
			offsetA += plan.strideA
			offsetB += plan.strideB
		}
	}
}

type view2ROPlan[A, BR any] struct {
	baseA, baseB     unsafe.Pointer
	strideA, strideB int
	rows             int
	fallback         iter.Seq2[*A, BR]
}

func prepareView2RO[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext) (plan view2ROPlan[A, BR]) {
	if !ctx.constraint.isValid() {
		return
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	depA, okA := ctx.info.getDep(itA)
	depB, okB := ctx.info.getDep(itB)
	if !okA || !okB || depA.readonly() || !depB.readonly() {
		return
	}

	aA, okA := ctx.world.archetypes.owner(itA)
	aB, okB := ctx.world.archetypes.owner(itB)
	if okA && okB && aA == aB {
		colA, _ := aA.colOf(itA)
		colB, _ := aA.colOf(itB)
		sizeA, sizeB := aA.sizes[colA], aA.sizes[colB]
		dataA, dataB := aA.columns[colA], aA.columns[colB]
		rows := aA.Len()
		isSubset := distinctComponentCount(itA, itB) < len(aA.types)
		if !isSubset {
			plan.baseA = unsafe.Pointer(unsafe.SliceData(dataA))
			plan.baseB = unsafe.Pointer(unsafe.SliceData(dataB))
			plan.strideA, plan.strideB, plan.rows = sizeA, sizeB, rows
			return
		}
		plan.fallback = func(yield func(*A, BR) bool) {
			for row := 0; row < rows; row++ {
				if !yield(
					(*A)(unsafe.Pointer(&dataA[row*sizeA])),
					zeroB.FromPtr(unsafe.Pointer(&dataB[row*sizeB])),
				) {
					return
				}
			}
			for ptrs := range iterResidualPtrs(ctx, itA, itB) {
				if !yield((*A)(ptrs[0]), zeroB.FromPtr(ptrs[1])) {
					return
				}
			}
		}
		return
	}

	if !okA && !okB {
		// 两列只存于 CSet，依赖权限在入口已校验；每行直接从这两个池取组件。
		setA, hasA := ctx.world.getComponentSet(itA)
		setB, hasB := ctx.world.getComponentSet(itB)
		if !hasA || !hasB {
			return
		}
		nativeA, okA := setA.(*CSet[A, PA])
		nativeB, okB := setB.(interface{ nativeCSet() ComponentSet })
		if okA && okB && nativeB.nativeCSet() == setB {
			if scan := view2ROCSet[A, BR](ctx, itA, itB, nativeA, setB); scan != nil {
				plan.fallback = scan
				return
			}
		}
		plan.fallback = queryView2ROCSet[A, BR, PA](ctx, itB, setA, setB)
		return
	}
	plan.fallback = prepareMixedView2RO[A, BR, PA](ctx, itA, itB, aA, aB)
	return
}

func prepareMixedView2RO[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext, itA, itB ComponentIntType, groupA, groupB *Archetype) iter.Seq2[*A, BR] {
	if groupA != nil && groupB == nil {
		setA, _ := ctx.world.getComponentSet(itA)
		setB, _ := ctx.world.getComponentSet(itB)
		nativeA, aOK := setA.(*CSet[A, PA])
		nativeB, bOK := setB.(interface{ nativeCSet() ComponentSet })
		if aOK && bOK && nativeB.nativeCSet() == setB {
			colA, _ := groupA.colOf(itA)
			return mixedView2RO[A, BR](ctx, itA, itB, groupA, colA, nativeA, setB)
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB)))
	return func(yield func(*A, BR) bool) {
		for index := range q.Iter() {
			a, okA := QueryGet[A, PA](&q, index)
			if !okA {
				continue
			}
			bv, okB := ctx.GetBuddyReadOnly[BR](index)
			if !okB {
				continue
			}
			if !yield(a, bv) {
				return
			}
		}
	}
}

// 原查询负责自定义池以及需要归并的池；消费时才捕获双方键切片与归并游标。
func queryView2ROCSet[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext, itB ComponentIntType, a, b ComponentSet) iter.Seq2[*A, BR] {
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB)))
	return func(yield func(*A, BR) bool) {
		var view BR
		for key := range q.Iter() {
			if !ctx.constraint.isValid() {
				return
			}
			pa, pb := a.get(key), b.get(key)
			if pa != nil && pb != nil && !yield((*A)(pa), view.FromPtr(pb)) {
				return
			}
		}
	}
}

// 每次按当前大小选驱动列，等长时使用类型编号较小的一列，与 Query 顺序一致。
func csetViewKeys(itA, itB ComponentIntType, a, b ComponentSet) ([]EntityIndex, pairStrategy) {
	keys, others := a.EntityIndexes(), b.EntityIndexes()
	if len(others) < len(keys) || (len(others) == len(keys) && itB < itA) {
		keys, others = others, keys
	}
	if len(keys) == 0 {
		return keys, pairDisjoint
	}
	if a.IsKeyOrdered() && b.IsKeyOrdered() {
		return keys, orderedPairStrategy(keys, others)
	}
	return keys, pairProbe
}

func view2ROCSet[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext, itA, itB ComponentIntType, a *CSet[A, PA], b ComponentSet) iter.Seq2[*A, BR] {
	if _, strategy := csetViewKeys(itA, itB, a, b); strategy == pairMerge {
		return nil
	}
	shape := FixedCompound2{min(itA, itB), max(itA, itB)}
	return func(yield func(*A, BR) bool) {
		if !ctx.constraint.isValid() {
			return
		}
		keys, strategy := csetViewKeys(itA, itB, a, b)
		if strategy == pairMerge {
			queryView2ROCSet[A, BR, PA](ctx, itB, a, b)(yield)
			return
		}
		ctx.info.getOptReporter().shapeUsageAddKey(shape)
		if strategy == pairDisjoint {
			return
		}
		var view BR
		for row, key := range keys {
			if !ctx.constraint.isValid() {
				return
			}
			var pa *A
			// 同一行在结构同步点后可能属于其他实体，只对当前行键匹配的覆盖区间直取。
			if strategy == pairCovered && row < len(a.idx2Key) && a.idx2Key[row] == int32(key) {
				pa = &a.data[row]
			} else {
				pa = a.getByIndex(key)
			}
			if pa == nil || !ctx.world.entities.SparseArray.Exist(key) {
				continue
			}
			if pb := b.get(key); pb != nil && !yield(pa, view.FromPtr(pb)) {
				return
			}
		}
	}
}

// Choose the same driving column and tie order as Query.Iter on every consumption.
// Re-read group rows and component addresses when a suspended consumer resumes.
func mixedView2RO[A any, BR ReadOnlyView[BR], PA ComponentPointer[A]](ctx *SystemContext, itA, itB ComponentIntType, group *Archetype, col int, a *CSet[A, PA], b ComponentSet) iter.Seq2[*A, BR] {
	shape := FixedCompound2{min(itA, itB), max(itA, itB)}
	return func(yield func(*A, BR) bool) {
		if !ctx.constraint.isValid() {
			return
		}
		ctx.info.getOptReporter().shapeUsageAddKey(shape)
		totalA, totalB := group.Len()+a.Len(), b.Len()
		if totalA == 0 || totalB == 0 {
			return
		}
		var view BR
		emit := func(key EntityIndex, pa *A) bool {
			if !ctx.constraint.isValid() {
				return false
			}
			if pa == nil || !ctx.world.entities.SparseArray.Exist(key) {
				return true
			}
			pb := b.get(key)
			return pb == nil || yield(pa, view.FromPtr(pb))
		}
		findA := func(key EntityIndex) *A {
			if row, ok := group.rowOf(key); ok {
				return (*A)(group.cellPtr(col, row))
			}
			return a.getByIndex(key)
		}
		if totalB < totalA || (totalB == totalA && itB < itA) {
			for _, key := range b.EntityIndexes() {
				if !emit(key, findA(key)) {
					return
				}
			}
			return
		}
		for row := 0; row < group.Len(); row++ {
			if !emit(group.entities[row], (*A)(group.cellPtr(col, row))) {
				return
			}
		}
		for row, key := range a.EntityIndexes() {
			var pa *A
			if row < len(a.idx2Key) && a.idx2Key[row] == int32(key) {
				pa = &a.data[row]
			} else {
				pa = findA(key)
			}
			if !emit(key, pa) {
				return
			}
		}
	}
}
