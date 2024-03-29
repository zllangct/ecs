package main

import (
	"context"
	"github.com/zllangct/ecs"
	"net/http"
	_ "net/http/pprof"
	"test_ecs_fake_server/client"
	"test_ecs_fake_server/game"
	"test_ecs_fake_server/gm"
)

func main() {
	ecs.Log.Info("game start...")
	go func() {
		ecs.Log.Info(http.ListenAndServe("localhost:8889", nil))
	}()

	ctx := context.Background()

	//my game
	game := game.NewGame()
	//game manager
	gm := gm.NewGM()
	////client manager
	cm := client.NewClient()

	go game.Run(ctx)
	go gm.Run(ctx, game)
	go cm.Run(ctx)

	<-ctx.Done()
	ecs.Log.Info("game end...")
}
