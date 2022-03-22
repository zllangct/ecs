package main

import (
	"github.com/zllangct/ecs"
	"math/rand"
)

type Caster struct {
	A *Action
	F *Force
	P *Position
	E *ecs.EntityInfo
}

type Target struct {
	P  *Position
	HP *HealthPoint
}

type DamageSystem struct {
	ecs.System[DamageSystem, *DamageSystem]
}

func (d *DamageSystem) Init() {
	d.SetRequirements(&Position{}, &HealthPoint{}, &Force{}, &Action{})
}

func (d *DamageSystem) DataMatch() ([]Caster, []Target) {
	iterAction := ecs.GetInterestedComponents[Action](d)
	if iterAction.Empty() {
		return nil, nil
	}

	//ecs.Log.Infof("iterAction: %v", iterAction.Empty())

	idTemp := map[ecs.Entity]struct{}{}
	var casters []Caster
	for a := iterAction.Begin(); !iterAction.End(); a = iterAction.Next() {
		caster := a.Owner()
		p := ecs.GetRelatedComponent[Position](d, caster)
		if p == nil {
			continue
		}
		f := ecs.GetRelatedComponent[Force](d, caster)
		if f == nil {
			continue
		}
		casters = append(casters, Caster{
			A: a,
			P: p,
			F: f,
			E: caster,
		})
		idTemp[caster.Entity()] = struct{}{}
	}

	iterPos := ecs.GetInterestedComponents[Position](d)
	if iterPos.Empty() {
		return nil, nil
	}
	var targets []Target
	for p := iterPos.Begin(); !iterPos.End(); p = iterPos.Next() {
		target := p.Owner()
		//if _, ok := idTemp[target.entity()]; ok {
		//	continue
		//}

		hp := ecs.GetRelatedComponent[HealthPoint](d, target)
		if hp == nil {
			continue
		}

		targets = append(targets, Target{
			P:  p,
			HP: hp,
		})
	}

	return casters, targets
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

	/*  聚合数据
	  - 分析'攻击'行为，我们需要考虑攻击的相关计算公式所需组件，如当前示例中的Force组件，包含基础攻击、力量、
	暴击倍率、暴击率、攻击范围等数据。考虑攻击范围时需要，知道位置相关信息，由Position组件提供数据支持，在全遍历
	所有位置关系时，消耗比较大，可通过AOI优化，减小遍历规模，优化搜索效率，此处示例不做额外处理。
	聚合数据时，需要匹配攻击者的位置、攻击基础信息，需要匹配被攻击者的位置和血量组件。此外提醒，Entity本身是聚合
	了'个体'的相关组件，天然聚合了相关组件，可以通过Entity得到攻击者或被攻击者的需要聚合的组件。
	*/
	//ecs.Log.Info("DamageSystem Update")
	casters, targets := d.DataMatch()

	// 伤害主逻辑
	for _, caster := range casters {
		for _, target := range targets {
			if caster.A.Owner().Entity() == target.P.Owner().Entity() {
				continue
			}
			//计算距离
			distance := Distance2D(caster.P, target.P)
			if distance > caster.F.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.F.PhysicalBaseAttack + caster.F.Strength
			critical := 0
			if rand.Intn(100) < caster.F.CriticalChange {
				critical = caster.F.PhysicalBaseAttack * caster.F.CriticalMultiple
			}
			damage = damage + critical
			//ecs.Log.Infof("Damage:%v", damage)
			target.HP.HP -= damage
			if target.HP.HP < 0 {
				target.HP.HP = 0
			}
		}
	}
}
