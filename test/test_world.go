package main

import (
	"encoding/json"
	"fmt"
	"github.com/zllangct/ecs"
	"reflect"
	"time"
)

/* 一般设计思路
	ecs有别于面向对象，ecs是数据驱动的设计思路，且提倡组合优于继承的设计原则。在面向对象盛行
的情况下，开发者可能会先入为主的用面向对象的逻辑去思考、设计逻辑，当切换到ecs时，可能会有更高
的心智成本。ecs不是通用的设计方式，在接近‘曾经’很简单的问题时，可能变得更复杂，‘曾经’很复杂的
问题可能在使用ecs后得到简单的解决或带来性能的提升。

1. 确定系统所需操作的组件
	- 系统设计需要考虑组件的组合关系，同时系统的设计也要反哺组件设计，合理配置组件中数据的粒
      度，因为组件的本质是一些列数据的组合，组件的数据粒度越小，系统使用越灵活，但同时可能被
      更多的系统所使用。理想状态是一个组件一个系统，实际业务中很难达到，越复杂的业务逻辑越难，
      尽力而为。
	- 通过SetRequirement()方法设置系统所需的组件，系统内仅能操作设置过的组件，以隔离数据
    - GetInterestedComponents() 仅能获取到所需组件列表内的组件，有可能操作到同一个组件
      的系统都会在同一线程执行，保证线程安全，无序加锁

2. 聚合匹配数据
	- 由于Entity本身聚合了组件，顾系统操作同一个Entity的若干个组件时,因为这些组件天然被对
	  应的Entity所聚合，可直接作为聚合过程，但Entity操作Component有加解锁过程，顾建议独
      立聚合过程。
    - 聚合后数据相互独立，互不影响，无组件间的数据竞争,可进一步在主逻辑中并行化

3. 系统执行处理逻辑
	- 主逻辑的设计，对数据操作，尽可能减少状态量，维持聚合数据间的独立性，可无锁化多线程处理
      逻辑中存在数据竞争时，需要加锁保护状态量，或者串行执行，但并不建议系统内继续并行化，因
      系统间本身已经经过并行处理，单个系统仅位于同一个线程，系统内并行，只为解决超耗时的系统
      具备多线程能力，充分利用系统资源。
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
	logger     ecs.IInternalLogger
}

func (m *MoveSystem) Init() {
	//m.logger = m.GetWorld().logger
	m.SetRequirements(Position{}, Movement{})
}

func (m *MoveSystem) Filter(ls map[reflect.Type][]ecs.ComponentOptResult) {
	if len(ls) > 0 {
		ecs.Log.Info("new component:", len(ls))
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
		移动系统关心的是Position、Movement两个组件，Position包含位置数据，Movement中包含速度、方向数据，
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
	for iter := ecs.NewIterator(csPosition); !iter.End(); iter.Next() {
		c := iter.Val()
		cb , _ := json.Marshal(c)
		ecs.Log.Info("position:", string(cb))
		if cd, ok := d[c.Owner().ID()]; ok {
			cd.P = c
		}else {
			d[c.Owner().ID()] = &MoveSystemData{P: c}
		}
	}
	ecs.Log.Info(csMovement.Len())
	for iter := ecs.NewIterator(csMovement); !iter.End(); iter.Next() {
		c := iter.Val()
		cb , _ := json.Marshal(c)
		_=cb

		ecs.Log.Info("movement: id:",c.Owner().ID(), string(cb), "type:", reflect.TypeOf(c))

		fmt.Printf("%+v \n", *c) //TODO 无法正确获取到 Movement

		if cd, ok := d[c.Owner().ID()]; ok {
			cd.M = c
		}else{
			d[c.Owner().ID()] = &MoveSystemData{M: c}
		}
	}
	b,_:= json.Marshal(d)
	ecs.Log.Info(string(b))
	/* MoveSystem 主逻辑
		- 移动计算公式：最终位置 = 当前位置 + 移动速度 * 移动方向
		- 本系统的移动逻辑非常简单，但可以进一步思考，各个实体的移动是相互独立的，无数据竞争，类似的
	      情况还有很多，此情况下，或者有意设计成独立逻辑的，且该系统计算量特别大的，导致本帧其他系统
	      已经执行完成，会全部等待本系统完成的情况下，可进一步将操作放入线程池中并行处理，充分利用计算资源。
	*/
	ecs.Log.Info("main logic")
	for e, data := range d {
		if data.M == nil || data.P == nil {
			//聚合数据组件不齐时，跳过处理
			continue
		}
		data.P.X = data.P.X + int(float64(data.M.Dir[0]*data.M.V)*delta.Seconds())
		data.P.Y = data.P.Y + int(float64(data.M.Dir[1]*data.M.V)*delta.Seconds())
		data.P.Z = data.P.Z + int(float64(data.M.Dir[2]*data.M.V)*delta.Seconds())

		ecs.Log.Info("target id:", e, " current position:", data.P.X, data.P.Y, data.P.Z)
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
	Type int
}

type DamageSystem struct {
	ecs.System[DamageSystem]
	actions []ecs.ComponentOptResult
}

func (d *DamageSystem) Init() {
	d.SetRequirements(Position{}, HealthPoint{}, Force{}, Action{})
}


func (d *DamageSystem) Filter(ls map[reflect.Type][]ecs.ComponentOptResult) {
	if len(ls) == 0 {
		return
	}

	as, ok := ls[ecs.GetType[Action]()]
	if !ok {
		return
	}

	d.actions = as
}

func (d *DamageSystem) DataMatch(){

}

// Update will be called every frame
func (d *DamageSystem) Update(event ecs.Event) {
	/*
		获取当前帧新所有创建、删除的组件
	    - 本系统（DamageSystem）根据行为类组件Action来处理实体的攻击操作，其Action.Type 可以
		是普通攻击、主动技能或者被动技能，根据需求开发者自行定义，伤害系统会根据实体是否挂载了该组件
		来决定本帧，该实体是否有攻击行为，当本次攻击完成后，删除该组件即可，这里以普通攻击为例。
	*/
	nc := d.GetInterestedNew()

	/*
		根据上面的描述，当实体被挂载Action组件时，该实体正在攻击，攻击介绍后，组件移除。即，当有Action
		组件被添加或者被移除时，是我们关心的，可在Filter中预处理、缓存相关数据，便于下一步聚合数据，快速
		筛选出关心的数据。
	*/
	d.Filter(nc)


	//todo logic 待续

}

//main function
func Runtime0() {
	// pprof
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:8888", nil))
	//}()

	// 创建运行时，运行时唯一
	rt := ecs.Runtime
	rt.Run()

	// 创建世界，世界可以多个
	world := rt.NewWorld()
	world.Run()

	// 注册系统
	world.Register(&MoveSystem{})

	// 创建实体并添加组件
	ee1 := world.NewEntity()
	ee2 := world.NewEntity()
	ee3 := world.NewEntity()

	ecs.Log.Info(ee1.ID(), ee2.ID(), ee3.ID())

	p1 := &Position{
		X: 100,
		Y: 100,
		Z: 100,
	}
	m1 := &Movement{
		V:   2000,
		Dir: []int{1,0,0},
	}
	world.NewEntity().AddByTemplate(p1, m1)
	//p2 := &Position{
	//	X: 100,
	//	Y: 100,
	//	Z: 100,
	//}
	//m2 := &Movement{
	//	V: 2000,
	//	Dir: []int{0,1,0},
	//}
	//e2 := world.NewEntity()
	//e2.AddByTemplate(p2, m2)

	//ecs.Log.Info("test entity id e2:", e2.ID())

	time.Sleep(time.Second * 1)
	//for {
	//	time.Sleep(time.Second * 3)
	//}
}

func main() {
	Runtime0()
}