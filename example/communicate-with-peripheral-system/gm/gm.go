package gm

import (
	"context"
	"test_ecs/game"
)

type GM struct {
	game *game.FakeGame
}

func NewGM() *GM {
	return &GM{}
}

func (g *GM) Run(ctx context.Context, game *game.FakeGame) {
	g.game = game
}

func (g *GM) ChangeMovementTimeScale() {
	g.game.ChangeMovementTimeScale()
}
