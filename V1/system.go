package main

import (
	"errors"
	"reflect"
	"sync"
)

type SystemType int

const (
	SYSTEM_TYPE_NONE SystemType = iota
	SYSTEM_TYPE_Default
	SYSTEM_TYPE_ONCE
	SYSTEM_TYPE_REPEAT
)

var (
	ErrSystemNotInit = errors.New("system not init")
)

type ISystemBase interface {
	Inject(sys ISystem)
	Filter(entity *Entity) //筛选感兴趣的组件
	Clean(entity *Entity)  //清理失效的组件
}

type ISystem interface {
	ISystemBase
	Init() //初始化
	Run()  //执行系统逻辑
}

type SystemBase struct {
	sync.RWMutex
	this         ISystem
	typ          SystemType
	es           map[*Entity]int
	data         [][]interface{}
	data0        map[*Entity][]interface{}
	requirements []reflect.Type
}

func (p *SystemBase) Inject(sys ISystem) {
	//Dependency inject
	p.this = sys
	//TODO init members
}

func (p *SystemBase) Filter(entity *Entity) {
	//check exist
	p.RLock()
	if _, ok := p.es[entity]; ok {
		p.RUnlock()
		return
	}
	p.RUnlock()

	//check requirements
	for _, rq := range p.requirements {
		if !entity.Has(rq) {
			return
		}
	}
	//generate data collection
	cmps := entity.GetComponents(p.requirements...)

	//add data
	p.Lock()
	defer p.Unlock()

	switch p.typ {
	case SYSTEM_TYPE_Default | SYSTEM_TYPE_REPEAT:
		p.data = append(p.data, cmps)
		p.es[entity] = len(p.data) - 1
	case SYSTEM_TYPE_ONCE:
		p.data0[entity] = cmps
	default:
		panic(ErrSystemNotInit)
	}

}

func (p *SystemBase) Clean(entity *Entity) {
	p.Lock()
	defer p.Unlock()

	index, ok := p.es[entity]
	if !ok {
		return
	}

	//remove the entity's data
	switch p.typ {
	case SYSTEM_TYPE_Default | SYSTEM_TYPE_REPEAT:
		length := len(p.data) - 1
		p.data[index], p.data[length] = p.data[length], p.data[index]
		p.data = p.data[:length]
		for key, value := range p.es {
			if value == len(p.data)-1 {
				p.es[key] = index
				break
			}
		}
	case SYSTEM_TYPE_ONCE:
		delete(p.data0, entity)
	default:
		panic(ErrSystemNotInit)
	}

}

type Start struct {
	SystemBase
}

func (p *Start) Init() {

}

func (*Start) Run() {
	panic("implement me")
}
