package ecs

import (
	"iter"
	"unsafe"
)

// Tuple3-Tuple8 多元组：Go 的 range-over-func 最多支持 2 个迭代变量，
// 3 元及以上的 View 以元组结构体产出（指针/视图值，拷贝无分配）。
type Tuple3[V1, V2, V3 any] struct {
	V1 V1
	V2 V2
	V3 V3
}

type Tuple4[V1, V2, V3, V4 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
}

type Tuple5[V1, V2, V3, V4, V5 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
	V5 V5
}

type Tuple6[V1, V2, V3, V4, V5, V6 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
	V5 V5
	V6 V6
}

type Tuple7[V1, V2, V3, V4, V5, V6, V7 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
	V5 V5
	V6 V6
	V7 V7
}

type Tuple8[V1, V2, V3, V4, V5, V6, V7, V8 any] struct {
	V1 V1
	V2 V2
	V3 V3
	V4 V4
	V5 V5
	V6 V6
	V7 V7
	V8 V8
}

// viewDepsWritable 全部类型均已声明为非只读依赖
func viewDepsWritable(ctx *SystemContext, its ...ComponentIntType) bool {
	for _, it := range its {
		dep, ok := ctx.info.getDep(it)
		if !ok || dep.readonly() {
			return false
		}
	}
	return true
}

// lockstepColumns 全部类型同属一张组表时，解析出各列数据/步长/行数
func lockstepColumns(ctx *SystemContext, its ...ComponentIntType) (datas [][]byte, sizes []int, rows int, ok bool) {
	a, ok := ctx.world.archetypes.owner(its[0])
	if !ok {
		return nil, nil, 0, false
	}
	datas = make([][]byte, len(its))
	sizes = make([]int, len(its))
	for i, it := range its {
		ai, ok := ctx.world.archetypes.owner(it)
		if !ok || ai != a {
			return nil, nil, 0, false
		}
		col, _ := a.colOf(it)
		datas[i] = a.columns[col]
		sizes[i] = a.sizes[col]
	}
	return datas, sizes, a.Len(), true
}

// ============================================================================
// View3-View8：组表多列 lockstep 迭代（全部可写依赖）。
// 语义与 View2 一致：同组表走 lockstep，否则回退 Query 通用路径；
// 依赖未声明/只读时静默返回空。产出为 TupleN（range 仅支持 2 变量）。
// ============================================================================

func View3[A, B, C any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C]](ctx *SystemContext) iter.Seq[Tuple3[*A, *B, *C]] {
	empty := func(yield func(Tuple3[*A, *B, *C]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC]()
	if !viewDepsWritable(ctx, itA, itB, itC) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC); ok {
		return func(yield func(Tuple3[*A, *B, *C]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple3[*A, *B, *C]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC]())
	return func(yield func(Tuple3[*A, *B, *C]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			if !k1 || !k2 || !k3 {
				continue
			}
			if !yield(Tuple3[*A, *B, *C]{V1: a, V2: b, V3: c}) {
				return
			}
		}
	}
}

func View4[A, B, C, D any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C], PD ComponentPointer[D]](ctx *SystemContext) iter.Seq[Tuple4[*A, *B, *C, *D]] {
	empty := func(yield func(Tuple4[*A, *B, *C, *D]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC, itD := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC](), GetIntType[D, PD]()
	if !viewDepsWritable(ctx, itA, itB, itC, itD) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD); ok {
		return func(yield func(Tuple4[*A, *B, *C, *D]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple4[*A, *B, *C, *D]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: (*D)(unsafe.Pointer(&datas[3][row*sizes[3]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC](), WithComp[D, PD]())
	return func(yield func(Tuple4[*A, *B, *C, *D]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			d, k4 := QueryGet[D, PD](&q, index)
			if !k1 || !k2 || !k3 || !k4 {
				continue
			}
			if !yield(Tuple4[*A, *B, *C, *D]{V1: a, V2: b, V3: c, V4: d}) {
				return
			}
		}
	}
}

func View5[A, B, C, D, E any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C], PD ComponentPointer[D], PE ComponentPointer[E]](ctx *SystemContext) iter.Seq[Tuple5[*A, *B, *C, *D, *E]] {
	empty := func(yield func(Tuple5[*A, *B, *C, *D, *E]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC, itD, itE := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC](), GetIntType[D, PD](), GetIntType[E, PE]()
	if !viewDepsWritable(ctx, itA, itB, itC, itD, itE) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE); ok {
		return func(yield func(Tuple5[*A, *B, *C, *D, *E]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple5[*A, *B, *C, *D, *E]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: (*D)(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: (*E)(unsafe.Pointer(&datas[4][row*sizes[4]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC](), WithComp[D, PD](), WithComp[E, PE]())
	return func(yield func(Tuple5[*A, *B, *C, *D, *E]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			d, k4 := QueryGet[D, PD](&q, index)
			e, k5 := QueryGet[E, PE](&q, index)
			if !k1 || !k2 || !k3 || !k4 || !k5 {
				continue
			}
			if !yield(Tuple5[*A, *B, *C, *D, *E]{V1: a, V2: b, V3: c, V4: d, V5: e}) {
				return
			}
		}
	}
}

func View6[A, B, C, D, E, F any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C], PD ComponentPointer[D], PE ComponentPointer[E], PF ComponentPointer[F]](ctx *SystemContext) iter.Seq[Tuple6[*A, *B, *C, *D, *E, *F]] {
	empty := func(yield func(Tuple6[*A, *B, *C, *D, *E, *F]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC, itD, itE, itF := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC](), GetIntType[D, PD](), GetIntType[E, PE](), GetIntType[F, PF]()
	if !viewDepsWritable(ctx, itA, itB, itC, itD, itE, itF) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF); ok {
		return func(yield func(Tuple6[*A, *B, *C, *D, *E, *F]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple6[*A, *B, *C, *D, *E, *F]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: (*D)(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: (*E)(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: (*F)(unsafe.Pointer(&datas[5][row*sizes[5]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC](), WithComp[D, PD](), WithComp[E, PE](), WithComp[F, PF]())
	return func(yield func(Tuple6[*A, *B, *C, *D, *E, *F]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			d, k4 := QueryGet[D, PD](&q, index)
			e, k5 := QueryGet[E, PE](&q, index)
			f, k6 := QueryGet[F, PF](&q, index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 {
				continue
			}
			if !yield(Tuple6[*A, *B, *C, *D, *E, *F]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f}) {
				return
			}
		}
	}
}

func View7[A, B, C, D, E, F, G any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C], PD ComponentPointer[D], PE ComponentPointer[E], PF ComponentPointer[F], PG ComponentPointer[G]](ctx *SystemContext) iter.Seq[Tuple7[*A, *B, *C, *D, *E, *F, *G]] {
	empty := func(yield func(Tuple7[*A, *B, *C, *D, *E, *F, *G]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC, itD, itE, itF, itG := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC](), GetIntType[D, PD](), GetIntType[E, PE](), GetIntType[F, PF](), GetIntType[G, PG]()
	if !viewDepsWritable(ctx, itA, itB, itC, itD, itE, itF, itG) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF, itG); ok {
		return func(yield func(Tuple7[*A, *B, *C, *D, *E, *F, *G]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple7[*A, *B, *C, *D, *E, *F, *G]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: (*D)(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: (*E)(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: (*F)(unsafe.Pointer(&datas[5][row*sizes[5]])),
					V7: (*G)(unsafe.Pointer(&datas[6][row*sizes[6]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC](), WithComp[D, PD](), WithComp[E, PE](), WithComp[F, PF](), WithComp[G, PG]())
	return func(yield func(Tuple7[*A, *B, *C, *D, *E, *F, *G]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			d, k4 := QueryGet[D, PD](&q, index)
			e, k5 := QueryGet[E, PE](&q, index)
			f, k6 := QueryGet[F, PF](&q, index)
			g, k7 := QueryGet[G, PG](&q, index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 || !k7 {
				continue
			}
			if !yield(Tuple7[*A, *B, *C, *D, *E, *F, *G]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f, V7: g}) {
				return
			}
		}
	}
}

func View8[A, B, C, D, E, F, G, H any, PA ComponentPointer[A], PB ComponentPointer[B], PC ComponentPointer[C], PD ComponentPointer[D], PE ComponentPointer[E], PF ComponentPointer[F], PG ComponentPointer[G], PH ComponentPointer[H]](ctx *SystemContext) iter.Seq[Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]] {
	empty := func(yield func(Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]) bool) {}
	if !ctx.constraint.isValid() {
		return empty
	}
	itA, itB, itC, itD, itE, itF, itG, itH := GetIntType[A, PA](), GetIntType[B, PB](), GetIntType[C, PC](), GetIntType[D, PD](), GetIntType[E, PE](), GetIntType[F, PF](), GetIntType[G, PG](), GetIntType[H, PH]()
	if !viewDepsWritable(ctx, itA, itB, itC, itD, itE, itF, itG, itH) {
		return empty
	}
	if datas, sizes, rows, ok := lockstepColumns(ctx, itA, itB, itC, itD, itE, itF, itG, itH); ok {
		return func(yield func(Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]) bool) {
			for row := 0; row < rows; row++ {
				if !yield(Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]{
					V1: (*A)(unsafe.Pointer(&datas[0][row*sizes[0]])),
					V2: (*B)(unsafe.Pointer(&datas[1][row*sizes[1]])),
					V3: (*C)(unsafe.Pointer(&datas[2][row*sizes[2]])),
					V4: (*D)(unsafe.Pointer(&datas[3][row*sizes[3]])),
					V5: (*E)(unsafe.Pointer(&datas[4][row*sizes[4]])),
					V6: (*F)(unsafe.Pointer(&datas[5][row*sizes[5]])),
					V7: (*G)(unsafe.Pointer(&datas[6][row*sizes[6]])),
					V8: (*H)(unsafe.Pointer(&datas[7][row*sizes[7]])),
				}) {
					return
				}
			}
		}
	}
	q := ctx.NewQuery(WithComp[A, PA](), WithComp[B, PB](), WithComp[C, PC](), WithComp[D, PD](), WithComp[E, PE](), WithComp[F, PF](), WithComp[G, PG](), WithComp[H, PH]())
	return func(yield func(Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]) bool) {
		for index := range q.Iter() {
			a, k1 := QueryGet[A, PA](&q, index)
			b, k2 := QueryGet[B, PB](&q, index)
			c, k3 := QueryGet[C, PC](&q, index)
			d, k4 := QueryGet[D, PD](&q, index)
			e, k5 := QueryGet[E, PE](&q, index)
			f, k6 := QueryGet[F, PF](&q, index)
			g, k7 := QueryGet[G, PG](&q, index)
			h, k8 := QueryGet[H, PH](&q, index)
			if !k1 || !k2 || !k3 || !k4 || !k5 || !k6 || !k7 || !k8 {
				continue
			}
			if !yield(Tuple8[*A, *B, *C, *D, *E, *F, *G, *H]{V1: a, V2: b, V3: c, V4: d, V5: e, V6: f, V7: g, V8: h}) {
				return
			}
		}
	}
}
