package ecs

type TUnorderedContainer[T any] []T

func NewTUnorderedContainer[T any]() *TUnorderedContainer[T] {
	return &TUnorderedContainer[T]{}

}

func (p *TUnorderedContainer[T]) Add(item T) (int, *T) {
	*p = append(*p, item)
	length := p.Len()
	return length, &(*p)[length-1]
}

func (p *TUnorderedContainer[T]) Remove(idx int) {
	length := p.Len()
	if idx < 0 || idx >= length {
		return
	}
	(*p)[idx] = (*p)[length-1]
	*p = (*p)[:length-1]
}

func (p *TUnorderedContainer[T]) Get(idx int) *T {
	length := p.Len()
	if idx < 0 || idx >= length {
		return nil
	}
	return &((*p)[idx])
}

func (p *TUnorderedContainer[T]) Len() int {
	return len(*p)
}

func (p *TUnorderedContainer[T]) Iterator() IIterator[T] {
	return NewIterator[T](p)
}
