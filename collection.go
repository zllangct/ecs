package ecs

import "reflect"

type Collection[T any] struct {
	data []T
	ids  map[int64]int64
	seq  int64
	typ  reflect.Type
}

func NewCollection[T any]() *Collection[T] {
	c := &Collection[T]{
		ids: map[int64]int64{},
		typ: reflect.TypeOf(*new(T)),
	}
	return c
}

func (c *Collection[T]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *Collection[T]) Add(element *T) (int64, *T) {
	idx := len(c.data)
	Log.Info("collection Add:", ObjectToString(element))
	c.data = append(c.data, *element)
	id := c.getID()
	c.ids[id] = int64(idx)
	c.ids[int64(-idx)] = -id
	ret := &(c.data[idx])
	//ss := (*T)(unsafe.Pointer(ret))
	//c.data[idx].setID(id)
	return id, ret
}

func (c *Collection[T]) Remove(id int64) {
	if id < 0 {
		return
	}
	idx, ok := c.ids[id]
	if !ok {
		return
	}
	l := len(c.data)
	c.ids[c.ids[int64(l)]] = idx
	delete(c.ids, c.ids[-idx])
	c.ids[-idx] = c.ids[int64(l)]
	delete(c.ids, int64(l))

	c.data[idx], c.data[l-1] = c.data[l-1], c.data[idx]
	c.data = c.data[:l-1]
}

func (c *Collection[T]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return &(c.data[idx])
}

func (c *Collection[T]) Len() int {
	return len(c.data)
}

func (c *Collection[T]) EleType() reflect.Type {
	return c.typ
}
