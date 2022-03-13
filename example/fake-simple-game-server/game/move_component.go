package game

import "github.com/zllangct/ecs"

type Position struct {
	ecs.Component[Position, *Position]
	X int
	Y int
	Z int
}

type Movement struct {
	ecs.Component[Movement, *Movement]
	V   int
	Dir []int
}
