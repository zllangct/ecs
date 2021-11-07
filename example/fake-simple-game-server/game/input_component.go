package game

import "github.com/zllangct/ecs"

type MoveChange struct {
	ecs.DisposableComponent[MoveChange]
	V   int
	Dir []int
}
