package main

import "reflect"

type Start struct {
	SystemBase
	//private
	components []IEventStart
}

func (p *Start) Init() {
	//base init
	p.SystemBase.Init()
	//inject system type info
	p.SetType(reflect.TypeOf(p))
	//initialise interest components
	p.components = make([]IEventStart,0,10)
	//set order
	p.SetOrder(ORDER_DEFAULT,PERIOD_PRE_START)
}

func (p *Start) SystemUpdate() {
	interval := len(p.components) / p.runtime.config.CpuNum
	remainder := len(p.components) % p.runtime.config.CpuNum
	offset := 0
	for i := 0; i<p.runtime.config.CpuNum;i++  {
		p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
			for _, event := range args[0].([]IEventStart) {
				event.Start()
			}
		},[]interface{}{p.components[offset:offset+interval]})
		offset += interval
	}
	for i := 0; i<remainder;i++  {
		p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
			args[0].(IEventStart).Start()
		},[]interface{}{p.components[offset]})
		offset+=1
	}
}

func (p *Start) Filter() {
	//clear old data
	p.components = p.components[0:0]
	//update new component
	coms:=p.runtime.components.GetComponentsAdded()
	for _, com := range coms {
		if v,ok:=com.(IEventStart);ok {
			p.components = append(p.components, v)
		}
	}
}
