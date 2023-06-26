package ecs

import (
	"reflect"
	"testing"
)

func BenchmarkIntTypeAndReflectType(b *testing.B) {
	infos := NewSparseArray[uint16, ComponentMetaInfo]()
	infos.Add(0, &ComponentMetaInfo{})

	m1 := map[uint16]struct{}{1: {}}
	m2 := map[reflect.Type]struct{}{reflect.TypeOf(0): {}}

	b.Run("int type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typ := *infos.Get(0)
			_ = m1[typ.it]
		}
	})
	b.Run("reflect type", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			typ := reflect.TypeOf(0)
			_ = m2[typ]
		}
	})
}
