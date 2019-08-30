package main

import (
	"container/list"
	"errors"
	"reflect"
	"sync"
)

type SystemLifeCircleType int
const (
	SYSTEM_LIFE_CIRCLE_TYPE_NONE SystemLifeCircleType = iota
	SYSTEM_LIFE_CIRCLE_TYPE_Default
	SYSTEM_LIFE_CIRCLE_TYPE_ONCE
	SYSTEM_LIFE_CIRCLE_TYPE_REPEAT
)

var (
	ErrSystemNotInit = errors.New("system not init")
)

type ISystem interface {
	Type() reflect.Type    //类型
	Init()                 //初始化
	Filter(entity *Entity) //筛选感兴趣的组件
	FrameUpdate()          //执行系统逻辑
}

type SystemBase struct {
	sync.RWMutex
	typ          SystemLifeCircleType
	pop          *list.List
	push         *list.List
	requirements []reflect.Type
}

func (p *SystemBase)SetOrder()  {

}

func (p *SystemBase)SetLifeCircleType()  {

}

type Start struct {
	SystemBase
	//private
	entityIndex  map[*Entity]int
	data         [][]interface{}
}

func (p *Start) Type() reflect.Type {
	return reflect.TypeOf(p)
}

func (p *Start) Init() {
	//TODO init members
}

func (p *Start) PreUpdate() {
	p.Lock()
	defer p.Unlock()

	for item := p.push.Front(); item != nil; item = item.Next() {
		kv := item.Value.(CollectionKV)
		p.data = append(p.data, kv.Data)
		p.es[kv.Entity] = len(p.data) - 1
	}
	p.push.Init()
	for item := p.pop.Front(); item != nil; item = item.Next() {
		kv := item.Value.(CollectionKV)
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

func (p *Start) FrameUpdate() {
	p.PreUpdate()
	//TODO slice task queue task_length / k * cpu_num

	//TODO 检查 COMPONENT_STATE_NONE 状态下，是否执行，
}

func (p *Start) Filter(entity *Entity) {
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

	p.push.PushBack(CollectionKV{Entity: entity, Data: cmps})
}

func (p *Start) Clean(entity *Entity) {
	p.Lock()
	defer p.Unlock()

	_, ok := p.es[entity]
	if !ok {
		return
	}
	p.pop.PushBack(entity)
}
