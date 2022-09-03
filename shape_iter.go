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

func (s *ShapeIter[T]) tryNext() *T {
	skip := false
	find := false
	var p unsafe.Pointer
	var ec *EmptyComponent
	for i := s.offset; i < s.maxLen; i++ {
		p = s.indices.containers[s.mainKeyIndex].getPointer(int64(s.offset))
		ec = (*EmptyComponent)(p)
		*(**byte)(unsafe.Pointer((uintptr)(unsafe.Pointer(s.cur)) + s.indices.subOffset[s.mainKeyIndex])) = (*byte)(p)
		entity := ec.Owner()
		skip = false
		for i := 0; i < len(s.indices.subTypes); i++ {
			if i == s.mainKeyIndex {
				continue
			}
			subPointer := s.indices.containers[i].getPointerByEntity(entity)
			if subPointer == nil {
				skip = true
				break
			}
			*(**byte)(unsafe.Pointer((uintptr)(unsafe.Pointer(s.cur)) + s.indices.subOffset[i])) = (*byte)(subPointer)
		}
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
