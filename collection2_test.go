package ecs

import (
	"fmt"
	"testing"
)

func TestShapeCollection_Add(t *testing.T) {
	type Position struct {
		Component[Position, *Position]
		X int
		Y int
		Z int
	}

	type Movement struct {
		Component[Movement, *Movement]
		V   int
		Dir [3]int
	}
	p1 := &Position{X: 1, Y: 2, Z: 3}
	p1.setID(1)

	c := NewCollection2[Position]()
	ret := c.Add(p1)
	Log.Infof("p:%p : %+v", &ret, ret)
	ret1 := c.Remove(1)
	Log.Infof("p:%p : %+v", &ret1, ret1)

	p1.setID(2)
	c.Add(p1)
	p1.setID(3)
	c.Add(p1)
	p1.setID(4)
	c.Add(p1)
	p1.setID(5)
	c.Add(p1)
	p1.setID(6)
	c.Add(p1)
	p1.setID(7)
	c.Add(p1)

	ret = c.Remove(5)

	_ = ret

	iter := NewIterator2[Position](c)
	for com := iter.Begin(); !iter.End(); com = iter.Next() {
		fmt.Printf("iter 2: %+v\n", com)
	}
}
