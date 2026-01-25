package main

import (
	"fmt"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/example1/components"
)

// MovementSystem - 根据速度更新位置
type MovementSystem struct{}

func (s *MovementSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("MovementSystem"),
		ecs.WithDep[components.Position](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
	)
	return nil
}

func (s *MovementSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := event.Delta.Seconds()
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 30.0 // 默认30fps
	}

	// 遍历所有有 Position 和 Velocity 组件的实体
	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Position](),
		ecs.WithComp[components.Velocity](),
	)
	for index, _ := range query.Iter() {
		pos, ok := ecs.GetBuddy[components.Position](ctx, index)
		if !ok {
			continue
		}
		vel, ok := ecs.GetBuddy[components.Velocity](ctx, index)
		if !ok {
			continue
		}

		// 根据速度和时间更新位置（直接访问公开字段）
		pos.X = pos.X + vel.X*float32(deltaSeconds)
		pos.Y = pos.Y + vel.Y*float32(deltaSeconds)
		pos.Z = pos.Z + vel.Z*float32(deltaSeconds)
	}

	return nil
}

// HealthSystem - 处理伤害事件
type HealthSystem struct{}

func (s *HealthSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("HealthSystem"),
		ecs.WithDep[components.Health](ecs.ReadWrite),
		ecs.WithDep[components.DamageEvent](ecs.ReadOnly),
	)
	return nil
}

func (s *HealthSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	// 遍历所有伤害事件
	for _, damage := range ecs.GetComponents[components.DamageEvent](ctx) {
		targetIndex := ecs.EntityIndex(damage.TargetId)

		// 获取目标的Health组件
		health, ok := ecs.GetBuddy[components.Health](ctx, targetIndex)
		if !ok {
			continue
		}

		// 应用伤害（直接访问公开字段）
		newHealth := health.Current - damage.Damage
		if newHealth < 0 {
			newHealth = 0
		}
		health.Current = newHealth

		fmt.Printf("Frame %d: Entity received %.1f damage, health: %.1f/%.1f\n",
			event.Frame, damage.Damage, health.Current, health.Max)
	}

	return nil
}

// RenderSystem - 打印实体状态（模拟渲染）
type RenderSystem struct {
	frameCount int
}

func (s *RenderSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("RenderSystem"),
		ecs.WithDep[components.Position](ecs.ReadOnly),
		ecs.WithDep[components.Player](ecs.ReadOnly),
	)
	return nil
}

func (s *RenderSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	s.frameCount++

	// 每10帧打印一次状态
	if s.frameCount%10 != 0 {
		return nil
	}

	fmt.Printf("\n=== Frame %d (Delta: %v) ===\n", event.Frame, event.Delta)

	// 遍历所有有 Position 和 Player 组件的实体
	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Position](),
		ecs.WithComp[components.Player](),
	)
	for index, _ := range query.Iter() {
		pos, ok := ecs.GetBuddy[components.Position](ctx, index)
		if !ok {
			continue
		}
		player, ok := ecs.GetBuddy[components.Player](ctx, index)
		if !ok {
			continue
		}

		// 从字节数组获取名字字符串
		name := string(player.Name[:])
		for i, b := range player.Name {
			if b == 0 {
				name = string(player.Name[:i])
				break
			}
		}

		fmt.Printf("  Player '%s' (Lv.%d): Position(%.2f, %.2f, %.2f) Health: %.1f/%.1f\n",
			name, player.Level,
			pos.X, pos.Y, pos.Z,
			player.Health, player.MaxHealth)
	}

	return nil
}

// LightRenderSystem - 轻量级系统示例（使用函数形式）
func LightRenderSystem(ctx *ecs.SystemContext, event ecs.Event) error {
	// 轻量级系统用于简单的打印任务
	if event.Frame%30 == 0 {
		fmt.Printf("[Light] Heartbeat at frame %d\n", event.Frame)
	}
	return nil
}
