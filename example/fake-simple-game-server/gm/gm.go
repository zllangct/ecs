package gm

import (
	"context"
	"math/rand"
	"test_ecs_fake_server/game"
	"time"
)

type GM struct {
	game *game.FakeGame
}

func NewGM() *GM {
	return &GM{}
}

func (g *GM) Run(ctx context.Context, game *game.FakeGame) {
	g.game = game

	timeScale := 0
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
		if timeScale == 0 {
			g.ChangeMovementTimeScale(1.2)
			timeScale = 1
		} else {
			g.ChangeMovementTimeScale(1.0)
			timeScale = 0
		}
	}
}

func (g *GM) ChangeMovementTimeScale(timeScale float64) {
	g.game.ChangeMovementTimeScale(timeScale)
}
