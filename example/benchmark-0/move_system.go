package main

import (
	"github.com/zllangct/ecs"
)

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem]
	getter *ecs.ShapeGetter[ecs.Shape2[Movement, Position], *ecs.Shape2[Movement, Position]]
}

func (m *MoveSystem) Init(initializer *ecs.SystemInitializer) {
	m.SetRequirements(initializer, &Position{}, &ecs.ReadOnly[Movement]{})
	getter, err := ecs.NewShapeGetter[ecs.Shape2[Movement, Position]](m)
	if err != nil {
		ecs.Log.Error(err)
	}
	m.getter = getter
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta
	_ = delta
	//获取系统所需的单个组件
	//方式 1:
	//iterPos := m.GetInterested(ecs.GetType[Position]()).(*ecs.UnorderedCollectionWithID[Position])
	//方式 2:
	//iterPos := ecs.GetInterestedComponents[Position](m)

	// 移动系统关心的是Position、Movement两个组件，Position包含位置数据，Movement中包含速度、方向数据，
	// 我们可以知道Position和Movement应该是同一个实体的两组件配合使用

	// 聚合方式 1: 通过ShapeGetter获取同一实体的Position和Movement

	// getter, err := ecs.NewShapeGetter[ecs.Shape2[Movement, Position]](m)
	// iter := m.getter.Get()
	// for shp := iter.Begin(); !iter.End(); shp = iter.Next() {
	//     p := shp.C1
	//	   mv := shp.C2
	//}

	count := 0
	iter := m.getter.Get()
	for shp := iter.Begin(); !iter.End(); shp = iter.Next() {
		mv := shp.C1
		p := shp.C2
		_, _ = p, mv
		p.X = p.X + int(float64(mv.Dir[0]*mv.V)*delta.Seconds())
		p.Y = p.Y + int(float64(mv.Dir[1]*mv.V)*delta.Seconds())
		p.Z = p.Z + int(float64(mv.Dir[2]*mv.V)*delta.Seconds())

		count++
		//e := p.Owner().Entity()
		//ecs.Log.Info("target id:", e, " delta:", delta, " current position:", p.X, p.Y, p.Z)
	}

	//聚合方式 2：直接从Entity聚合相关组件

	//iterPos := ecs.GetInterestedComponents[Position](m)
	//iterMov := ecs.GetInterestedComponents[Movement](m)
	//
	//if iterPos.Empty() || iterMov.Empty() {
	//	return
	//}
	//
	//d := map[ecs.Entity]*MoveSystemData{}
	//
	//for iter := iterPos; !iter.End(); iter.Next() {
	//	position := iter.Val()
	//	owner := position.Owner()
	//	/*
	//	  无法通过Entity直接获取到所有组件，故意如此设计，保证在系统中错误修改非必须的组件，GetRelatedComponent
	//	  能够检查需要获取的组件是否是该系统所必须组件。
	//	*/
	//	movement := ecs.GetRelatedComponent[Movement](m, owner)
	//	if movement == nil {
	//		continue
	//	}
	//
	//	d[position.Owner().Entity()] = &MoveSystemData{P: position, M: movement}
	//}
	//
	///* MoveSystem 主逻辑
	//    - 移动计算公式：最终位置 = 当前位置 + 移动速度 * 移动方向
	//	- 本系统的移动逻辑非常简单，但可以进一步思考，各个实体的移动是相互独立的，无数据竞争，类似的
	//  情况还有很多，此情况下，或者有意设计成独立逻辑的，且该系统计算量特别大的，导致本帧其他系统
	//  已经执行完成，会全部等待本系统完成的情况下，可进一步将操作放入线程池中并行处理，充分利用计算资源。
	//*/
	//
	//for e, data := range d {
	//	if data.M == nil || data.P == nil {
	//		//聚合数据组件不齐时，跳过处理
	//		continue
	//	}
	//	data.P.X = data.P.X + int(float64(data.M.Dir[0]*data.M.V)*delta.Seconds())
	//	data.P.Y = data.P.Y + int(float64(data.M.Dir[1]*data.M.V)*delta.Seconds())
	//	data.P.Z = data.P.Z + int(float64(data.M.Dir[2]*data.M.V)*delta.Seconds())
	//
	//	_ = e
	//	//ecs.Log.Info("target id:", e, " delta:", delta, " current position:", data.P.X, data.P.Y, data.P.Z)
	//}

}
