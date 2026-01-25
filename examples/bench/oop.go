package bench

// OOP方式的实体和组件定义，用于与ECS对比

// OOPTransform 变换组件（OOP版本）
type OOPTransform struct {
	PosX   float32
	PosY   float32
	PosZ   float32
	RotX   float32
	RotY   float32
	RotZ   float32
	ScaleX float32
	ScaleY float32
	ScaleZ float32
}

// OOPVelocity 速度组件（OOP版本）
type OOPVelocity struct {
	LinearX  float32
	LinearY  float32
	LinearZ  float32
	AngularX float32
	AngularY float32
	AngularZ float32
}

// OOPPhysics 物理组件（OOP版本）
type OOPPhysics struct {
	Mass        float32
	Drag        float32
	AngularDrag float32
	UseGravity  int8
	IsKinematic int8
}

// OOPRender 渲染组件（OOP版本）
type OOPRender struct {
	MeshId        int32
	MaterialId    int32
	Visible       int8
	CastShadow    int8
	ReceiveShadow int8
}

// OOPHealth 生命组件（OOP版本）
type OOPHealth struct {
	Current float32
	Max     float32
	Regen   float32
}

// OOPAI AI组件（OOP版本）
type OOPAI struct {
	State       int32
	TargetId    int64
	AggroRange  float32
	AttackRange float32
}

// OOPMovement 移动组件（OOP版本）
type OOPMovement struct {
	Speed        float32
	MaxSpeed     float32
	Acceleration float32
	Deceleration float32
	TurnSpeed    float32
}

// OOPCombat 战斗组件（OOP版本）
type OOPCombat struct {
	Attack      float32
	Defense     float32
	CritRate    float32
	CritDamage  float32
	AttackSpeed float32
}

// ==========================
// OOP实体定义（继承式组合）
// ==========================

// OOPEntity 基础实体（包含所有可能组件的指针）
type OOPEntity struct {
	ID        int64
	Transform *OOPTransform
	Velocity  *OOPVelocity
	Physics   *OOPPhysics
	Render    *OOPRender
	Health    *OOPHealth
	AI        *OOPAI
	Movement  *OOPMovement
	Combat    *OOPCombat
}

// OOPEntityMoveable 可移动实体（直接包含组件）
type OOPEntityMoveable struct {
	ID        int64
	Transform OOPTransform
	Velocity  OOPVelocity
}

// OOPEntityPhysical 物理实体（直接包含组件）
type OOPEntityPhysical struct {
	ID        int64
	Transform OOPTransform
	Velocity  OOPVelocity
	Physics   OOPPhysics
}

// OOPEntityRenderable 可渲染实体（直接包含组件）
type OOPEntityRenderable struct {
	ID        int64
	Transform OOPTransform
	Render    OOPRender
}

// OOPEntityActor 完整的游戏对象（直接包含多个组件）
type OOPEntityActor struct {
	ID        int64
	Transform OOPTransform
	Velocity  OOPVelocity
	Physics   OOPPhysics
	Render    OOPRender
	Health    OOPHealth
	AI        OOPAI
	Movement  OOPMovement
	Combat    OOPCombat
}

// ==========================
// OOP系统/管理器
// ==========================

// OOPWorld OOP方式的世界管理器
type OOPWorld struct {
	Entities           []*OOPEntity
	MoveableEntities   []*OOPEntityMoveable
	PhysicalEntities   []*OOPEntityPhysical
	RenderableEntities []*OOPEntityRenderable
	ActorEntities      []*OOPEntityActor
}

// NewOOPWorld 创建OOP世界
func NewOOPWorld() *OOPWorld {
	return &OOPWorld{
		Entities:           make([]*OOPEntity, 0),
		MoveableEntities:   make([]*OOPEntityMoveable, 0),
		PhysicalEntities:   make([]*OOPEntityPhysical, 0),
		RenderableEntities: make([]*OOPEntityRenderable, 0),
		ActorEntities:      make([]*OOPEntityActor, 0),
	}
}

// NewOOPMoveable 创建可移动实体
func (w *OOPWorld) NewOOPMoveable(id int64) *OOPEntityMoveable {
	e := &OOPEntityMoveable{ID: id}
	w.MoveableEntities = append(w.MoveableEntities, e)
	return e
}

// NewOOPPhysical 创建物理实体
func (w *OOPWorld) NewOOPPhysical(id int64) *OOPEntityPhysical {
	e := &OOPEntityPhysical{ID: id}
	w.PhysicalEntities = append(w.PhysicalEntities, e)
	return e
}

// NewOOPRenderable 创建可渲染实体
func (w *OOPWorld) NewOOPRenderable(id int64) *OOPEntityRenderable {
	e := &OOPEntityRenderable{ID: id}
	w.RenderableEntities = append(w.RenderableEntities, e)
	return e
}

// NewOOPActor 创建完整游戏对象
func (w *OOPWorld) NewOOPActor(id int64) *OOPEntityActor {
	e := &OOPEntityActor{ID: id}
	w.ActorEntities = append(w.ActorEntities, e)
	return e
}

// NewOOPEntity 创建通用实体（指针组合）
func (w *OOPWorld) NewOOPEntity(id int64) *OOPEntity {
	e := &OOPEntity{ID: id}
	w.Entities = append(w.Entities, e)
	return e
}

// ==========================
// OOP系统更新方法
// ==========================

// UpdateMoveableEntities 更新可移动实体位置
func (w *OOPWorld) UpdateMoveableEntities(deltaSeconds float32) {
	for _, e := range w.MoveableEntities {
		e.Transform.PosX += e.Velocity.LinearX * deltaSeconds
		e.Transform.PosY += e.Velocity.LinearY * deltaSeconds
		e.Transform.PosZ += e.Velocity.LinearZ * deltaSeconds
		e.Transform.RotX += e.Velocity.AngularX * deltaSeconds
		e.Transform.RotY += e.Velocity.AngularY * deltaSeconds
		e.Transform.RotZ += e.Velocity.AngularZ * deltaSeconds
	}
}

// UpdatePhysicalEntities 更新物理实体
func (w *OOPWorld) UpdatePhysicalEntities(deltaSeconds float32) {
	for _, e := range w.PhysicalEntities {
		// 简单物理模拟
		if e.Physics.UseGravity == 1 {
			e.Velocity.LinearY -= 9.8 * deltaSeconds
		}
		// 阻力
		e.Velocity.LinearX *= (1 - e.Physics.Drag*deltaSeconds)
		e.Velocity.LinearY *= (1 - e.Physics.Drag*deltaSeconds)
		e.Velocity.LinearZ *= (1 - e.Physics.Drag*deltaSeconds)
		// 更新位置
		e.Transform.PosX += e.Velocity.LinearX * deltaSeconds
		e.Transform.PosY += e.Velocity.LinearY * deltaSeconds
		e.Transform.PosZ += e.Velocity.LinearZ * deltaSeconds
	}
}

// UpdateActorEntities 更新完整Actor实体（包含移动、物理、AI、战斗等）
func (w *OOPWorld) UpdateActorEntities(deltaSeconds float32) {
	for _, e := range w.ActorEntities {
		// 移动更新
		e.Transform.PosX += e.Velocity.LinearX * deltaSeconds
		e.Transform.PosY += e.Velocity.LinearY * deltaSeconds
		e.Transform.PosZ += e.Velocity.LinearZ * deltaSeconds

		// 物理更新
		if e.Physics.UseGravity == 1 {
			e.Velocity.LinearY -= 9.8 * deltaSeconds
		}

		// 生命恢复
		if e.Health.Current < e.Health.Max {
			e.Health.Current += e.Health.Regen * deltaSeconds
			if e.Health.Current > e.Health.Max {
				e.Health.Current = e.Health.Max
			}
		}
	}
}

// UpdateGenericEntities 更新通用实体（使用指针，需要空检查）
func (w *OOPWorld) UpdateGenericEntities(deltaSeconds float32) {
	for _, e := range w.Entities {
		// 需要检查组件是否存在
		if e.Transform != nil && e.Velocity != nil {
			e.Transform.PosX += e.Velocity.LinearX * deltaSeconds
			e.Transform.PosY += e.Velocity.LinearY * deltaSeconds
			e.Transform.PosZ += e.Velocity.LinearZ * deltaSeconds
		}

		if e.Physics != nil && e.Velocity != nil {
			if e.Physics.UseGravity == 1 {
				e.Velocity.LinearY -= 9.8 * deltaSeconds
			}
		}

		if e.Health != nil {
			if e.Health.Current < e.Health.Max {
				e.Health.Current += e.Health.Regen * deltaSeconds
			}
		}
	}
}
