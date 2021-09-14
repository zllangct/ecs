package ecs

import (
	"reflect"
	"sync"
)

type CollectionOperate uint8

const (
	COLLECTION_OPERATE_NONE   CollectionOperate = iota
	COLLECTION_OPERATE_ADD                      //add component operation
	COLLECTION_OPERATE_DELETE                   //delete component operation
)

type ComponentOperate = CollectionOperate

const (
	COMPONENT_OPERATE_ADD    = COLLECTION_OPERATE_ADD    //add component operation
	COMPONENT_OPERATE_DELETE = COLLECTION_OPERATE_DELETE //delete component operation
)

type TemplateOperateInfo struct {
	target *Entity
	template    IComponentTemplate
	op     CollectionOperate
	typ reflect.Type
}

func NewTemplateOperateInfo(entity *Entity, template IComponentTemplate, typ reflect.Type, op CollectionOperate) TemplateOperateInfo {
	return TemplateOperateInfo{target: entity, template: template, op: op, typ: typ}
}

type ComponentOptResult struct {
	com IComponent
	opInfo TemplateOperateInfo
}

type ComponentCollection struct {
	collections map[reflect.Type]interface{}
	//new component cache
	locks []sync.Mutex
	base  int64
	optTemp []map[reflect.Type][]TemplateOperateInfo
	componentsNew map[reflect.Type][]ComponentOptResult
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collections: map[reflect.Type]interface{}{},
	}

	for i := 1; ; i++ {
		if c := int64(1 << i); int64(k) < c {
			cc.base = c - 1
			break
		}
	}

	cc.locks = make([]sync.Mutex, cc.base+1)
	for i:= int64(0); i < cc.base + 1; i++ {
		cc.locks[i] = sync.Mutex{}
	}
	cc.optTemp =  make([]map[reflect.Type][]TemplateOperateInfo, cc.base+1)
	cc.resetOptTemp()

	cc.componentsNew =  make(map[reflect.Type][]ComponentOptResult)
	return cc
}

func (c *ComponentCollection) resetOptTemp() {
	for index := range c.optTemp {
		c.optTemp[index] = make(map[reflect.Type][]TemplateOperateInfo)
	}
}

func (c *ComponentCollection) TempTemplateOperate(entity *Entity, template IComponentTemplate, op CollectionOperate) {
	hash := entity.ID() & c.base

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	typ := template.ComponentType()
	newOpt := NewTemplateOperateInfo(entity, template, typ, op)
	b := c.optTemp[hash]
	if _, ok := b[typ]; ok {
		b[typ] = append(b[typ], newOpt)
	} else {
		b[typ] = []TemplateOperateInfo{ newOpt }
	}
}

func (c *ComponentCollection) GetTempTasks() []func()(reflect.Type, []ComponentOptResult) {
	combination := make(map[reflect.Type][]TemplateOperateInfo)

	for i:=0; i < len(c.optTemp); i++ {
		for typ, op := range c.optTemp[i] {
			if len(op) == 0 {
				continue
			}
			if _, ok := combination[typ]; ok {
				combination[typ] = append(combination[typ], op...)
			} else {
				combination[typ] = op
			}
		}
	}

	var tasks []func() (reflect.Type, []ComponentOptResult)
	for typ, opList := range combination {
		typTemp := typ
		oopList := opList
		collection, ok := c.collections[typTemp]
		if !ok {
			c.collections[typTemp] = oopList[0].template.NewCollection()
			collection = c.collections[typTemp]
		}

		fn := func() (reflect.Type, []ComponentOptResult) {
			n := make([]ComponentOptResult, len(oopList))
			var t reflect.Type
			for _, operate := range oopList {
				t = operate.typ
				//add to component container
				ret := operate.template.AddToCollection(collection)
				//add to entity
				operate.target.componentAdded(t, ret)

				n = append(n, ComponentOptResult{com: ret, opInfo: operate})
			}
			return t, n
		}
		tasks = append(tasks, fn)
	}
	return tasks
}

func (c *ComponentCollection) TempTasksDone(newList map[reflect.Type][]ComponentOptResult) {
	c.componentsNew = newList
	c.resetOptTemp()
}

func (c *ComponentCollection) GetNewComponentsAll() map[reflect.Type][]ComponentOptResult {
	return c.componentsNew
}

func (c *ComponentCollection) GetNewComponents(typ reflect.Type) []ComponentOptResult {
	return c.componentsNew[typ]
}

func (c *ComponentCollection) GetCollection(typ reflect.Type) interface{} {
	return c.collections[typ]
}
