package main

import (
	"container/list"
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

type CollectionKV struct {
	Entity *Entity
	Data   []interface{}
}

type ISystemBase interface {
	Inject(sys ISystem)
	Filter(entity *Entity) //筛选感兴趣的组件
}

type ISystem interface {
	ISystemBase
	Init() //初始化
	Run()  //执行系统逻辑
}

type SystemBase struct {
	sync.RWMutex
	//private
	this         ISystem
	typ          SystemType
	// traverse ratio: slice = 10 * map
	es           map[*Entity]int
	data         [][]interface{}
	pop          *list.List
	push         *list.List
	requirements []reflect.Type
}

func (p *SystemBase) Inject(sys ISystem) {
	//Dependency inject
	p.this = sys
	//TODO init members
}

func (p *SystemBase)PreUpdate()  {
	p.Lock()
	defer p.Unlock()

	for item:=p.push.Front();item != nil ;item = item.Next()  {
		kv:=item.Value.(CollectionKV)
		p.data = append(p.data, kv.Data)
		p.es[kv.Entity] = len(p.data) - 1
	}
	p.push.Init()
	for item:=p.pop.Front();item != nil ;item = item.Next()  {
		kv:=item.Value.(CollectionKV)
		length := len(p.data) - 1
		index, ok := p.es[kv.Entity]
		if !ok {
			continue
		}
		p.data[index], p.data[length] = p.data[length], p.data[index]
		p.data = p.data[:length]
		for key, value := range p.es {
			if value == len(p.data)-1 {
				p.es[key] = index
				break
			}
		}
	}
	p.pop.Init()
}

func (p *SystemBase)Update()  {
	p.PreUpdate()
	//TODO slice task queue task_length / k * cpu_num
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

	p.push.PushBack(CollectionKV{Entity:entity,Data:cmps})

	//	p.data = append(p.data, cmps)
	//	p.es[entity] = len(p.data) - 1


}

func (p *SystemBase) Clean(entity *Entity) {
	p.Lock()
	defer p.Unlock()

	_, ok := p.es[entity]
	if !ok {
		return
	}
	p.pop.PushBack(entity)

	//remove the entity's data
	//	length := len(p.data) - 1
	//	p.data[index], p.data[length] = p.data[length], p.data[index]
	//	p.data = p.data[:length]
	//	for key, value := range p.es {
	//		if value == len(p.data)-1 {
	//			p.es[key] = index
	//			break
	//		}
	//	}


}

type Start struct {
	SystemBase
}

func (p *Start) Init() {

}

func (*Start) Run() {

}
