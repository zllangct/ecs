package main

import (
	"reflect"
	"sync"
)

type ISystem interface {
	Filter()                      //筛选感兴趣的组件
	Clean()                       //清理失效的组件
	Run()                         //执行系统逻辑
	Requirements() []reflect.Type //系统需要的组件
}

type SystemBase struct {
	sync.RWMutex
	es map[*Entity]struct{}
	data
}


func (p *SystemBase) Requirements() []reflect.Type {
	panic("implement me")
}

func (p *SystemBase) Filter(entity *Entity) {
	//check exist
	if _,ok:=p.es[entity];ok {
		return
	}
	//check requirements
	rqs:=p.Requirements()
	for _, rq := range rqs {
		if !entity.Has(rq) {
			return
		}
	}
	//generate data collection


}

func (*SystemBase) Clean() {
	panic("implement me")
}

func (*SystemBase) Run() {
	panic("implement me")
}
