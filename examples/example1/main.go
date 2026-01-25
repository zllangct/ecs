package main

import (
	"fmt"
	"time"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/example1/components"
)

func main() {
	fmt.Println("=== ECS Example 1: 完整的ECS流程 ===")
	fmt.Println()

	// 注意：组件已经在 components 包的 init() 函数中自动注册
	// 无需手动调用 ecs.RegisterComponent

	// 1. 创建World
	world := ecs.NewWorld(
		ecs.WithWorldSyncMode(), // 同步模式，便于调试
		ecs.WithWorldDefaultUpdateRate(30),
	)

	// 2. 注册Systems
	// 标准System（实现Init和Update接口）
	if err := world.RegisterStandard(&MovementSystem{}); err != nil {
		panic(err)
	}
	if err := world.RegisterStandard(&HealthSystem{}); err != nil {
		panic(err)
	}
	if err := world.RegisterStandard(&RenderSystem{}); err != nil {
		panic(err)
	}

	// 轻量级System（函数形式）
	if err := world.RegisterLight(LightRenderSystem,
		ecs.WithName("LightRenderSystem"),
	); err != nil {
		panic(err)
	}

	// 4. 创建实体
	// 创建玩家实体
	player1Pos := &components.Position{}
	player1Pos.X = 0
	player1Pos.Y = 0
	player1Pos.Z = 0

	player1Vel := &components.Velocity{}
	player1Vel.X = 1.0
	player1Vel.Y = 0.5
	player1Vel.Z = 0

	player1Info := &components.Player{}
	copy(player1Info.Name[:], "Hero")
	player1Info.Level = 10
	player1Info.Health = 100
	player1Info.MaxHealth = 100

	player1Health := &components.Health{}
	player1Health.Current = 100
	player1Health.Max = 100

	player1 := world.NewEntity(ecs.WithComponents(
		player1Pos,
		player1Vel,
		player1Info,
		player1Health,
	))

	// 创建第二个玩家实体
	player2Pos := &components.Position{}
	player2Pos.X = 10
	player2Pos.Y = 5
	player2Pos.Z = 0

	player2Vel := &components.Velocity{}
	player2Vel.X = -0.5
	player2Vel.Y = 1.0
	player2Vel.Z = 0

	player2Info := &components.Player{}
	copy(player2Info.Name[:], "Warrior")
	player2Info.Level = 15
	player2Info.Health = 150
	player2Info.MaxHealth = 150

	player2Health := &components.Health{}
	player2Health.Current = 150
	player2Health.Max = 150

	player2 := world.NewEntity(ecs.WithComponents(
		player2Pos,
		player2Vel,
		player2Info,
		player2Health,
	))

	// 创建一个只有位置和速度的NPC（非玩家实体）
	npcPos := &components.Position{}
	npcPos.X = 5
	npcPos.Y = 5
	npcPos.Z = 0

	npcVel := &components.Velocity{}
	npcVel.X = 0.2
	npcVel.Y = -0.2
	npcVel.Z = 0

	world.NewEntity(ecs.WithComponents(npcPos, npcVel))

	fmt.Printf("Created entities: Player1(Entity: %d), Player2(Entity: %d)\n",
		player1.Entity().ToInt64(), player2.Entity().ToInt64())
	fmt.Println()

	// 5. 运行World主循环
	fmt.Println("Starting game loop (100 frames)...")
	fmt.Println()

	for i := 0; i < 100; i++ {
		// 在第50帧添加伤害事件
		if i == 50 {
			damageEvent := &components.DamageEvent{}
			damageEvent.Damage = 25
			damageEvent.SourceId = player2.Entity().ToInt64()
			damageEvent.TargetId = int64(player1.Entity().Index())

			world.NewEntity(ecs.WithComponents(damageEvent))
			fmt.Println("\n>>> Added damage event at frame 50")
		}

		// 更新World
		if err := world.Update(); err != nil {
			panic(err)
		}

		// 控制帧率
		time.Sleep(time.Millisecond * 10)
	}

	fmt.Println()
	fmt.Println("=== Game loop finished ===")

	// 6. 销毁World
	if err := world.Destroy(); err != nil {
		panic(err)
	}
}
