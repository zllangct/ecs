package ecs

import "reflect"

type Collection[T ComponentObject, TP ComponentPointer[T]] struct {
	data []T
	ids  map[int64]int64
	seq  int64
	typ  reflect.Type
	len  int64
}

func NewCollection[T ComponentObject, TP ComponentPointer[T]]() *Collection[T, TP] {
	c := &Collection[T, TP]{
		ids: map[int64]int64{},
		typ: TypeOf[T](),
	}
	return c
}

func (c *Collection[T, TP]) getID() int64 {
	ok := false
	for !ok {
		c.seq++
		if _, exist := c.ids[c.seq]; !exist {
			break
		}
	}
	return c.seq
}

func (c *Collection[T, TP]) Add(element *T) (int64, *T) {
	idx := len(c.data)
	//Log.Info("collection Add:", ObjectToString(element))
	if int64(len(c.data)) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	id := c.getID()
	c.ids[id] = int64(idx)
	c.ids[int64(-idx)] = -id
	ret := TP(&(c.data[idx]))
	c.len++
	ret.setID(id)
	ret.setState(ComponentStateActive)
	return id, (*T)(ret)
}

func (c *Collection[T, TP]) Remove(id int64) {
	if id < 0 {
		return
	}
	idx, ok := c.ids[id]
	if !ok {
		return
	}
	//Log.Info("collection Remove:", ObjectToString(c.data[idx]))
	l := len(c.data)
	c.ids[c.ids[int64(l)]] = idx
	delete(c.ids, c.ids[-idx])
	c.ids[-idx] = c.ids[int64(l)]
	delete(c.ids, int64(l))
	c.data[idx], c.data[l-1] = c.data[l-1], c.data[idx]
	TP(&(c.data[l-1])).setState(ComponentStateDisable)
	c.shrink()
	c.len--
}

func (c *Collection[T, TP]) shrink() {
	var threshold int64
	if len(c.data) < 1024 {
		threshold = c.len * 2
	} else {
		threshold = int64(float64(c.len) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		c.data = c.data[:c.len]
	}
}

func (c *Collection[T, TP]) Get(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	return &(c.data[idx])
}

func (c *Collection[T, TP]) Len() int {
	return int(c.len)
}

func (c *Collection[T, TP]) EleType() reflect.Type {
	return c.typ
}
