package main

import (
	"github.com/zllangct/ecs"
	_ "net/http/pprof"
	"testing"
	"time"
	_ "unsafe"
)

func TestFrame(t *testing.T) {
	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	game.init(config)

	var delta time.Duration
	var ts time.Time
	for i := 0; i < 10; i++ {
		ts = time.Now()
		doFrame(game.world, uint64(i), delta)
		game.attack()
		delta = time.Since(ts)
		//ecs.Log.Info("===== Frame:", i, "=====", delta)
	}
}

func TestEcsOptimizer(t *testing.T) {
	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	game.init(config)

	var frameInterval = time.Millisecond * 33
	var delta time.Duration
	var ts time.Time
	for i := 0; i < 10; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()
		doFrame(game.world, uint64(i), delta)
		game.attack()
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			game.world.Optimize(frameInterval - delta)
			time.Sleep(frameInterval - delta)
			delta = frameInterval
		}
	}
}
