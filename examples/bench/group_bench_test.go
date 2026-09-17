package bench

import (
	"testing"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/bench/components"
)

// ============================================================================
// 组表（Archetype Group）存储对照基准
// 对照组: 现状 CSet + optimizer 排序（WithoutGroups）
// 实验组: GroupMovementSystem 声明 WithGroup(Transform, Velocity)
// ============================================================================

// GroupMovementSystem 与 MovementSystem 逻辑一致，仅多声明固定访问组
type GroupMovementSystem struct{}

func (s *GroupMovementSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("GroupMovementSystem"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
		ecs.WithGroup(ecs.Dep[components.Transform](), ecs.Dep[components.Velocity]()),
	)
	return nil
}

func (s *GroupMovementSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	query := ctx.NewQuery(
		ecs.WithComp[components.Transform](),
		ecs.WithComp[components.Velocity](),
	)
	for index := range query.Iter() {
		transform, ok := ctx.GetBuddy[components.Transform](index)
		if !ok {
			continue
		}
		velocity, ok := ctx.GetBuddyReadOnly[components.VelocityReadOnly](index)
		if !ok {
			continue
		}

		transform.PosX += velocity.LinearX() * deltaSeconds
		transform.PosY += velocity.LinearY() * deltaSeconds
		transform.PosZ += velocity.LinearZ() * deltaSeconds
		transform.RotX += velocity.AngularX() * deltaSeconds
		transform.RotY += velocity.AngularY() * deltaSeconds
		transform.RotZ += velocity.AngularZ() * deltaSeconds
	}
	return nil
}

// View2ROMovementSystem 与 GroupMovementSystem 逻辑一致，改用 View2RO lockstep 迭代
type View2ROMovementSystem struct{}

func (s *View2ROMovementSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("View2ROMovementSystem"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
		ecs.WithGroup(ecs.Dep[components.Transform](), ecs.Dep[components.Velocity]()),
	)
	return nil
}

func (s *View2ROMovementSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}
	for t, v := range ecs.View2RO[components.Transform, components.VelocityReadOnly](ctx) {
		t.PosX += v.LinearX() * deltaSeconds
		t.PosY += v.LinearY() * deltaSeconds
		t.PosZ += v.LinearZ() * deltaSeconds
		t.RotX += v.AngularX() * deltaSeconds
		t.RotY += v.AngularY() * deltaSeconds
		t.RotZ += v.AngularZ() * deltaSeconds
	}
	return nil
}

func setupViewBenchWorld(entityCount int) *ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.Register[View2ROMovementSystem]()
	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{
			PosX: float32(i), PosY: 0, PosZ: 0,
			ScaleX: 1, ScaleY: 1, ScaleZ: 1,
		}
		velocity := &components.Velocity{
			LinearX: 1.0, LinearY: 0.5, LinearZ: 0.2,
			AngularX: 0.1, AngularY: 0.2, AngularZ: 0.3,
		}
		world.NewEntity(ecs.WithComponents(transform, velocity))
	}
	return world
}

func BenchmarkView2RO_Group_10000(b *testing.B) {
	world := setupViewBenchWorld(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkView2RO_Group_100000(b *testing.B) {
	world := setupViewBenchWorld(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func setupGroupBenchWorld(entityCount int, withGroup bool) (*ecs.World, []ecs.Entity) {
	opts := []ecs.WorldOption{ecs.WithWorldSyncMode()}
	if !withGroup {
		opts = append(opts, ecs.WithoutGroups())
	}
	world := ecs.NewWorld(opts...)
	world.Register[GroupMovementSystem]()

	entities := make([]ecs.Entity, 0, entityCount)
	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{
			PosX: float32(i), PosY: 0, PosZ: 0,
			ScaleX: 1, ScaleY: 1, ScaleZ: 1,
		}
		velocity := &components.Velocity{
			LinearX: 1.0, LinearY: 0.5, LinearZ: 0.2,
			AngularX: 0.1, AngularY: 0.2, AngularZ: 0.3,
		}
		entities = append(entities, world.NewEntity(ecs.WithComponents(transform, velocity)))
	}
	return world, entities
}

// churn 扰动：销毁并重建一半实体（隔一个删一个，再以新实体补回）。
// CSet 的密集数组经 swap-remove 后顺序与 EntityIndex 脱节，
// 两个组件池的 buddy 查找退化为双流随机访问；组表行保持紧凑，双列按行号同步顺序访问。
func churn(world *ecs.World, entities []ecs.Entity) []ecs.Entity {
	alive := make([]ecs.Entity, 0, len(entities)/2)
	for i, e := range entities {
		if i%2 == 0 {
			world.DestroyEntity(e)
		} else {
			alive = append(alive, e)
		}
	}
	world.Update() // flush 销毁
	// 补回新实体（新 EntityIndex，追加在尾部）
	for i := 0; i < len(entities)/2; i++ {
		transform := &components.Transform{PosX: float32(i), ScaleX: 1, ScaleY: 1, ScaleZ: 1}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5, LinearZ: 0.2}
		alive = append(alive, world.NewEntity(ecs.WithComponents(transform, velocity)))
	}
	world.Update() // flush 新增（组表在此收拢）
	return alive
}

func BenchmarkGroupStorage_CSet_10000(b *testing.B) {
	world, _ := setupGroupBenchWorld(10000, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkGroupStorage_Group_10000(b *testing.B) {
	world, _ := setupGroupBenchWorld(10000, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkGroupStorage_CSet_100000(b *testing.B) {
	world, _ := setupGroupBenchWorld(100000, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkGroupStorage_Group_100000(b *testing.B) {
	world, _ := setupGroupBenchWorld(100000, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkGroupStorage_CSet_Churned_100000(b *testing.B) {
	world, entities := setupGroupBenchWorld(100000, false)
	churn(world, entities)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkGroupStorage_Group_Churned_100000(b *testing.B) {
	world, entities := setupGroupBenchWorld(100000, true)
	churn(world, entities)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}
