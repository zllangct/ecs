package ecs

import (
	"errors"
	"reflect"
	"sync"
	"time"
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
	Init(runtime *Runtime) //init
	GetBase() *SystemBase  //get system base data
	GetType() reflect.Type
	GetOrder() Order
	GetRequirements() []reflect.Type
	Filter(component IComponent, op CollectionOperate) //interest filter of component
	SystemUpdate(delta time.Duration)                  //update every frame
	Call(label int) interface{}
}

type SystemBase struct {
	sync.Mutex
	requirements []reflect.Type
	order        Order
	runtime      *Runtime
	typ          reflect.Type
	isPreFilter  bool
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

func (p *SystemBase) Init(runtime *Runtime) {
	p.requirements = []reflect.Type{}
	p.SetOrder(ORDER_DEFAULT)
	p.runtime = runtime
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

func (p *SystemBase) SetOrder(order Order) {
	p.Lock()
	p.order = order
	p.Unlock()
}

func (p *SystemBase) GetOrder() Order {
	p.Lock()
	defer p.Unlock()
	return p.order
}

func (p *SystemBase) IsConcerned(com IComponent) bool {
	concerned := true
	ctyp := reflect.TypeOf(com)
	for _, typ := range p.requirements {
		if typ != ctyp && !com.GetOwner().Has(typ) {
			concerned = false
			break
		}
	}
	return concerned
}
