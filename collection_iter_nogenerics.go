package ecs

import "reflect"

type IShapeIterator interface {
	Begin() []IComponent
	Val() []IComponent
	Next() []IComponent
	End() bool
	Empty() bool
}

type ShapeIter struct {
	c        ICollection
	len      int
	offset   int
	typ      []reflect.Type
	cur      []IComponent
	end      bool
	readOnly bool
}

func EmptyShapeIter() IShapeIterator {
	return &ShapeIter{}
}

func NewShapeIterator(collection ICollection, typ []reflect.Type, readOnly ...bool) IShapeIterator {
	iter := &ShapeIter{
		c:      collection,
		len:    collection.Len(),
		cur:    make([]IComponent, len(typ)),
		typ:    typ,
		offset: 0,
	}
	if len(readOnly) > 0 {
		iter.readOnly = readOnly[0]
	}
	if iter.len != 0 {
		if iter.readOnly {
			com := collection.getByIndex(0)
			iter.cur[0] = com.Clone()
			for i, t := range iter.typ {
				cc := com.Owner().getComponentByTypeInSystem(t)
				if cc != nil {
					iter.cur[i] = cc.Clone()
				} else {
					iter.cur[i] = nil
				}
			}
		} else {
			com := collection.getByIndex(0)
			iter.cur[0] = com
			for i, t := range iter.typ {
				iter.cur[i] = com.Owner().getComponentByTypeInSystem(t)
			}
		}
	}

	return iter
}

func (i *ShapeIter) Empty() bool {
	if i.len == 0 {
		return true
	}
	return false
}

func (i *ShapeIter) End() bool {
	if i.offset > i.len-1 || i.len == 0 {
		return true
	}
	return false
}

func (i *ShapeIter) Begin() []IComponent {
	if i.len != 0 {
		i.offset = 0
		if i.readOnly {
			com := i.c.getByIndex(int64(i.offset))
			i.cur[0] = com.Clone()
			for idx, t := range i.typ {
				cc := com.Owner().getComponentByTypeInSystem(t)
				if cc != nil {
					i.cur[idx] = cc.Clone()
				} else {
					i.cur[idx] = nil
				}
			}
		} else {
			com := i.c.getByIndex(int64(i.offset))
			i.cur[0] = com
			for idx, t := range i.typ {
				i.cur[idx] = com.Owner().getComponentByTypeInSystem(t)
			}
		}
	}
	return i.cur
}

func (i *ShapeIter) Val() []IComponent {
	return i.cur
}

func (i *ShapeIter) Next() []IComponent {
	i.offset++
	if !i.End() {
		if i.readOnly {
			com := i.c.getByIndex(int64(i.offset))
			i.cur[0] = com.Clone()
			for idx, t := range i.typ {
				cc := com.Owner().getComponentByTypeInSystem(t)
				if cc != nil {
					i.cur[idx] = cc.Clone()
				} else {
					i.cur[idx] = nil
				}
			}
		} else {
			com := i.c.getByIndex(int64(i.offset))
			i.cur[0] = com
			for idx, t := range i.typ {
				i.cur[idx] = com.Owner().getComponentByTypeInSystem(t)
			}
		}
	}
	return i.cur
}
