package main

import (
	"github.com/zllangct/ecs"
	"time"
	_ "unsafe"
)

type GameECS struct {
	world    ecs.IWorld
	entities []ecs.Entity
}

//go:linkname doFrame github.com/zllangct/ecs.doFrameForBenchmark
func doFrame(w ecs.IWorld, frame uint64, lastDelta time.Duration)

func (g *GameECS) init(config *ecs.WorldConfig) {
	ecs.Configure(ecs.NewDefaultRuntimeConfig())
	ecs.Run()

	g.world = ecs.CreateWorld(config)

	ecs.RegisterSystem[MoveSystem](g.world)
	ecs.RegisterSystem[DamageSystem](g.world)
	ecs.RegisterSystem[Test1System](g.world)
	ecs.RegisterSystem[Test2System](g.world)
	ecs.RegisterSystem[Test3System](g.world)
	ecs.RegisterSystem[Test4System](g.world)
	ecs.RegisterSystem[Test5System](g.world)

	DataGenerateECS(g)
}

func (g *GameECS) attack() {
	act := &Action{
		ActionType: 1,
	}
	for _, entity := range g.entities {
		info := ecs.GetEntityInfo(g.world, entity)
		err := info.Add(act)
		if err != nil {
			ecs.Log.Infof("%+v", err)
		}
	}
}

func DataGenerateECS(game *GameECS) {
	for i := 0; i < PlayerCount; i++ {
		p := &Position{
			X: 0,
			Y: 0,
			Z: 0,
		}
		m := &Movement{
			V:   100,
			Dir: [3]int{1, 0, 0},
		}
		h := &HealthPoint{
			HP: 100,
		}
		f := &Force{
			AttackRange:        10000,
			PhysicalBaseAttack: 10,
		}

		t1 := &Test1{}
		t2 := &Test2{}
		t3 := &Test3{}
		t4 := &Test4{}
		t5 := &Test5{}

		e := game.world.NewEntity()
		e.Add(p, m, h, f, t1, t2, t3, t4, t5)
		game.entities = append(game.entities, e.Entity())
	}
}
