package main

import (
	"github.com/zllangct/ecs"
	"math/rand"
)

type Caster struct {
	Action   *Action
	Position *Position
	Force    *Force
}

type Target struct {
	HealthPoint *HealthPoint
	Position    *Position
}

type DamageSystem struct {
	ecs.System[DamageSystem]
	casterGetter *ecs.Shape[Caster]
	targetGetter *ecs.Shape[Target]
}

func (d *DamageSystem) Init(si ecs.SystemInitializer) {
	d.SetRequirements(
		si,
		&ecs.ReadOnly[Position]{},
		&ecs.ReadOnly[Force]{},
		&ecs.ReadOnly[Action]{},
		&HealthPoint{})
	d.casterGetter = ecs.NewShape[Caster](si).SetGuide(&Action{})
	d.targetGetter = ecs.NewShape[Target](si)

	if !d.casterGetter.IsValid() || !d.targetGetter.IsValid() {
		si.SetBroken("invalid shape getter")
	}
}

// Update will be called every frame
func (d *DamageSystem) Update(event ecs.Event) {
	/*  伤害结算逻辑
	  - 分析'攻击'行为，我们需要考虑攻击的相关计算公式所需组件，如当前示例中的Force组件，包含基础攻击、力量、
	暴击倍率、暴击率、攻击范围等数据。考虑攻击范围时需要，知道位置相关信息，由Position组件提供数据支持，在全遍历
	所有位置关系时，消耗比较大，可通过AOI优化，减小遍历规模，优化搜索效率，此处示例不做额外处理。
	*/
	reporter := d.World().GetMetrics().NewReporter("damage system")
	reporter.Start()
	casterIter := d.casterGetter.Get()
	targetIter := d.targetGetter.Get()
	count := 0
	for caster := casterIter.Begin(); !casterIter.End(); caster = casterIter.Next() {
		count++
		for target := targetIter.Begin(); !targetIter.End(); target = targetIter.Next() {
			if caster.Action.Owner() == target.HealthPoint.Owner() {
				continue
			}

			//计算距离
			distance := (caster.Position.X-target.Position.X)*(caster.Position.X-target.Position.X) +
				(caster.Position.Y-target.Position.Y)*(caster.Position.Y-target.Position.Y)
			if distance > caster.Force.AttackRange*caster.Force.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.Force.PhysicalBaseAttack + caster.Force.Strength
			critical := 0
			if rand.Intn(100) < caster.Force.CriticalChange {
				critical = caster.Force.PhysicalBaseAttack * caster.Force.CriticalMultiple
			}
			damage = damage + critical
			//ecs.Log.Infof("Damage:%v", damage)
			target.HealthPoint.HP -= damage
			if target.HealthPoint.HP < 0 {
				target.HealthPoint.HP = 0
			}
		}
	}
	reporter.Sample("damage")
	reporter.Stop()
	reporter.Print()
}
