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
	Init()                //init
	GetBase() *SystemBase //get system common data
	Filter(*Entity)       //interest filter of component
	SystemUpdate()        //update every frame
	GetRequirements() []reflect.Type
	Call(label int)interface{}
}

type SystemBase struct {
	sync.Mutex
	pop          *list.List
	push         *list.List
	requirements []reflect.Type
	order        SystemOrder
	runtime      *Runtime
	typ          reflect.Type
}

func (p *SystemBase) Call(label int)interface{}{
	return nil
}

func (p *SystemBase) SetRequirements(rqs ... reflect.Type){
	p.requirements = rqs
}

func (p *SystemBase) GetRequirements() []reflect.Type {
	return p.requirements
}

func (p *SystemBase)Init()  {
	p.pop = list.New()
	p.push = list.New()
	p.requirements = []reflect.Type{}
}

func (p *Start) SetType(typ reflect.Type){
	p.Lock()
	defer p.Unlock()

	p.typ = typ
}

func (p *Start) GetType() reflect.Type {
	p.Lock()
	defer p.Unlock()

	return reflect.TypeOf(p)
}

func (p *SystemBase) GetBase() *SystemBase {
	return p
}

func (p *SystemBase) SetOrder(order int, period ...SystemPeriod) {
	mPeriod := PERIOD_DEFAULT
	if len(period) > 0 {
		mPeriod = period[0]
	}
	p.Lock()
	p.order = SystemOrder(mPeriod<<32) + SystemOrder(order)
	p.Unlock()
}

func (p *SystemBase) GetOrder() (SystemPeriod, int) {
	p.Lock()
	defer p.Unlock()
	return SystemPeriod(p.order >> 32), int(p.order & 0xffff)
}

func (p *SystemBase) Clean(entity *Entity) {
	p.Lock()
	defer p.Unlock()

	p.pop.PushBack(entity)
}

type Start struct {
	SystemBase
	//private
	entityLock sync.RWMutex
	entityIndex map[*Entity]int
	data        [][]interface{}
}

func (p *Start) Init() {
	//inject system type info
	p.SetType(reflect.TypeOf(p))


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

func (p *Start) SystemUpdate() {
	p.PreUpdate()
	//TODO slice task queue task_length / k * cpu_num

	//TODO 检查 COMPONENT_STATE_NONE 状态下，是否执行，
}

func (p *Start) Filter(entity *Entity) {
	//check exist
	p.entityLock.RLock()
	if _, ok := p.entityIndex[entity]; ok {
		p.entityLock.RUnlock()
		return
	}
	p.entityLock.RUnlock()

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


