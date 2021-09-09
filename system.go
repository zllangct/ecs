package ecs

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
	Init(runtime *Runtime) //init
	GetBase() *SystemBase  //get system base data
	GetType() reflect.Type
	GetOrder() Order
	GetRequirements() map[reflect.Type]struct{}
	Call(label int) interface{}
}

type SystemBase struct {
	sync.Mutex
	requirements map[reflect.Type]struct{}
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

func (p *SystemBase) SetRequirements(rqs ...IComponentType) {
	if p.requirements == nil {
		p.requirements = map[reflect.Type]struct{}{}
	}
	for _, value := range rqs {
		p.requirements[reflect.TypeOf(value)] = struct{}{}
	}
}

func (p *SystemBase) GetRequirements() map[reflect.Type]struct{} {
	return p.requirements
}

func (p *SystemBase) Init(runtime *Runtime) {
	p.requirements = map[reflect.Type]struct{}{}
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
	defer p.Unlock()

	p.order = order
}

func (p *SystemBase) GetOrder() Order {
	p.Lock()
	defer p.Unlock()

	return p.order
}

func (p *SystemBase) IsConcerned(com IComponent) bool {
	cType := com.GetType()
	if _, concerned := p.requirements[cType]; concerned {
		for r, _ := range p.requirements {
			if r != cType {
				if !com.GetOwner().Has(r) {
					concerned = false
					break
				}
			}
		}
		return concerned
	}
	return false
}

func (p *SystemBase) GetRuntime() *Runtime {
	return p.runtime
}

func (p *SystemBase) GetNewComponent(op CollectionOperate) map[reflect.Type][]CollectionOperateInfo {
	temp := map[reflect.Type][]CollectionOperateInfo{}
	for typ, _ := range p.requirements {
		temp[typ] = p.runtime.getNewComponents(op, typ)
	}
	return temp
}
