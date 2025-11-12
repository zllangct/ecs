package ecs

import (
	"errors"
	"iter"
	"slices"
)

type SystemTraverser interface {
	count() int
	all() []SystemInfo
	independentGroups() iter.Seq[[]SystemInfo]
	add(sys SystemInfo) error
	remove(sys SystemInfo) error
	setOrder(order Order)
	getOrder() Order
	getBatchCount() int
	getMaxCountPeerBatch() int
}

type SystemDefaultList struct {
	order   Order
	systems []SystemInfo
}

func NewSystemDefaultList() SystemTraverser {
	return &SystemDefaultList{}
}

func (s *SystemDefaultList) count() int {
	return len(s.systems)
}

func (s *SystemDefaultList) all() []SystemInfo {
	return s.systems
}

func (s *SystemDefaultList) getBatchCount() int {
	return 1
}

func (s *SystemDefaultList) getMaxCountPeerBatch() int {
	return len(s.systems)
}

func (s *SystemDefaultList) independentGroups() iter.Seq[[]SystemInfo] {
	return func(yield func([]SystemInfo) bool) {
		yield(s.all())
	}
}

func (s *SystemDefaultList) getOrder() Order {
	return s.order
}

func (s *SystemDefaultList) setOrder(order Order) {
	s.order = order
}

func (s *SystemDefaultList) add(sys SystemInfo) error {
	if slices.Contains(s.systems, sys) {
		return errors.New("repeated system")
	}

	s.systems = append(s.systems, sys)
	return nil
}

func (s *SystemDefaultList) remove(sys SystemInfo) error {
	index := slices.Index(s.systems, sys)
	if index == -1 {
		return errors.New("system not found in the system list")
	}
	s.systems = append(s.systems[:index], s.systems[index+1:]...)
	return nil
}
