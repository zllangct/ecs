package game

import "github.com/zllangct/ecs"

type MoveChange struct {
	ecs.DisposableComponent[MoveChange]
	Entity ecs.Entity
	V   int
	Dir []int
}
