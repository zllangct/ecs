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
	Init()                                 //init
	Type() reflect.Type                    //type of system
	GetOrder() (SystemPeriod, int) //order of system
	Filter(*Entity)                 	   //interest filter of component
	FrameUpdate()                          //update every frame
}

type SystemBase struct {
	sync.RWMutex
	pop          *list.List
	push         *list.List
	requirements []reflect.Type
	order        SystemOrder
}

func (p *SystemBase) SetOrder(order int, period ...SystemPeriod) {
	mPeriod := SYSTEM_PERIOD_DEFAULT
	if len(period) > 0 {
		mPeriod = period[0]
	}
	p.Lock()
	p.order = SystemOrder(mPeriod<<32) + SystemOrder(order)
	p.Unlock()
}

func (p *SystemBase) GetOrder() (SystemPeriod, int) {
	p.RLock()
	defer p.RUnlock()
	return SystemPeriod(p.order >> 32), int(p.order & 0xffff)
}


type Start struct {
	SystemBase
	//private
	entityIndex map[*Entity]int
	data        [][]interface{}
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
		p.entityIndex[kv.Entity] = len(p.data) - 1
	}
	p.push.Init()
	for item := p.pop.Front(); item != nil; item = item.Next() {
		kv := item.Value.(CollectionKV)
		length := len(p.data) - 1
		index, ok := p.entityIndex[kv.Entity]
		if !ok {
			continue
		}
		p.data[index], p.data[length] = p.data[length], p.data[index]
		p.data = p.data[:length]
		for key, value := range p.entityIndex {
			if value == len(p.data)-1 {
				p.entityIndex[key] = index
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
	if _, ok := p.entityIndex[entity]; ok {
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

	_, ok := p.entityIndex[entity]
	if !ok {
		return
	}
	p.pop.PushBack(entity)
}
