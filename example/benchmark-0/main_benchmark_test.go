package main

import (
	"github.com/zllangct/ecs"
	"net/http"
	"testing"
	"time"
)

func BenchmarkNormal(b *testing.B) {
	game := &GameNormal{
		players: make(map[int64]*Player),
	}
	game.init()
	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.doFrame(false, uint64(i), delta)
		delta = time.Since(ts)
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
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	config.CollectionVersion = 1
	game.init(config)

	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.attack()
		doFrame(game.world, uint64(i), delta)
		delta = time.Since(ts)
	}
}
