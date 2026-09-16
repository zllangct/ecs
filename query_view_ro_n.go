package ecs

import (
	"iter"
	"unsafe"
)

// ============================================================================
// View3RO-View8RO：组表多列 lockstep 迭代（首组件可写 + 其余只读视图）。
// 语义与 View2RO 一致：A 须非只读依赖，其余须只读依赖；
// 同组表走 lockstep，否则回退通用路径（GetBuddyReadOnly）。
// 产出为 TupleN（range 仅支持 2 变量）。
// ============================================================================

// viewRODeps 校验：首类型非只读依赖，其余类型只读依赖
func viewRODeps(ctx *SystemContext, itA ComponentIntType, itsRO ...ComponentIntType) bool {
	depA, ok := ctx.info.getDep(itA)
	if !ok || depA.readonly() {
		return false
	}
	for _, it := range itsRO {
		dep, ok := ctx.info.getDep(it)
		if !ok || !dep.readonly() {
			return false
		}
	}
	return true
}

func View3RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple3[*A, BR, CR]] {
	empty := func(yield func(Tuple3[*A, BR, CR]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC); ok {
		return func(yield func(Tuple3[*A, BR, CR]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple3[*A, BR, CR]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC)))
	return func(yield func(Tuple3[*A, BR, CR]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			if !k1 || !k2 || !k3 {
				continue
			}
			if !yield(Tuple3[*A, BR, CR]{V1: a, V2: b, V3: c}) {
				return
			}
		}
	}
}

func View4RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], DR ReadOnlyView[DR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple4[*A, BR, CR, DR]] {
	empty := func(yield func(Tuple4[*A, BR, CR, DR]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	var zeroD DR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	itD := ComponentIntType(zeroD.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC, itD) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD); ok {
		return func(yield func(Tuple4[*A, BR, CR, DR]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple4[*A, BR, CR, DR]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: zeroD.FromPtr(unsafe.Pointer(&datas[3][row*sizes[3]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC, itD)))
	return func(yield func(Tuple4[*A, BR, CR, DR]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			d, k4 := ctx.GetBuddyReadOnly[DR](index)
			if !k1 || !k2 || !k3 || !k4 {
				continue
			}
			if !yield(Tuple4[*A, BR, CR, DR]{V1: a, V2: b, V3: c, V4: d}) {
				return
			}
		}
	}
}

func View5RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], DR ReadOnlyView[DR], ER ReadOnlyView[ER], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple5[*A, BR, CR, DR, ER]] {
	empty := func(yield func(Tuple5[*A, BR, CR, DR, ER]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	var zeroD DR
	var zeroE ER
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	itD := ComponentIntType(zeroD.ComponentPacketIdentifier())
	itE := ComponentIntType(zeroE.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC, itD, itE) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE); ok {
		return func(yield func(Tuple5[*A, BR, CR, DR, ER]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple5[*A, BR, CR, DR, ER]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: zeroD.FromPtr(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: zeroE.FromPtr(unsafe.Pointer(&datas[4][row*sizes[4]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC, itD, itE)))
	return func(yield func(Tuple5[*A, BR, CR, DR, ER]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			d, k4 := ctx.GetBuddyReadOnly[DR](index)
			e, k5 := ctx.GetBuddyReadOnly[ER](index)
			if !k1 || !k2 || !k3 || !k4 || !k5 {
				continue
			}
			if !yield(Tuple5[*A, BR, CR, DR, ER]{V1: a, V2: b, V3: c, V4: d, V5: e}) {
				return
			}
		}
	}
}

func View6RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], DR ReadOnlyView[DR], ER ReadOnlyView[ER], FR ReadOnlyView[FR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple6[*A, BR, CR, DR, ER, FR]] {
	empty := func(yield func(Tuple6[*A, BR, CR, DR, ER, FR]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	var zeroD DR
	var zeroE ER
	var zeroF FR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	itD := ComponentIntType(zeroD.ComponentPacketIdentifier())
	itE := ComponentIntType(zeroE.ComponentPacketIdentifier())
	itF := ComponentIntType(zeroF.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC, itD, itE, itF) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF); ok {
		return func(yield func(Tuple6[*A, BR, CR, DR, ER, FR]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple6[*A, BR, CR, DR, ER, FR]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: zeroD.FromPtr(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: zeroE.FromPtr(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: zeroF.FromPtr(unsafe.Pointer(&datas[5][row*sizes[5]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC, itD, itE, itF)))
	return func(yield func(Tuple6[*A, BR, CR, DR, ER, FR]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			d, k4 := ctx.GetBuddyReadOnly[DR](index)
			e, k5 := ctx.GetBuddyReadOnly[ER](index)
			f, k6 := ctx.GetBuddyReadOnly[FR](index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 {
				continue
			}
			if !yield(Tuple6[*A, BR, CR, DR, ER, FR]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f}) {
				return
			}
		}
	}
}

func View7RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], DR ReadOnlyView[DR], ER ReadOnlyView[ER], FR ReadOnlyView[FR], GR ReadOnlyView[GR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple7[*A, BR, CR, DR, ER, FR, GR]] {
	empty := func(yield func(Tuple7[*A, BR, CR, DR, ER, FR, GR]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	var zeroD DR
	var zeroE ER
	var zeroF FR
	var zeroG GR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	itD := ComponentIntType(zeroD.ComponentPacketIdentifier())
	itE := ComponentIntType(zeroE.ComponentPacketIdentifier())
	itF := ComponentIntType(zeroF.ComponentPacketIdentifier())
	itG := ComponentIntType(zeroG.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC, itD, itE, itF, itG) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF, itG); ok {
		return func(yield func(Tuple7[*A, BR, CR, DR, ER, FR, GR]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple7[*A, BR, CR, DR, ER, FR, GR]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: zeroD.FromPtr(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: zeroE.FromPtr(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: zeroF.FromPtr(unsafe.Pointer(&datas[5][row*sizes[5]])),
					V7: zeroG.FromPtr(unsafe.Pointer(&datas[6][row*sizes[6]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC, itD, itE, itF, itG)))
	return func(yield func(Tuple7[*A, BR, CR, DR, ER, FR, GR]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			d, k4 := ctx.GetBuddyReadOnly[DR](index)
			e, k5 := ctx.GetBuddyReadOnly[ER](index)
			f, k6 := ctx.GetBuddyReadOnly[FR](index)
			g, k7 := ctx.GetBuddyReadOnly[GR](index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 || !k7 {
				continue
			}
			if !yield(Tuple7[*A, BR, CR, DR, ER, FR, GR]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f, V7: g}) {
				return
			}
		}
	}
}

func View8RO[A any, BR ReadOnlyView[BR], CR ReadOnlyView[CR], DR ReadOnlyView[DR], ER ReadOnlyView[ER], FR ReadOnlyView[FR], GR ReadOnlyView[GR], HR ReadOnlyView[HR], PA ComponentPointer[A]](ctx *SystemContext) iter.Seq[Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]] {
	empty := func(yield func(Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA := GetIntType[A, PA]()
	var zeroB BR
	var zeroC CR
	var zeroD DR
	var zeroE ER
	var zeroF FR
	var zeroG GR
	var zeroH HR
	itB := ComponentIntType(zeroB.ComponentPacketIdentifier())
	itC := ComponentIntType(zeroC.ComponentPacketIdentifier())
	itD := ComponentIntType(zeroD.ComponentPacketIdentifier())
	itE := ComponentIntType(zeroE.ComponentPacketIdentifier())
	itF := ComponentIntType(zeroF.ComponentPacketIdentifier())
	itG := ComponentIntType(zeroG.ComponentPacketIdentifier())
	itH := ComponentIntType(zeroH.ComponentPacketIdentifier())
	if !viewRODeps(ctx, itA, itB, itC, itD, itE, itF, itG, itH) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF, itG, itH); ok {
		return func(yield func(Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: zeroB.FromPtr(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: zeroC.FromPtr(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: zeroD.FromPtr(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: zeroE.FromPtr(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: zeroF.FromPtr(unsafe.Pointer(&datas[5][row*sizes[5]])),
					V7: zeroG.FromPtr(unsafe.Pointer(&datas[6][row*sizes[6]])),
					V8: zeroH.FromPtr(unsafe.Pointer(&datas[7][row*sizes[7]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithBuddies(*NewQueryBuddies(itB, itC, itD, itE, itF, itG, itH)))
	return func(yield func(Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := ctx.GetBuddyReadOnly[BR](index)
			c, k3 := ctx.GetBuddyReadOnly[CR](index)
			d, k4 := ctx.GetBuddyReadOnly[DR](index)
			e, k5 := ctx.GetBuddyReadOnly[ER](index)
			f, k6 := ctx.GetBuddyReadOnly[FR](index)
			g, k7 := ctx.GetBuddyReadOnly[GR](index)
			h, k8 := ctx.GetBuddyReadOnly[HR](index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 || !k7 || !k8 {
				continue
			}
			if !yield(Tuple8[*A, BR, CR, DR, ER, FR, GR, HR]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f, V7: g, V8: h}) {
				return
			}
		}
	}
}
