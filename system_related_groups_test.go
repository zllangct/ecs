package ecs

import (
	"strconv"
	"testing"
)

func genTestCase(sg SystemTraverser) {
	handler := func(ctx *SystemContext, event Event) error { return nil }
	args := [][]struct {
		i ComponentIntType
		b bool
	}{
		{
			{1, false},
			{2, false},
		},
		{
			{2, true},
			{3, false},
		},
		{
			{3, false},
			{5, false},
		},
		{
			{4, false},
			{3, false},
			{6, false},
		},
		{
			{7, false},
		},
		{
			{9, false},
			{10, false},
		},
		{
			{6, false},
		},
		{
			{1, false},
			{5, false},
		},
		{
			{4, false},
			{6, false},
		},
		{
			{7, false},
			{5, false},
		},
		{
			{1, true},
		},
	}
	for i, deps := range args {
		var compDeps []Dependency
		for _, dep := range deps {
			compDeps = append(compDeps, NewDependency(dep.i, dep.b))
		}
		s := newSystem(nil, handler, SystemTypeLight)
		s.init(WithDeps(compDeps...), WithName(strconv.Itoa(i)))
		sg.add(s)
	}
}

func TestNewSystemGroupIterEmpty(t *testing.T) {
	sg := NewSystemRelatedGroups()
	t.Logf("========== system count %d, Batch count: %d, Max peer Batch: %d:", sg.count(), sg.getBatchCount(), sg.getMaxCountPeerBatch())

	for ss := range sg.independentGroups() {
		t.Logf("========== batch:")
		for _, s := range ss {
			t.Logf("%s\n", s.name())
		}
	}
}

func TestNewSystemGroupIter(t *testing.T) {
	sg := NewSystemRelatedGroups()
	genTestCase(sg)
	t.Logf("========== system count %d, Batch count: %d, Max peer Batch: %d:", sg.count(), sg.getBatchCount(), sg.getMaxCountPeerBatch())

	for ss := range sg.independentGroups() {
		t.Logf("========== batch: count: %d", len(ss))
		for _, s := range ss {
			t.Logf("%s\n", s.name())
		}
	}
}

func BenchmarkSystemGroupIter(b *testing.B) {
	sg := NewSystemRelatedGroups()
	genTestCase(sg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for ss := range sg.independentGroups() {
			for _, s := range ss {
				_ = s
			}
		}
	}
}
