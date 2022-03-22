package ecs

import (
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

	shp := &Shape2[Position, Movement]{
		C1: p1,
		C2: m1,
	}
	shp.entity = 1

	Log.Infof("p:%p : %+v", &shp, shp)
	c := NewShapeCollection[Shape2[Position, Movement]]([]reflect.Type{TypeOf[Position](), TypeOf[Movement]()})
	ret := c.Add(shp, shp.entity)
	Log.Infof("p:%p : %+v", &ret, ret)
	ret1 := c.RemoveAndReturn(shp.entity)
	Log.Infof("p:%p : %+v", &ret1, ret1)
}
