package main

import (
	"github.com/zllangct/ecs"
)

type MoveSystemData struct {
	P *Position
	M *Movement
}

type MoveSystem struct {
	ecs.System[MoveSystem, *MoveSystem]
}

func (m *MoveSystem) Init() {
	m.SetRequirements(&Position{}, &Movement{})
}

func (m *MoveSystem) Update(event ecs.Event) {
	delta := event.Delta
	//获取系统所需的组件
	//方式 1:
	//iterPos := m.GetInterested(ecs.GetType[Position]()).(*ecs.Collection[Position])
	//方式 2:
	//iterPos := ecs.GetInterestedComponents[Position](m)

	/* 聚合数据
	  - 移动系统关心的是Position、Movement两个组件，Position包含位置数据，Movement中包含速度、方向数据，
	我们可以知道Position和Movement应该是同一个实体的两组件配合使用，组件的配对，在数据聚合阶段完成，可知，
	当两组件拥有相同的Owner时，完成匹配。
	*/
	iterPos := ecs.GetInterestedComponents[Position](m)
	iterMov := ecs.GetInterestedComponents[Movement](m)

	if iterPos.Empty() || iterMov.Empty() {
		return
	}

	d := map[ecs.Entity]*MoveSystemData{}

	//聚合方式 1：根据组件的Owner（EntityInfo.Entity）来匹配数据
	//for iter := iterPos; !iter.End(); iter.Next() {
	//	c := iter.Val()
	//	if cd, ok := d[c.Owner().Entity()]; ok {
	//		cd.P = c
	//	}else {
	//		d[c.Owner().Entity()] = &MoveSystemData{P: c}
	//	}
	//}
	//for iter := iterMov; !iter.End(); iter.Next() {
	//	c := iter.Val()
	//	if cd, ok := d[c.Owner().Entity()]; ok {
	//		cd.M = c
	//	}else{
	//		d[c.Owner().Entity()] = &MoveSystemData{M: c}
	//	}
	//}

	//聚合方式 2：直接从Entity聚合相关组件
	for iter := iterPos; !iter.End(); iter.Next() {
		position := iter.Val()
		owner := position.Owner()
		/*
		  无法通过Entity直接获取到所有组件，故意如此设计，保证在系统中错误修改非必须的组件，GetRelatedComponent
		  能够检查需要获取的组件是否是该系统所必须组件。
		*/
		movement := ecs.GetRelatedComponent[Movement](m, owner)
		if movement == nil {
			continue
		}

		d[position.Owner().Entity()] = &MoveSystemData{P: position, M: movement}
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

		_ = e
		//ecs.Log.Info("target id:", e, " delta:", delta, " current position:", data.P.X, data.P.Y, data.P.Z)
	}
}
