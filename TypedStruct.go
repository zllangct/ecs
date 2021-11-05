package ecs

func ToTypedParam(in interface{}) int {
	return 0
}

type TS1[T1 any] struct {
	A1 T1
}

func (t TS1[T1]) Expand() T1 {
	return t.A1
}

func (t TS1[T1]) Count() int {
	return 1
}

type TS2[T1, T2 any] struct {
	A1 T1
	A2 T2
}

type TS3[T1, T2, T3 any] struct {
	A1 T1
	A2 T2
	A3 T3
}

type TS4[T1, T2, T3, T4 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
}

type TS5[T1, T2, T3, T4, T5 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
	A5 T5
}

type TS6[T1, T2, T3, T4, T5, T6 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
	A5 T5
	A6 T6
}

type TS7[T1, T2, T3, T4, T5, T6, T7 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
	A5 T5
	A6 T6
	A7 T7
}

type TS8[T1, T2, T3, T4, T5, T6, T7, T8 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
	A5 T5
	A6 T6
	A7 T7
	A8 T8
}

type TS9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any] struct {
	A1 T1
	A2 T2
	A3 T3
	A4 T4
	A5 T5
	A6 T6
	A7 T7
	A8 T8
	A9 T9
}
