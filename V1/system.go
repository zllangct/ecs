package main

import (
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
	Init()                           //init
	GetBase() *SystemBase            //get system base data
	GetType() reflect.Type
	GetOrder() (SystemPeriod, Order)
	GetRequirements() []reflect.Type
	Filter()       //interest filter of component
	SystemUpdate() //update every frame
	Call(label int) interface{}
}

type SystemBase struct {
	sync.Mutex
	requirements []reflect.Type
	order        SystemOrder
	runtime      *Runtime
	typ          reflect.Type
}

func (p *SystemBase) Call(label int) interface{} {
	return nil
}

func (p *SystemBase) GetBase() *SystemBase {
	return p
}

func (p *SystemBase) SetRequirements(rqs ...IComponent) {
	for _, value := range rqs {
		p.requirements = append(p.requirements, reflect.TypeOf(value))
	}
}

func (p *SystemBase) GetRequirements() []reflect.Type {
	return p.requirements
}

func (p *SystemBase) Init() {
	p.requirements = []reflect.Type{}
	p.SetOrder(ORDER_DEFAULT, PERIOD_DEFAULT)
}

func (p *SystemBase) SetType(typ reflect.Type) {
	p.Lock()
	defer p.Unlock()

	p.typ = typ
}

func (p *SystemBase) GetType() reflect.Type {
	p.Lock()
	defer p.Unlock()

	return reflect.TypeOf(p)
}

func (p *SystemBase) SetOrder(order Order, period ...SystemPeriod) {
	mPeriod := PERIOD_DEFAULT
	if len(period) > 0 {
		mPeriod = period[0]
	}
	p.Lock()
	p.order = SystemOrder(mPeriod)<<32 + SystemOrder(order)
	p.Unlock()
}

func (p *SystemBase) GetOrder() (SystemPeriod, Order) {
	p.Lock()
	defer p.Unlock()
	return SystemPeriod(p.order >> 32), Order(p.order & 0xffff)
}
