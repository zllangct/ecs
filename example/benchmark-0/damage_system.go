package main

import (
	"github.com/zllangct/ecs"
	"math/rand"
)

type DamageSystem struct {
	ecs.System[DamageSystem, *DamageSystem]
	casterGetter *ecs.ShapeGetter[ecs.Shape3[Action, Position, Force],
		*ecs.Shape3[Action, Position, Force]]
	targetGetter *ecs.ShapeGetter[ecs.Shape2[HealthPoint, Position],
		*ecs.Shape2[HealthPoint, Position]]
}

func (d *DamageSystem) Init() {
	d.SetRequirements(
		&ecs.ReadOnly[Position]{},
		&ecs.ReadOnly[Force]{},
		&ecs.ReadOnly[Action]{},
		&HealthPoint{})
	var err error
	d.casterGetter, err = ecs.NewShapeGetter[ecs.Shape3[Action, Position, Force]](d)
	if err != nil {
		ecs.Log.Error(err)
	}
	d.targetGetter, err = ecs.NewShapeGetter[ecs.Shape2[HealthPoint, Position]](d)
	if err != nil {
		ecs.Log.Error(err)
	}
}

// Update will be called every frame
func (d *DamageSystem) Update(event ecs.Event) {
	/*
	    - 当前帧中挂载Action组件的实体即为正在进行攻击的实体，攻击结束后Action组件移除，符合'组件'即'能力'
	  的设计思路，挂载攻击组件（Action）具备攻击能力，移除该组件后，失去攻击能力。当然这不是唯一的方法，
	  当前example中攻击每帧结算，由于ecs系统中有缓存每帧新添加、新移除等组件操作，符合攻击每帧结算的特征. 方法
	  很多，如果考虑攻击间隔时，也可以将cooldown作为Attack组件的一个数据项，系统中进行逻辑判断，这样可以无需Action组件，
	  此处使用Action方式，仅简单处理。
	*/

	/*  伤害结算逻辑
	  - 分析'攻击'行为，我们需要考虑攻击的相关计算公式所需组件，如当前示例中的Force组件，包含基础攻击、力量、
	暴击倍率、暴击率、攻击范围等数据。考虑攻击范围时需要，知道位置相关信息，由Position组件提供数据支持，在全遍历
	所有位置关系时，消耗比较大，可通过AOI优化，减小遍历规模，优化搜索效率，此处示例不做额外处理。
	*/
	reporter := ecs.Runtime.GetMetrics().NewReporter("damage system")
	reporter.Start()
	casterIter := d.casterGetter.Iter()
	targetIter := d.targetGetter.Iter()
	_, _ = casterIter, targetIter
	count := 0
	for caster := casterIter.Begin(); !casterIter.End(); caster = casterIter.Next() {
		casterPos := caster.C2
		casterForce := caster.C3
		_, _ = casterPos, casterForce
		count++
		for target := targetIter.Begin(); !targetIter.End(); target = targetIter.Next() {
			if caster.C1.Owner().Entity() == target.C1.Owner().Entity() {
				continue
			}
			targetHp := target.C1
			targetPos := target.C2
			//计算距离
			distance := (casterPos.X-targetPos.X)*(casterPos.X-targetPos.X) + (casterPos.Y-targetPos.Y)*(casterPos.Y-targetPos.Y)
			if distance > casterForce.AttackRange*casterForce.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := casterForce.PhysicalBaseAttack + casterForce.Strength
			critical := 0
			if rand.Intn(100) < casterForce.CriticalChange {
				critical = casterForce.PhysicalBaseAttack * casterForce.CriticalMultiple
			}
			damage = damage + critical
			//ecs.Log.Infof("Damage:%v", damage)
			targetHp.HP -= damage
			if targetHp.HP < 0 {
				targetHp.HP = 0
			}
		}
	}
	reporter.Sample("damage")
	reporter.Stop()
	reporter.Print()
}
