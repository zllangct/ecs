package ecs

import "unsafe"

type IShapeIterator[T any] interface {
	Begin() *T
	Val() *T
	Next() *T
	End() bool
}

type ShapeIter[T any] struct {
	indices      ShapeIndices
	maxLen       int
	offset       int
	begin        int
	mainKeyIndex int
	cur          *T
}

func EmptyShapeIter[T any]() IShapeIterator[T] {
	return &ShapeIter[T]{}
}

func NewShapeIterator[T any](indices ShapeIndices, mainKeyIndex int) IShapeIterator[T] {
	iter := &ShapeIter[T]{
		indices:      indices,
		maxLen:       indices.containers[mainKeyIndex].Len(),
		mainKeyIndex: mainKeyIndex,
		offset:       0,
	}

	return iter
}

// TODO 热点
func (s *ShapeIter[T]) tryNext() *T {
	skip := false
	find := false
	var p unsafe.Pointer
	var ec *EmptyComponent
	for i := s.offset; i < s.maxLen; i++ {
		//TODO check if this is the best way to do this
		p = s.indices.containers[s.mainKeyIndex].getPointerByIndex(int64(s.offset))
		ec = (*EmptyComponent)(p)
		if s.indices.readOnly[s.mainKeyIndex] {
			*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.indices.subOffset[s.mainKeyIndex])) = &(*(*byte)(p))
		} else {
			*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.indices.subOffset[s.mainKeyIndex])) = (*byte)(p)
		}
		entity := ec.Owner()
		skip = s.getSiblings(entity)
		if !skip {
			s.offset = i
			find = true
			break
		}
	}
	if !find {
		s.cur = nil
	}

	return s.cur
}

func (s *ShapeIter[T]) getSiblings(entity Entity) bool {
	for i := 0; i < len(s.indices.subTypes); i++ {
		if i == s.mainKeyIndex {
			continue
		}
		subPointer := s.indices.containers[i].getPointerByEntity(entity)
		if subPointer == nil {
			return true
		}
		s.trans(i, subPointer)
	}
	return false
}

func (s *ShapeIter[T]) trans(i int, subPointer unsafe.Pointer) {
	if s.indices.readOnly[i] {
		*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.indices.subOffset[i])) = &(*(*byte)(subPointer))
	} else {
		*(**byte)(unsafe.Add(unsafe.Pointer(s.cur), s.indices.subOffset[i])) = (*byte)(subPointer)
	}
}

func (s *ShapeIter[T]) End() bool {
	if s.cur == nil {
		return true
	}
	return false
}

func (s *ShapeIter[T]) Begin() *T {
	if s.maxLen != 0 {
		s.offset = 0
		s.cur = new(T)
		s.tryNext()
	}
	return s.cur
}

func (s *ShapeIter[T]) Val() *T {
	if s.cur == nil || !s.End() {
		s.Begin()
	}
	return s.cur
}

func (s *ShapeIter[T]) Next() *T {
	s.offset++
	if !s.End() {
		s.tryNext()
	} else {
		s.cur = nil
	}
	return s.cur
}
