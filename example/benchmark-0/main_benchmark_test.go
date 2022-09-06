package main

import (
	"github.com/zllangct/ecs"
	"net/http"
	runtime2 "runtime"
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

func BenchmarkEcs(b *testing.B) {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	config.CollectionVersion = 1
	game.init(config)

	game.world.Startup()

	b.ResetTimer()

	var delta time.Duration
	_ = delta
	var ts time.Time
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.attack()
		game.world.Update()
		delta = time.Since(ts)
	}
}

func BenchmarkEcsSingleCore(b *testing.B) {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	config.CollectionVersion = 1
	game.init(config)

	game.world.Startup()

	runtime2.GOMAXPROCS(1)
	b.ResetTimer()

	var delta time.Duration
	_ = delta
	var ts time.Time
	for i := 0; i < b.N; i++ {
		ts = time.Now()
		game.attack()
		game.world.Update()
		delta = time.Since(ts)
	}
}
