package game

import "github.com/zllangct/ecs"

type MoveChange struct {
	ecs.Component[MoveChange]
	V   int
	Dir []int
}
