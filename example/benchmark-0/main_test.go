package main

import (
	"testing"
	"time"
	_ "unsafe"
)

func BenchmarkNormal(b *testing.B) {
	game := &GameNormal{
		players: make(map[int64]*Player),
	}
	game.init()
	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.doFrame(false, uint64(i), frameInterval)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}

func BenchmarkNormalParallel(b *testing.B) {
	game := &GameNormal{
		players: make(map[int64]*Player),
	}
	game.init()
	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.doFrame(true, uint64(i), frameInterval)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}

func BenchmarkEcs(b *testing.B) {
	game := &GameECS{}
	game.init()

	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < b.N; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()
		game.attack()
		doFrame(game.world, uint64(i), frameInterval)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}

func TestEcs(t *testing.T) {
	game := &GameECS{}
	game.init()

	var delta time.Duration
	var ts time.Time
	var frameInterval time.Duration = time.Millisecond * 33
	for i := 0; i < 10; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()
		doFrame(game.world, uint64(i), frameInterval)
		game.attack()
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			delta = frameInterval
		}
	}
}
