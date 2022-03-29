package ecs

import (
	"reflect"
)

const (
	InitMaxSize = 1024 * 16
)

type ICollection interface {
	Len() int
	ElementType() reflect.Type

	getByIndex(idx int64) IComponent
}

type Collection[T ComponentObject, TP ComponentPointer[T]] struct {
	data []T
	ids  map[int64]int64
	seq  int64
	len  int64
}

func NewCollection[T ComponentObject, TP ComponentPointer[T]]() *Collection[T, TP] {
	size := InitMaxSize / TypeOf[T]().Size()
	c := &Collection[T, TP]{
		ids:  map[int64]int64{},
		data: make([]T, 0, size),
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

func (c *Collection[T, TP]) Add(element *T) *T {
	//Log.Info("collection Add:", ObjectToString(element))
	if int64(len(c.data)) > c.len {
		c.data[c.len] = *element
	} else {
		c.data = append(c.data, *element)
	}
	idx := c.len
	id := TP(element).ID()
	if id == 0 {
		id = c.getID()
		TP(element).setID(id)
	}
	c.ids[id] = idx
	c.ids[-idx] = -id
	ret := TP(&(c.data[idx]))
	c.len++
	return (*T)(ret)
}

func (c *Collection[T, TP]) Remove(id int64) *T {
	if id < 0 {
		return nil
	}
	idx, ok := c.ids[id]
	if !ok {
		return nil
	}
	lastIdx := c.len - 1
	lastId := -c.ids[-lastIdx]

	c.ids[lastId] = idx
	c.ids[-idx] = -lastId
	delete(c.ids, -lastIdx)
	delete(c.ids, id)

	c.data[idx], c.data[lastIdx] = c.data[lastIdx], c.data[idx]
	c.shrink()
	c.len--
	return &c.data[lastIdx]
}

func (c *Collection[T, TP]) RemoveAndReturn(id int64) *T {
	cpy := *c.Remove(id)
	return &cpy
}

func (c *Collection[T, TP]) shrink() {
	var threshold int64
	if len(c.data) < 1024 {
		threshold = c.len * 2
	} else {
		threshold = int64(float64(c.len) * 1.25)
	}
	if int64(len(c.data)) > threshold {
		c.data = c.data[:threshold]
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

func (c *Collection[T, TP]) getByIndex(idx int64) IComponent {
	return TP(&(c.data[idx]))
}

func (c *Collection[T, TP]) Len() int {
	return int(c.len)
}

func (c *Collection[T, TP]) ElementType() reflect.Type {
	return TypeOf[T]()
}
