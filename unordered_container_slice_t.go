package ecs

type TUnorderedContainer[T any] struct {
	data  []T
	len  int
}

func NewTUnorderedContainer[T any]() *TUnorderedContainer[T] {
	return &TUnorderedContainer[T]{
		data:  make([]T, 0),
		len:  0,
	}
}

func (p *TUnorderedContainer[T]) Add(item T) (int, *T) {
	p.data = append(p.data, item)
	p.len++
	return p.len - 1, &(p.data[p.len - 1])
}

func (p *TUnorderedContainer[T]) Remove(idx int) {
	if idx < 0 || idx >= p.len {
		return
	}
	p.data[idx] = p.data[p.len - 1]
	p.len -= 1
	p.data = p.data[:p.len - 1]
}

func (p *TUnorderedContainer[T]) Get(idx int) *T {
	if idx < 0 || idx >= p.len {
		return nil
	}
	return &(p.data[idx])
}

func (p *TUnorderedContainer[T]) Len() int {
	return p.len
}
