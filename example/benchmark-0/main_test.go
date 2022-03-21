package main

import (
	"github.com/zllangct/ecs"
	"net/http"
	_ "net/http/pprof"
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

func BenchmarkEcsCollectionV1(b *testing.B) {
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

func BenchmarkEcsCollectionV2(b *testing.B) {
	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	config.CollectionVersion = 2
	game.init(config)

	b.ResetTimer()

	var delta time.Duration
	var ts time.Time
	for i := 0; i < b.N; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()
		game.attack()
		doFrame(game.world, uint64(i), delta)
		delta = time.Since(ts)
	}
}

func TestEcs(t *testing.T) {
	game := &GameECS{}
	config := ecs.NewDefaultWorldConfig()
	game.init(config)

	var delta time.Duration
	var ts time.Time
	for i := 0; i < 10; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()
		doFrame(game.world, uint64(i), delta)
		game.attack()
		delta = time.Since(ts)
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
