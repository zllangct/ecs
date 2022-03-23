package ecs

import (
	"fmt"
	"reflect"
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
	m1 := &Movement{V: 1, Dir: [3]int{1, 2, 3}}

	shp1 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp1.entity = 1

	shp2 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp2.entity = 2

	shp3 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp3.entity = 3

	shp4 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp4.entity = 4

	shp5 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp5.entity = 5

	shp6 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp6.entity = 6

	shp7 := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp7.entity = 7

	Log.Infof("p:%p : %+v", &shp1, shp1)
	c := NewShapeCollection[Shape2[Position, Movement]]([]reflect.Type{TypeOf[Position](), TypeOf[Movement]()})
	ret := c.Add(shp1)
	Log.Infof("p:%p : %+v", &ret, ret)
	ret1 := c.RemoveAndReturn(1)
	Log.Infof("p:%p : %+v", &ret1, ret1)

	c.Add(shp2)
	c.Add(shp3)
	c.Add(shp4)
	c.Add(shp5)
	c.Add(shp6)
	c.Add(shp7)

	ret = c.RemoveAndReturn(5) //TODO 有bug

	iter := NewShapeIterator[Shape2[Position, Movement]](c)
	for shp := iter.Begin(); !iter.End(); shp = iter.Next() {
		fmt.Printf("style 2: %+v, C1:%+v, C2: %+v\n", shp, shp.C1, shp.C2)
	}
}
