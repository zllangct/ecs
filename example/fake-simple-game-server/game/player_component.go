package game

import "github.com/zllangct/ecs"

type PlayerComponent struct {
	ecs.Component[PlayerComponent, *PlayerComponent]
	Name      string
	Level     int
	SessionID int
}
