package main

import (
	"github.com/zllangct/ecs"
	"math/rand"
	"reflect"
	"time"
)

/* 一般设计思路

    - ecs有别于面向对象，ecs是数据驱动的设计思路，且提倡组合优于继承的设计原则。在面向对象盛行
  的情况下，开发者可能会先入为主的用面向对象的逻辑去思考、设计逻辑，当切换到ecs时，可能会有更高
  的心智成本。ecs不是通用的设计方式，在接近‘曾经’很简单的问题时，可能变得更复杂，‘曾经’很复杂的
  问题可能在使用ecs后得到简单的解决或带来性能的提升。

1. 确定系统所需操作的组件
    - 系统设计需要考虑组件的组合关系，同时系统的设计也要反哺组件设计，合理配置组件中数据的粒
  度，因为组件的本质是一些列数据的组合，组件的数据粒度越小，系统使用越灵活，但同时可能被更多的系统
  所使用。理想状态是一个组件一个系统，实际业务中很难达到，越复杂的业务逻辑越难，
  尽力而为。
    - 通过SetRequirement()方法设置系统所需的组件，系统内仅能操作设置过的组件，以隔离数据
    - GetInterestedComponents() 仅能获取到所需组件列表内的组件，有可能操作到同一个组件
  的系统都会在同一线程执行，保证线程安全，无序加锁

2. 聚合匹配数据
    - 由于Entity本身聚合了组件，顾系统操作同一个Entity的若干个组件时,因为这些组件天然被对应
  的Entity所聚合，可直接作为聚合过程.
    - 聚合后数据相互独立，互不影响，无组件间的数据竞争,可进一步在主逻辑中并行化

3. 系统执行处理逻辑
    - 主逻辑的设计，对数据操作，尽可能减少状态量，维持聚合数据间的独立性，可无锁化多线程处理
  逻辑中存在数据竞争时，需要加锁保护状态量，或者串行执行，但并不建议系统内继续并行化，因系统间
  本身已经经过并行处理，单个系统仅位于同一个线程，系统内并行，只为解决超耗时的系统具备多线程能
  力，充分利用系统资源。
*/


//position component
type Position struct {
	ecs.Component[Position]
	X int
	Y int
	Z int
}

type Movement struct {
	ecs.Component[Movement]
	V   int
	Dir []int
}

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	logger     ecs.ILogger
}

func (m *MoveSystem) Init() {
	//m.logger = m.GetWorld().logger
	m.SetRequirements(&Position{}, &Movement{})
}

func (m *MoveSystem) Filter(ls map[reflect.Type][]ecs.OperateInfo) {
	if len(ls) > 0 {
		//ecs.Log.Info("new component:", len(ls))
	}
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta

	//当前帧新所有创建、删除的组件
	nc := m.GetInterestedNew()
	m.Filter(nc)

	//获取系统所需的组件
	//方式 1:
	//csPosition := m.GetInterested(ecs.GetType[Position]()).(*ecs.Collection[Position])
	//方式 2:
	//csPosition := ecs.GetInterestedComponents[Position](m)

	/* 聚合数据
	    - 移动系统关心的是Position、Movement两个组件，Position包含位置数据，Movement中包含速度、方向数据，
	  我们可以知道Position和Movement应该是同一个实体的两组件配合使用，组件的配对，在数据聚合阶段完成，可知，
	  当两组件拥有相同的Owner时，完成匹配。
	*/
	csPosition := ecs.GetInterestedComponents[Position](m)
	if csPosition == nil {
		return
	}
	csMovement := ecs.GetInterestedComponents[Movement](m)
	if csMovement == nil {
		return
	}

	d := map[int64]*MoveSystemData{}

	//聚合方式 1：根据组件的Owner（Entity.GetID）来匹配数据
	//for iter := ecs.NewIterator(csPosition); !iter.End(); iter.Next() {
	//	c := iter.Val()
	//	if cd, ok := d[c.Owner().GetID()]; ok {
	//		cd.P = c
	//	}else {
	//		d[c.Owner().GetID()] = &MoveSystemData{P: c}
	//	}
	//}
	//for iter := ecs.NewIterator(csMovement); !iter.End(); iter.Next() {
	//	c := iter.Val()
	//	if cd, ok := d[c.Owner().GetID()]; ok {
	//		cd.M = c
	//	}else{
	//		d[c.Owner().GetID()] = &MoveSystemData{M: c}
	//	}
	//}

	//聚合方式 2：直接从Entity聚合相关组件
	for iter := ecs.NewIterator(csPosition); !iter.End(); iter.Next() {
		position := iter.Val()
		owner := position.Owner()
		/*
		  无法通过Entity直接获取到所有组件，故意如此设计，保证在系统中错误修改非必须的组件，CheckComponent
		  能够检查需要获取的组件是否是该系统所必须组件。
		 */
		movement := ecs.CheckComponent[Movement](m, owner)
		if movement == nil {
			continue
		}

		d[position.Owner().GetID()] = &MoveSystemData{P: position, M: movement}
	}

	/* MoveSystem 主逻辑
	    - 移动计算公式：最终位置 = 当前位置 + 移动速度 * 移动方向
		- 本系统的移动逻辑非常简单，但可以进一步思考，各个实体的移动是相互独立的，无数据竞争，类似的
	  情况还有很多，此情况下，或者有意设计成独立逻辑的，且该系统计算量特别大的，导致本帧其他系统
	  已经执行完成，会全部等待本系统完成的情况下，可进一步将操作放入线程池中并行处理，充分利用计算资源。
	*/

	for e, data := range d {
		if data.M == nil || data.P == nil {
			//聚合数据组件不齐时，跳过处理
			continue
		}
		data.P.X = data.P.X + int(float64(data.M.Dir[0]*data.M.V)*delta.Seconds())
		data.P.Y = data.P.Y + int(float64(data.M.Dir[1]*data.M.V)*delta.Seconds())
		data.P.Z = data.P.Z + int(float64(data.M.Dir[2]*data.M.V)*delta.Seconds())

		ecs.Log.Info("target id:", e, "delta:", delta, " current position:", data.P.X, data.P.Y, data.P.Z)
	}
}

//hp component
type HealthPoint struct {
	ecs.Component[HealthPoint]
	HP int
}

type Force struct {
	ecs.Component[Force]
	AttackRange int
	PhysicalBaseAttack int
	Strength int
	CriticalChange int
	CriticalMultiple int
}

type Action struct {
	ecs.Component[Action]
	ActionType int
}

type Caster struct {
	A *Action
	F *Force
	P *Position
	E *ecs.Entity
}

type Target struct {
	P *Position
	HP *HealthPoint
}

type DamageSystem struct {
	ecs.System[DamageSystem]
}

func (d *DamageSystem) Init() {
	d.SetRequirements(&Position{}, &HealthPoint{}, &Force{}, &Action{})
}

func (d *DamageSystem) Filter() []ecs.OperateInfo{
	/*
	  获取当前帧新所有创建、删除的组件
	*/
	nc := d.GetInterestedNew()
	if len(nc) == 0 {
		return nil
	}

	as, ok := nc[ecs.GetType[Action]()]
	if !ok {
		return nil
	}

	return as
}

func (d *DamageSystem) DataMatch() ([]Caster, []Target) {
	action := ecs.GetInterestedComponents[Action](d)
	if action == nil {
		return nil, nil
	}

	idTemp := map[int64]struct{}{}
	var casters []Caster
	iter:= ecs.NewIterator(action)
	for a := iter.Begin(); !iter.End(); iter.Next() {
		caster := a.Owner()
		p := ecs.CheckComponent[Position](d, caster)
		if p == nil {
			continue
		}
		f := ecs.CheckComponent[Force](d, caster)
		if f == nil {
			continue
		}
		casters = append(casters, Caster{
			A: a,
			P: p,
			F: f,
			E: caster,
		})
		idTemp[caster.GetID()] = struct{}{}
	}

	position := ecs.GetInterestedComponents[Position](d)
	if position == nil {
		return nil, nil
	}
	var targets []Target
	pIter := ecs.NewIterator(position)
	for p := pIter.Begin(); !pIter.End(); pIter.Next() {
		target := p.Owner()
		if _, ok := idTemp[target.GetID()]; ok {
			continue
		}

		hp := ecs.CheckComponent[HealthPoint](d, target)
		if hp == nil {
			continue
		}

		targets = append(targets, Target{
			P: p,
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
	casters, targets := d.DataMatch()

	// 伤害主逻辑
	for _, caster := range casters{
		for _, target := range targets {
			//计算距离
			distance := Distance2D(caster.P, target.P)
			if distance > caster.F.AttackRange {
				continue
			}

			//伤害公式：伤害=（基础攻击+力量）+ 暴击伤害， 暴击伤害=基础攻击 * 2
			damage := caster.F.PhysicalBaseAttack * caster.F.Strength
			critical := 0
			if rand.Intn(100) < caster.F.CriticalChange {
				critical = caster.F.PhysicalBaseAttack * caster.F.CriticalMultiple
			}
			damage = damage + critical
			target.HP.HP -= damage
			if target.HP.HP < 0 {
				target.HP.HP =0
			}
		}

		caster.E.Remove(caster.A)

	}

}

//main function
func Runtime0() {
	// 配置运行时，运行时唯一
	ecs.RuntimeConfigure(ecs.NewDefaultRuntimeConfig())
	ecs.Run()

	// 创建世界，世界可以多个
	world := ecs.CreateWorld(ecs.NewDefaultWorldConfig())
	world.Run()

	// 注册系统
	//world.Register(&MoveSystem{})
	ecs.RegisterSystem[MoveSystem](world)


	// 创建实体并添加组件
	ee1 := world.NewEntity()
	ee2 := world.NewEntity()
	ee3 := world.NewEntity()

	ecs.Log.Info(ee1.GetID(), ee2.GetID(), ee3.GetID())

	p1 := &Position{
		X: 100,
		Y: 100,
		Z: 100,
	}
	m1 := &Movement{
		V:   2000,
		Dir: []int{1,0,0},
	}
	world.NewEntity().Add(p1, m1)
	p2 := &Position{
		X: 100,
		Y: 100,
		Z: 100,
	}
	m2 := &Movement{
		V: 2000,
		Dir: []int{0,1,0},
	}
	world.NewEntity().Add(p2, m2)

	//示例仅运行1秒
	time.Sleep(time.Second * 1)
}

func main() {
	Runtime0()
}
