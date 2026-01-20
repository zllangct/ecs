package bench

import (
	"testing"

	"github.com/zllangct/ecs"
	"github.com/zllangct/ecs/example/bench/components"
)

// ============================================================================
// 测试维度说明：
// 1. 实体数量维度：100, 1000, 10000, 100000 实体
// 2. 组件数量维度：单组件、双组件、多组件（8组件）
// 3. 访问模式维度：顺序访问、随机访问
// 4. 系统类型维度：简单系统（位置更新）、复杂系统（物理+AI+战斗）
// 5. 内存布局维度：ECS紧凑布局 vs OOP指针布局 vs OOP内联布局
// ============================================================================

// ============================================================================
// 全局变量，防止编译器优化掉计算结果
// ============================================================================
var sinkFloat32 float32

// ============================================================================
// ECS System 定义
// ============================================================================

// MovementSystem - 简单移动系统
type MovementSystem struct{}

func (s *MovementSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("MovementSystem"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
	)
	return nil
}

func (s *MovementSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Transform](),
		ecs.WithComp[components.Velocity](),
	)
	for index := range query.Iter() {
		transform, ok := ecs.GetBuddy[components.Transform](ctx, index)
		if !ok {
			continue
		}
		velocity, ok := ecs.GetBuddy[components.Velocity](ctx, index)
		if !ok {
			continue
		}

		transform.PosX += velocity.LinearX * deltaSeconds
		transform.PosY += velocity.LinearY * deltaSeconds
		transform.PosZ += velocity.LinearZ * deltaSeconds
		transform.RotX += velocity.AngularX * deltaSeconds
		transform.RotY += velocity.AngularY * deltaSeconds
		transform.RotZ += velocity.AngularZ * deltaSeconds
	}
	return nil
}

// PhysicsSystem - 物理系统
type PhysicsSystem struct{}

func (s *PhysicsSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("PhysicsSystem"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadWrite),
		ecs.WithDep[components.Physics](ecs.ReadOnly),
	)
	return nil
}

func (s *PhysicsSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Transform](),
		ecs.WithComp[components.Velocity](),
		ecs.WithComp[components.Physics](),
	)
	for index := range query.Iter() {
		transform, ok := ecs.GetBuddy[components.Transform](ctx, index)
		if !ok {
			continue
		}
		velocity, ok := ecs.GetBuddy[components.Velocity](ctx, index)
		if !ok {
			continue
		}
		physics, ok := ecs.GetBuddy[components.Physics](ctx, index)
		if !ok {
			continue
		}

		// 重力
		if physics.UseGravity == 1 {
			velocity.LinearY -= 9.8 * deltaSeconds
		}
		// 阻力
		velocity.LinearX *= (1 - physics.Drag*deltaSeconds)
		velocity.LinearY *= (1 - physics.Drag*deltaSeconds)
		velocity.LinearZ *= (1 - physics.Drag*deltaSeconds)
		// 更新位置
		transform.PosX += velocity.LinearX * deltaSeconds
		transform.PosY += velocity.LinearY * deltaSeconds
		transform.PosZ += velocity.LinearZ * deltaSeconds
	}
	return nil
}

// HealthRegenSystem - 生命恢复系统
type HealthRegenSystem struct{}

func (s *HealthRegenSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("HealthRegenSystem"),
		ecs.WithDep[components.Health](ecs.ReadWrite),
	)
	return nil
}

func (s *HealthRegenSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	for _, health := range ecs.GetComponents[components.Health](ctx) {
		if health.Current < health.Max {
			health.Current += health.Regen * deltaSeconds
			if health.Current > health.Max {
				health.Current = health.Max
			}
		}
	}
	return nil
}

// ComplexActorSystem - 复杂Actor更新系统（包含移动、物理、生命）
type ComplexActorSystem struct{}

func (s *ComplexActorSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("ComplexActorSystem"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadWrite),
		ecs.WithDep[components.Physics](ecs.ReadOnly),
		ecs.WithDep[components.Health](ecs.ReadWrite),
	)
	return nil
}

func (s *ComplexActorSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Transform](),
		ecs.WithComp[components.Velocity](),
		ecs.WithComp[components.Physics](),
		ecs.WithComp[components.Health](),
	)
	for index := range query.Iter() {
		transform, _ := ecs.GetBuddy[components.Transform](ctx, index)
		velocity, _ := ecs.GetBuddy[components.Velocity](ctx, index)
		physics, _ := ecs.GetBuddy[components.Physics](ctx, index)
		health, _ := ecs.GetBuddy[components.Health](ctx, index)

		// 移动
		transform.PosX += velocity.LinearX * deltaSeconds
		transform.PosY += velocity.LinearY * deltaSeconds
		transform.PosZ += velocity.LinearZ * deltaSeconds

		// 物理
		if physics.UseGravity == 1 {
			velocity.LinearY -= 9.8 * deltaSeconds
		}

		// 生命恢复
		if health.Current < health.Max {
			health.Current += health.Regen * deltaSeconds
			if health.Current > health.Max {
				health.Current = health.Max
			}
		}
	}
	return nil
}

// ============================================================================
// 辅助函数
// ============================================================================

// 创建ECS World 并注册系统
func setupECSWorldSimple(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&MovementSystem{})

	// 创建实体
	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{
			PosX: float32(i), PosY: 0, PosZ: 0,
			RotX: 0, RotY: 0, RotZ: 0,
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

func setupECSWorldPhysics(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&PhysicsSystem{})

	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 10.0, LinearZ: 0.2}
		physics := &components.Physics{Mass: 1.0, Drag: 0.1, UseGravity: 1}
		world.NewEntity(ecs.WithComponents(transform, velocity, physics))
	}
	return world
}

func setupECSWorldComplex(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&ComplexActorSystem{})

	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 10.0}
		physics := &components.Physics{Mass: 1.0, Drag: 0.1, UseGravity: 1}
		health := &components.Health{Current: 50, Max: 100, Regen: 1.0}
		world.NewEntity(ecs.WithComponents(transform, velocity, physics, health))
	}
	return world
}

func setupECSWorldFull(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&MovementSystem{})
	world.RegisterStandard(&HealthRegenSystem{})

	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5}
		physics := &components.Physics{Mass: 1.0, Drag: 0.1, UseGravity: 1}
		render := &components.Render{MeshId: 1, MaterialId: 1, Visible: 1}
		health := &components.Health{Current: 50, Max: 100, Regen: 1.0}
		ai := &components.AI{State: 0, AggroRange: 10, AttackRange: 2}
		movement := &components.Movement{Speed: 5, MaxSpeed: 10, Acceleration: 2}
		combat := &components.Combat{Attack: 10, Defense: 5, CritRate: 0.1}
		world.NewEntity(ecs.WithComponents(
			transform, velocity, physics, render,
			health, ai, movement, combat,
		))
	}
	return world
}

// 创建OOP World
func setupOOPWorldSimple(entityCount int) *OOPWorld {
	world := NewOOPWorld()
	for i := 0; i < entityCount; i++ {
		e := world.NewOOPMoveable(int64(i))
		e.Transform.PosX = float32(i)
		e.Velocity.LinearX = 1.0
		e.Velocity.LinearY = 0.5
		e.Velocity.LinearZ = 0.2
		e.Velocity.AngularX = 0.1
		e.Velocity.AngularY = 0.2
		e.Velocity.AngularZ = 0.3
	}
	return world
}

func setupOOPWorldPhysics(entityCount int) *OOPWorld {
	world := NewOOPWorld()
	for i := 0; i < entityCount; i++ {
		e := world.NewOOPPhysical(int64(i))
		e.Transform.PosX = float32(i)
		e.Velocity.LinearX = 1.0
		e.Velocity.LinearY = 10.0
		e.Velocity.LinearZ = 0.2
		e.Physics.Mass = 1.0
		e.Physics.Drag = 0.1
		e.Physics.UseGravity = 1
	}
	return world
}

func setupOOPWorldComplex(entityCount int) *OOPWorld {
	world := NewOOPWorld()
	for i := 0; i < entityCount; i++ {
		e := world.NewOOPActor(int64(i))
		e.Transform.PosX = float32(i)
		e.Velocity.LinearX = 1.0
		e.Velocity.LinearY = 10.0
		e.Physics.Mass = 1.0
		e.Physics.Drag = 0.1
		e.Physics.UseGravity = 1
		e.Health.Current = 50
		e.Health.Max = 100
		e.Health.Regen = 1.0
	}
	return world
}

func setupOOPWorldGeneric(entityCount int) *OOPWorld {
	world := NewOOPWorld()
	for i := 0; i < entityCount; i++ {
		e := world.NewOOPEntity(int64(i))
		e.Transform = &OOPTransform{PosX: float32(i)}
		e.Velocity = &OOPVelocity{LinearX: 1.0, LinearY: 0.5}
		e.Physics = &OOPPhysics{Mass: 1.0, Drag: 0.1, UseGravity: 1}
		e.Health = &OOPHealth{Current: 50, Max: 100, Regen: 1.0}
	}
	return world
}

// ============================================================================
// 维度1: 实体数量对比 - 简单移动系统
// ============================================================================

func BenchmarkEntityCount_ECS_100(b *testing.B) {
	world := setupECSWorldSimple(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkEntityCount_OOP_100(b *testing.B) {
	world := setupOOPWorldSimple(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

func BenchmarkEntityCount_ECS_1000(b *testing.B) {
	world := setupECSWorldSimple(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkEntityCount_OOP_1000(b *testing.B) {
	world := setupOOPWorldSimple(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

func BenchmarkEntityCount_ECS_10000(b *testing.B) {
	world := setupECSWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkEntityCount_OOP_10000(b *testing.B) {
	world := setupOOPWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

func BenchmarkEntityCount_ECS_100000(b *testing.B) {
	world := setupECSWorldSimple(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkEntityCount_OOP_100000(b *testing.B) {
	world := setupOOPWorldSimple(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

// ============================================================================
// 维度2: 组件数量对比 - 双组件(简单) vs 三组件(物理) vs 四组件(复杂)
// ============================================================================

func BenchmarkComponentCount_ECS_2Comp_10000(b *testing.B) {
	world := setupECSWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkComponentCount_OOP_2Comp_10000(b *testing.B) {
	world := setupOOPWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

func BenchmarkComponentCount_ECS_3Comp_10000(b *testing.B) {
	world := setupECSWorldPhysics(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkComponentCount_OOP_3Comp_10000(b *testing.B) {
	world := setupOOPWorldPhysics(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdatePhysicalEntities(1.0 / 60.0)
	}
}

func BenchmarkComponentCount_ECS_4Comp_10000(b *testing.B) {
	world := setupECSWorldComplex(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkComponentCount_OOP_4Comp_10000(b *testing.B) {
	world := setupOOPWorldComplex(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateActorEntities(1.0 / 60.0)
	}
}

// ============================================================================
// 维度3: 内存布局对比 - ECS紧凑 vs OOP内联 vs OOP指针
// ============================================================================

func BenchmarkMemoryLayout_ECS_10000(b *testing.B) {
	world := setupECSWorldComplex(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkMemoryLayout_OOP_Inline_10000(b *testing.B) {
	world := setupOOPWorldComplex(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateActorEntities(1.0 / 60.0)
	}
}

func BenchmarkMemoryLayout_OOP_Pointer_10000(b *testing.B) {
	world := setupOOPWorldGeneric(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateGenericEntities(1.0 / 60.0)
	}
}

// ============================================================================
// 维度4: 大规模实体测试
// ============================================================================

func BenchmarkLargeScale_ECS_50000(b *testing.B) {
	world := setupECSWorldSimple(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkLargeScale_OOP_50000(b *testing.B) {
	world := setupOOPWorldSimple(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

// ============================================================================
// 维度5: 实体创建性能对比
// ============================================================================

func BenchmarkEntityCreation_ECS_1000(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world := ecs.NewWorld(ecs.WithWorldSyncMode())
		for j := 0; j < 1000; j++ {
			transform := &components.Transform{PosX: float32(j)}
			velocity := &components.Velocity{LinearX: 1.0}
			world.NewEntity(ecs.WithComponents(transform, velocity))
		}
	}
}

func BenchmarkEntityCreation_OOP_1000(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world := NewOOPWorld()
		for j := 0; j < 1000; j++ {
			e := world.NewOOPMoveable(int64(j))
			e.Transform.PosX = float32(j)
			e.Velocity.LinearX = 1.0
		}
	}
}

// ============================================================================
// 维度6: 组件访问性能 - 纯遍历不做复杂计算
// ============================================================================

func BenchmarkPureIteration_ECS_10000(b *testing.B) {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())

	// 使用轻量系统进行纯遍历
	world.RegisterLight(func(ctx *ecs.SystemContext, event ecs.Event) error {
		var sum float32
		query := ecs.NewQuery(ctx,
			ecs.WithComp[components.Transform](),
		)
		for index := range query.Iter() {
			transform, ok := ecs.GetBuddy[components.Transform](ctx, index)
			if ok {
				sum += transform.PosX
			}
		}
		sinkFloat32 = sum
		return nil
	}, ecs.WithName("PureIterationSystem"),
		ecs.WithDep[components.Transform](ecs.ReadOnly))

	for i := 0; i < 10000; i++ {
		transform := &components.Transform{PosX: float32(i)}
		world.NewEntity(ecs.WithComponents(transform))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkPureIteration_OOP_10000(b *testing.B) {
	type SimpleEntity struct {
		PosX float32
	}
	entities := make([]SimpleEntity, 10000)
	for i := 0; i < 10000; i++ {
		entities[i].PosX = float32(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for j := range entities {
			sum += entities[j].PosX
		}
		sinkFloat32 = sum
	}
}

func BenchmarkPureIteration_OOP_Pointer_10000(b *testing.B) {
	type SimpleEntity struct {
		PosX float32
	}
	entities := make([]*SimpleEntity, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = &SimpleEntity{PosX: float32(i)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for _, e := range entities {
			sum += e.PosX
		}
		sinkFloat32 = sum
	}
}

// ============================================================================
// 维度7: 完整系统更新对比（8组件实体）
// ============================================================================

func BenchmarkFullEntity_ECS_10000(b *testing.B) {
	world := setupECSWorldFull(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkFullEntity_OOP_10000(b *testing.B) {
	world := NewOOPWorld()
	for i := 0; i < 10000; i++ {
		e := world.NewOOPActor(int64(i))
		e.Transform.PosX = float32(i)
		e.Velocity.LinearX = 1.0
		e.Health.Current = 50
		e.Health.Max = 100
		e.Health.Regen = 1.0
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 移动更新
		for _, e := range world.ActorEntities {
			e.Transform.PosX += e.Velocity.LinearX * (1.0 / 60.0)
			e.Transform.PosY += e.Velocity.LinearY * (1.0 / 60.0)
			e.Transform.PosZ += e.Velocity.LinearZ * (1.0 / 60.0)
		}
		// 生命恢复
		for _, e := range world.ActorEntities {
			if e.Health.Current < e.Health.Max {
				e.Health.Current += e.Health.Regen * (1.0 / 60.0)
			}
		}
	}
}

// ============================================================================
// 维度8: 稀疏组件查询（只有部分实体有目标组件）
// ============================================================================

func BenchmarkSparseQuery_ECS_10000(b *testing.B) {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&HealthRegenSystem{})

	// 只有10%的实体有Health组件
	for i := 0; i < 10000; i++ {
		transform := &components.Transform{PosX: float32(i)}
		if i%10 == 0 {
			health := &components.Health{Current: 50, Max: 100, Regen: 1.0}
			world.NewEntity(ecs.WithComponents(transform, health))
		} else {
			world.NewEntity(ecs.WithComponents(transform))
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkSparseQuery_OOP_10000(b *testing.B) {
	world := NewOOPWorld()

	// 只有10%的实体有Health组件
	for i := 0; i < 10000; i++ {
		e := world.NewOOPEntity(int64(i))
		e.Transform = &OOPTransform{PosX: float32(i)}
		if i%10 == 0 {
			e.Health = &OOPHealth{Current: 50, Max: 100, Regen: 1.0}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 需要检查组件是否存在
		for _, e := range world.Entities {
			if e.Health != nil {
				if e.Health.Current < e.Health.Max {
					e.Health.Current += e.Health.Regen * (1.0 / 60.0)
				}
			}
		}
	}
}

// ============================================================================
// 维度9: 排除ECS框架开销 - 直接遍历组件存储（CSet）
// 对比ECS系统更新与直接遍历CSet存储的性能差异
// ============================================================================

func BenchmarkNoFrameworkOverhead_ECS_System_10000(b *testing.B) {
	// 使用完整的ECS系统（包含框架开销）
	world := setupECSWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkNoFrameworkOverhead_ECS_DirectCSet_10000(b *testing.B) {
	// 直接使用CSet存储，排除框架开销
	cset := ecs.NewCSet[components.Transform]()
	velocityCset := ecs.NewCSet[components.Velocity]()

	for i := 0; i < 10000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5, LinearZ: 0.2}
		cset.Add(entity, transform)
		velocityCset.Add(entity, velocity)
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 直接遍历CSet，无框架开销
		for idx, transform := range cset.Iter() {
			velocity := velocityCset.Get(ecs.Entity(idx))
			if velocity == nil {
				continue
			}
			vel := velocity.(*components.Velocity)
			transform.PosX += vel.LinearX * deltaSeconds
			transform.PosY += vel.LinearY * deltaSeconds
			transform.PosZ += vel.LinearZ * deltaSeconds
		}
	}
}

func BenchmarkNoFrameworkOverhead_OOP_Slice_10000(b *testing.B) {
	// OOP方式使用切片存储
	world := setupOOPWorldSimple(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

func BenchmarkNoFrameworkOverhead_ECS_DirectCSet_100000(b *testing.B) {
	// 大规模测试：直接使用CSet存储
	cset := ecs.NewCSet[components.Transform]()
	velocityCset := ecs.NewCSet[components.Velocity]()

	for i := 0; i < 100000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5}
		cset.Add(entity, transform)
		velocityCset.Add(entity, velocity)
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for idx, transform := range cset.Iter() {
			velocity := velocityCset.Get(ecs.Entity(idx))
			if velocity == nil {
				continue
			}
			vel := velocity.(*components.Velocity)
			transform.PosX += vel.LinearX * deltaSeconds
			transform.PosY += vel.LinearY * deltaSeconds
			transform.PosZ += vel.LinearZ * deltaSeconds
		}
	}
}

func BenchmarkNoFrameworkOverhead_OOP_Slice_100000(b *testing.B) {
	world := setupOOPWorldSimple(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.UpdateMoveableEntities(1.0 / 60.0)
	}
}

// ============================================================================
// 维度10: 并行模式对比 - WithWorldASyncMode
// 测试异步/并行执行模式的性能
// ============================================================================

// ParallelMovementSystem - 并行移动系统1
type ParallelMovementSystem1 struct{}

func (s *ParallelMovementSystem1) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("ParallelMovementSystem1"),
		ecs.WithDep[components.Transform](ecs.ReadWrite),
		ecs.WithDep[components.Velocity](ecs.ReadOnly),
	)
	return nil
}

func (s *ParallelMovementSystem1) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	query := ecs.NewQuery(ctx,
		ecs.WithComp[components.Transform](),
		ecs.WithComp[components.Velocity](),
	)
	for index := range query.Iter() {
		transform, ok := ecs.GetBuddy[components.Transform](ctx, index)
		if !ok {
			continue
		}
		velocity, ok := ecs.GetBuddy[components.Velocity](ctx, index)
		if !ok {
			continue
		}
		transform.PosX += velocity.LinearX * deltaSeconds
		transform.PosY += velocity.LinearY * deltaSeconds
		transform.PosZ += velocity.LinearZ * deltaSeconds
	}
	return nil
}

// ParallelHealthSystem - 并行生命系统（与移动系统无依赖）
type ParallelHealthSystem struct{}

func (s *ParallelHealthSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("ParallelHealthSystem"),
		ecs.WithDep[components.Health](ecs.ReadWrite),
	)
	return nil
}

func (s *ParallelHealthSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	for _, health := range ecs.GetComponents[components.Health](ctx) {
		if health.Current < health.Max {
			health.Current += health.Regen * deltaSeconds
			if health.Current > health.Max {
				health.Current = health.Max
			}
		}
	}
	return nil
}

// ParallelCombatSystem - 并行战斗系统（与其他系统无依赖）
type ParallelCombatSystem struct{}

func (s *ParallelCombatSystem) Init(ctx *ecs.SystemInitContext) error {
	ctx.SetOption(
		ecs.WithName("ParallelCombatSystem"),
		ecs.WithDep[components.Combat](ecs.ReadWrite),
	)
	return nil
}

func (s *ParallelCombatSystem) Update(ctx *ecs.SystemContext, event ecs.Event) error {
	deltaSeconds := float32(event.Delta.Seconds())
	if deltaSeconds == 0 {
		deltaSeconds = 1.0 / 60.0
	}

	for _, combat := range ecs.GetComponents[components.Combat](ctx) {
		// 模拟攻击速度影响的计算
		combat.Attack *= (1 + combat.AttackSpeed*deltaSeconds*0.001)
	}
	return nil
}

func setupECSWorldSyncMultiSystem(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldSyncMode())
	world.RegisterStandard(&ParallelMovementSystem1{})
	world.RegisterStandard(&ParallelHealthSystem{})
	world.RegisterStandard(&ParallelCombatSystem{})

	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5}
		health := &components.Health{Current: 50, Max: 100, Regen: 1.0}
		combat := &components.Combat{Attack: 10, AttackSpeed: 1.0}
		world.NewEntity(ecs.WithComponents(transform, velocity, health, combat))
	}
	return world
}

func setupECSWorldAsyncMultiSystem(entityCount int) ecs.World {
	world := ecs.NewWorld(ecs.WithWorldASyncMode())
	world.RegisterStandard(&ParallelMovementSystem1{})
	world.RegisterStandard(&ParallelHealthSystem{})
	world.RegisterStandard(&ParallelCombatSystem{})

	for i := 0; i < entityCount; i++ {
		transform := &components.Transform{PosX: float32(i)}
		velocity := &components.Velocity{LinearX: 1.0, LinearY: 0.5}
		health := &components.Health{Current: 50, Max: 100, Regen: 1.0}
		combat := &components.Combat{Attack: 10, AttackSpeed: 1.0}
		world.NewEntity(ecs.WithComponents(transform, velocity, health, combat))
	}
	return world
}

func BenchmarkParallelMode_Sync_10000(b *testing.B) {
	world := setupECSWorldSyncMultiSystem(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkParallelMode_Async_10000(b *testing.B) {
	world := setupECSWorldAsyncMultiSystem(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkParallelMode_Sync_50000(b *testing.B) {
	world := setupECSWorldSyncMultiSystem(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkParallelMode_Async_50000(b *testing.B) {
	world := setupECSWorldAsyncMultiSystem(50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkParallelMode_Sync_100000(b *testing.B) {
	world := setupECSWorldSyncMultiSystem(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

func BenchmarkParallelMode_Async_100000(b *testing.B) {
	world := setupECSWorldAsyncMultiSystem(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update()
	}
}

// ============================================================================
// 维度11: CSet存储 vs OOP存储 - 纯遍历性能对比
// 排除所有框架开销，单纯对比数据存储结构的遍历性能
// ============================================================================

func BenchmarkStorageIteration_CSet_10000(b *testing.B) {
	// 使用CSet存储组件（连续内存布局）
	cset := ecs.NewCSet[components.Transform]()

	for i := 0; i < 10000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
		cset.Add(entity, transform)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for _, transform := range cset.Iter() {
			sum += transform.PosX + transform.PosY + transform.PosZ
		}
		sinkFloat32 = sum
	}
}

func BenchmarkStorageIteration_OOP_Slice_10000(b *testing.B) {
	// OOP方式：切片存储结构体（连续内存）
	entities := make([]OOPTransform, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = OOPTransform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for j := range entities {
			sum += entities[j].PosX + entities[j].PosY + entities[j].PosZ
		}
		sinkFloat32 = sum
	}
}

func BenchmarkStorageIteration_OOP_PointerSlice_10000(b *testing.B) {
	// OOP方式：切片存储指针（非连续内存）
	entities := make([]*OOPTransform, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = &OOPTransform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for _, e := range entities {
			sum += e.PosX + e.PosY + e.PosZ
		}
		sinkFloat32 = sum
	}
}

func BenchmarkStorageIteration_CSet_100000(b *testing.B) {
	cset := ecs.NewCSet[components.Transform]()

	for i := 0; i < 100000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
		cset.Add(entity, transform)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for _, transform := range cset.Iter() {
			sum += transform.PosX + transform.PosY + transform.PosZ
		}
		sinkFloat32 = sum
	}
}

func BenchmarkStorageIteration_OOP_Slice_100000(b *testing.B) {
	entities := make([]OOPTransform, 100000)
	for i := 0; i < 100000; i++ {
		entities[i] = OOPTransform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for j := range entities {
			sum += entities[j].PosX + entities[j].PosY + entities[j].PosZ
		}
		sinkFloat32 = sum
	}
}

func BenchmarkStorageIteration_OOP_PointerSlice_100000(b *testing.B) {
	entities := make([]*OOPTransform, 100000)
	for i := 0; i < 100000; i++ {
		entities[i] = &OOPTransform{PosX: float32(i), PosY: float32(i * 2), PosZ: float32(i * 3)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sum float32
		for _, e := range entities {
			sum += e.PosX + e.PosY + e.PosZ
		}
		sinkFloat32 = sum
	}
}

// ============================================================================
// 维度12: CSet存储 vs OOP存储 - 带修改操作的遍历对比
// ============================================================================

func BenchmarkStorageModify_CSet_10000(b *testing.B) {
	cset := ecs.NewCSet[components.Transform]()

	for i := 0; i < 10000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i)}
		cset.Add(entity, transform)
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, transform := range cset.Iter() {
			transform.PosX += 1.0 * deltaSeconds
			transform.PosY += 0.5 * deltaSeconds
			transform.PosZ += 0.2 * deltaSeconds
		}
	}
}

func BenchmarkStorageModify_OOP_Slice_10000(b *testing.B) {
	entities := make([]OOPTransform, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = OOPTransform{PosX: float32(i)}
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range entities {
			entities[j].PosX += 1.0 * deltaSeconds
			entities[j].PosY += 0.5 * deltaSeconds
			entities[j].PosZ += 0.2 * deltaSeconds
		}
	}
}

func BenchmarkStorageModify_OOP_PointerSlice_10000(b *testing.B) {
	entities := make([]*OOPTransform, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = &OOPTransform{PosX: float32(i)}
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.PosX += 1.0 * deltaSeconds
			e.PosY += 0.5 * deltaSeconds
			e.PosZ += 0.2 * deltaSeconds
		}
	}
}

func BenchmarkStorageModify_CSet_100000(b *testing.B) {
	cset := ecs.NewCSet[components.Transform]()

	for i := 0; i < 100000; i++ {
		entity := ecs.Entity(i)
		transform := &components.Transform{PosX: float32(i)}
		cset.Add(entity, transform)
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, transform := range cset.Iter() {
			transform.PosX += 1.0 * deltaSeconds
			transform.PosY += 0.5 * deltaSeconds
			transform.PosZ += 0.2 * deltaSeconds
		}
	}
}

func BenchmarkStorageModify_OOP_Slice_100000(b *testing.B) {
	entities := make([]OOPTransform, 100000)
	for i := 0; i < 100000; i++ {
		entities[i] = OOPTransform{PosX: float32(i)}
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := range entities {
			entities[j].PosX += 1.0 * deltaSeconds
			entities[j].PosY += 0.5 * deltaSeconds
			entities[j].PosZ += 0.2 * deltaSeconds
		}
	}
}

func BenchmarkStorageModify_OOP_PointerSlice_100000(b *testing.B) {
	entities := make([]*OOPTransform, 100000)
	for i := 0; i < 100000; i++ {
		entities[i] = &OOPTransform{PosX: float32(i)}
	}

	deltaSeconds := float32(1.0 / 60.0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, e := range entities {
			e.PosX += 1.0 * deltaSeconds
			e.PosY += 0.5 * deltaSeconds
			e.PosZ += 0.2 * deltaSeconds
		}
	}
}
