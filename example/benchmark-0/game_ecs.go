package main

import (
	"github.com/zllangct/ecs"
	_ "unsafe"
)

type GameECS struct {
	world    *ecs.SyncWorld
	entities []ecs.Entity
}

func (g *GameECS) init(config *ecs.WorldConfig) {
	g.world = ecs.NewSyncWorld(config)

	ecs.RegisterSystem[MoveSystem](g.world)
	ecs.RegisterSystem[DamageSystem](g.world)
	ecs.RegisterSystem[Test1System](g.world)
	ecs.RegisterSystem[Test2System](g.world)
	ecs.RegisterSystem[Test3System](g.world)
	ecs.RegisterSystem[Test4System](g.world)
	ecs.RegisterSystem[Test5System](g.world)
	ecs.RegisterSystem[Test6System](g.world)
	ecs.RegisterSystem[Test7System](g.world)
	ecs.RegisterSystem[Test8System](g.world)
	ecs.RegisterSystem[Test9System](g.world)
	ecs.RegisterSystem[Test10System](g.world)

	DataGenerateECS(g)
}

func (g *GameECS) attack() {
	act := &Action{
		ActionType: 1,
	}
	for _, entity := range g.entities {
		info, _ := g.world.GetEntityInfo(entity)
		info.Add(g.world, act)
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
		t6 := &Test6{}
		t7 := &Test7{}
		t8 := &Test8{}
		t9 := &Test9{}
		t10 := &Test10{}

		e := game.world.NewEntity()
		game.world.Add(e, p, m, h, f, t1, t2, t3, t4, t5, t6, t7, t8, t9, t10)

		game.entities = append(game.entities, e)
	}
}
