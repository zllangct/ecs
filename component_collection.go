package ecs

import (
	"reflect"
	"sync"
	"unsafe"
)

type CollectionOperate uint8

const (
	CollectionOperateNone   CollectionOperate = iota
	CollectionOperateAdd                      //add component operation
	CollectionOperateDelete                   //delete component operation
)

type OperateInfo struct {
	target *EntityInfo
	com    IComponent
	op     CollectionOperate
}

func NewTemplateOperateInfo(entity *EntityInfo, template IComponent, op CollectionOperate) OperateInfo {
	return OperateInfo{target: entity, com: template, op: op}
}

type ComponentCollection struct {
	collections map[reflect.Type]interface{}
	//new component cache
	bucket  int64
	locks   []sync.RWMutex
	optTemp []map[reflect.Type][]OperateInfo
	componentsNew map[reflect.Type][]OperateInfo
	once    []map[reflect.Type]struct{}
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collections: map[reflect.Type]interface{}{},
	}

	for i := 1; ; i++ {
		if c := int64(1 << i); int64(k) < c {
			cc.bucket = c - 1
			break
		}
	}

	cc.locks = make([]sync.RWMutex, cc.bucket+1)
	for i := int64(0); i < cc.bucket+1; i++ {
		cc.locks[i] = sync.RWMutex{}
	}
	cc.optTemp = make([]map[reflect.Type][]OperateInfo, cc.bucket+1)
	cc.resetOptTemp()
	cc.once = make([]map[reflect.Type]struct{}, cc.bucket+1)
	cc.resetOnce()

	cc.componentsNew = make(map[reflect.Type][]OperateInfo)

	return cc
}

func (c *ComponentCollection) resetOptTemp() {
	for index := range c.optTemp {
		c.locks[index].Lock()
		c.optTemp[index] = make(map[reflect.Type][]OperateInfo)
		c.locks[index].Unlock()
	}
}

func (c *ComponentCollection) resetOnce() {
	for index := range c.once {
		c.locks[index].Lock()
		c.once[index] = make(map[reflect.Type]struct{})
		c.locks[index].Unlock()
	}
}

func (c *ComponentCollection) TempTemplateOperate(entity *EntityInfo, component IComponent, op CollectionOperate) {
	var hash int64
	var isOnce bool
	switch component.getComponentType() {
	case ComponentTypeFree:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
		isOnce = false
	case ComponentTypeFreeDisposable:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
		isOnce = true
	case ComponentTypeNormal:
		if entity == nil {
			Log.Errorf("invalid operate, entity is nil")
			return
		}
		hash = entity.hashKey() & c.bucket
		isOnce = false
	case ComponentTypeDisposable:
		if entity == nil {
			Log.Errorf("invalid operate, entity is nil")
			return
		}
		hash = entity.hashKey() & c.bucket
		isOnce = true
	}

	typ := component.Type()
	newOpt := NewTemplateOperateInfo(entity, component, op)

	b := c.optTemp[hash]
	o := c.once[hash]

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	if _, ok := b[typ]; ok {
		b[typ] = append(b[typ], newOpt)
	} else {
		b[typ] = []OperateInfo{newOpt}
	}

	if isOnce {
		if _, ok := o[typ]; !ok {
			o[typ] = Empty
		}
	}
}

func (c *ComponentCollection) ClearDisposable() {
	for index := range c.once {
		c.locks[index].Lock()
		m :=c.once[index]
		for typ, _ := range m {
			c.RemoveAllByType(typ)
		}
		c.locks[index].Unlock()
	}
}

func (c *ComponentCollection) GetTempTasks() []func() (reflect.Type, []OperateInfo) {
	combination := make(map[reflect.Type][]OperateInfo)

	for i := 0; i < len(c.optTemp); i++ {
		c.locks[i].RLock()
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
		c.locks[i].RUnlock()
	}

	var tasks []func() (reflect.Type, []OperateInfo)
	for typ, opList := range combination {
		typTemp := typ
		oopList := opList
		collection, ok := c.collections[typTemp]
		if !ok {
			c.collections[typTemp] = oopList[0].com.newCollection()
			collection = c.collections[typTemp]
		}

		fn := func() (reflect.Type, []OperateInfo) {
			n := make([]OperateInfo, len(oopList))
			var t reflect.Type
			for _, operate := range oopList {
				t = operate.com.Type()
				switch operate.op {
				case CollectionOperateAdd:
					ret := operate.com.addToCollection(collection)
					switch operate.com.getComponentType() {
					case ComponentTypeNormal, ComponentTypeDisposable:
						operate.target.componentAdded(t, ret)
					}
					operate.com = ret
					n = append(n, operate)
				case CollectionOperateDelete:
					operate.com.deleteFromCollection(collection)
					switch operate.com.getComponentType() {
					case ComponentTypeNormal, ComponentTypeDisposable:
						operate.target.componentDeleted(t, operate.com)
					}
					n = append(n, operate)
				}
			}
			return t, n
		}
		tasks = append(tasks, fn)
	}
	return tasks
}

func (c *ComponentCollection) TempTasksDone(newList map[reflect.Type][]OperateInfo) {
	c.componentsNew = newList
	c.resetOptTemp()
}

func (c *ComponentCollection) GetNewComponentsAll() map[reflect.Type][]OperateInfo {
	return c.componentsNew
}

func (c *ComponentCollection) GetNewComponents(typ reflect.Type) []OperateInfo {
	return c.componentsNew[typ]
}

func (c *ComponentCollection) GetCollection(typ reflect.Type) interface{} {
	return c.collections[typ]
}

func (c *ComponentCollection) RemoveAllByType(typ reflect.Type) {
	delete(c.collections,typ)
}
