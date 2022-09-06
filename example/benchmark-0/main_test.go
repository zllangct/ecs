package main

import (
	"fmt"
	"github.com/zllangct/ecs"
	_ "net/http/pprof"
	"reflect"
	"testing"
	"time"
	"unsafe"
	_ "unsafe"
)

func TestFrame(t *testing.T) {
	ecs.EnableMainThreadDebug()

	game := &GameECS{}
	game.init(ecs.NewDefaultWorldConfig())

	game.world.Startup()

	var delta time.Duration
	_ = delta
	var ts time.Time
	for i := 0; i < 10; i++ {
		ts = time.Now()
		game.world.Update()
		game.attack()
		delta = time.Since(ts)
		//ecs.Log.Info("===== Frame:", i, "=====", delta)
	}
}

func TestEcsOptimizer(t *testing.T) {
	ecs.EnableMainThreadDebug()

	game := &GameECS{}
	game.init(ecs.NewDefaultWorldConfig())

	game.world.Startup()

	var frameInterval = time.Millisecond * 33
	var delta time.Duration
	var ts time.Time
	for i := 0; i < 10; i++ {
		//ecs.Log.Info("===== Frame:", i)
		ts = time.Now()

		game.world.Update()
		game.attack()
		delta = time.Since(ts)
		ecs.Log.Info("===== Frame:", i, "=====", delta)
		if frameInterval-delta != 0 {
			game.world.Optimize(frameInterval-delta, true)
			time.Sleep(frameInterval - delta)
			delta = frameInterval
		}
	}
}

func TestOthers(t *testing.T) {
	var a1 = [...]reflect.Type{ecs.TypeOf[Test1](), ecs.TypeOf[Test2](), ecs.TypeOf[Test3]()}
	var a2 = [...]reflect.Type{ecs.TypeOf[Test1](), ecs.TypeOf[Test2](), ecs.TypeOf[Test3]()}
	m := map[interface{}]string{}
	m[a1] = "this is a1"
	m[a2] = "this is a2"
	fmt.Printf("%v\n", m)

	println(unsafe.Sizeof(ecs.Component[Test1]{}))
}
