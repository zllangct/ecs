package ecs

import (
	"testing"
)

func TestUniqueID(t *testing.T) {
	m := make(map[int64]struct{})
	count := 0
	for i := 0; i < 5000000; i++ {
		id := UniqueID()
		if _, ok := m[id]; ok {
			count += 1
			println("repeat:", count, id)
		} else {
			m[id] = struct{}{}
		}
	}
}

func BenchmarkUniqueID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := UniqueID()
		_ = id
	}
}
