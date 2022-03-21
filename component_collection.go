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

type IComponentCollection interface {
	operate(op CollectionOperate, entity *EntityInfo, component IComponent)
	getTempTasks() []func()
	collectorRun()
	clearDisposable()
	getCollection(typ reflect.Type) interface{}
	getCollections() map[reflect.Type]interface{}
}

type ComponentCollection struct {
	collections map[reflect.Type]interface{}
	bucket      int64
	locks       []sync.RWMutex
	opLog       []map[reflect.Type]*opTaskList
	once        []map[reflect.Type]struct{}
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
	cc.opLog = make([]map[reflect.Type]*opTaskList, cc.bucket+1)
	cc.initOptTemp()
	cc.once = make([]map[reflect.Type]struct{}, cc.bucket+1)
	cc.initOnce()

	return cc
}

func (c *ComponentCollection) initOptTemp() {
	for index := range c.opLog {
		c.locks[index].Lock()
		c.opLog[index] = make(map[reflect.Type]*opTaskList)
		c.locks[index].Unlock()
	}
}

func (c *ComponentCollection) initOnce() {
	for index := range c.once {
		c.locks[index].Lock()
		c.once[index] = make(map[reflect.Type]struct{})
		c.locks[index].Unlock()
	}
}

func (c *ComponentCollection) operate(op CollectionOperate, entity *EntityInfo, component IComponent) {
	var hash int64
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
	case ComponentTypeNormal, ComponentTypeDisposable:
		if entity == nil {
			Log.Errorf("invalid operate, entity is nil")
			return
		}
		hash = entity.hashKey() & c.bucket
	}

	typ := component.Type()
	newOpt := opTaskPool.Get()
	newOpt.target = entity
	newOpt.com = component
	newOpt.op = op

	b := c.opLog[hash]

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	tl, ok := b[typ]
	if !ok {
		tl = &opTaskList{}
		b[typ] = tl
	}

	tl.Append(newOpt)
}

func (c *ComponentCollection) clearDisposable() {
	for i := 0; i < len(c.once); i++ {
		c.locks[i].Lock()
		m := c.once[i]
		if len(m) > 0 {
			for typ, _ := range m {
				c.removeAllByType(typ)
				delete(m, typ)
				//Log.Info("collection Remove (Disposable):", typ.String())
			}
		}
		c.locks[i].Unlock()
	}
}

func (c *ComponentCollection) disposableTemp(com IComponent, typ reflect.Type) {
	var hash int64
	switch com.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
	case ComponentTypeNormal, ComponentTypeDisposable:
		info := com.Owner()
		if info == nil {
			Log.Errorf("invalid operate, entity is nil")
			return
		}
		hash = info.hashKey() & c.bucket
	}
	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	if _, ok := c.once[hash][typ]; !ok {
		c.once[hash][typ] = struct{}{}
	}
}

func (c *ComponentCollection) getTempTasks() []func() {
	combination := make(map[reflect.Type]*opTaskList)

	for i := 0; i < len(c.opLog); i++ {
		c.locks[i].RLock()
		for typ, list := range c.opLog[i] {
			if list.Len() == 0 {
				continue
			}
			if _, ok := combination[typ]; ok {
				combination[typ].Combine(list)
			} else {
				combination[typ] = list
			}
		}
		c.opLog[i] = make(map[reflect.Type]*opTaskList)
		c.locks[i].RUnlock()
	}

	var tasks []func()
	for typ, list := range combination {
		typTemp := typ
		taskList := list
		collection, ok := c.collections[typTemp]
		if !ok {
			c.collections[typTemp] = taskList.head.com.newCollection()
			collection = c.collections[typTemp]
		}

		fn := func() {
			c.opExecute(taskList, collection)
		}
		tasks = append(tasks, fn)
	}
	return tasks
}

func (c *ComponentCollection) opExecute(taskList *opTaskList, collection any) {
	var t reflect.Type
	for task := taskList.head; task != nil; task = task.next {
		t = task.com.Type()
		switch task.op {
		case CollectionOperateAdd:
			ret := task.com.addToCollection(collection)
			switch task.com.getComponentType() {
			case ComponentTypeNormal:
				task.target.componentAdded(t, ret)
			case ComponentTypeDisposable:
				c.disposableTemp(task.com, t)
				task.target.componentAdded(t, ret)
			case ComponentTypeFreeDisposable:
				c.disposableTemp(task.com, t)
			}
			task.com = ret
		case CollectionOperateDelete:
			task.com.deleteFromCollection(collection)
			switch task.com.getComponentType() {
			case ComponentTypeNormal, ComponentTypeDisposable:
				task.target.componentDeleted(t, task.com.getComponentType())
			}
		}
	}
	next := taskList.head
	for next != nil {
		task := next
		next = next.next
		opTaskPool.Put(task)
	}
	taskList.Reset()
}

func (c *ComponentCollection) getCollection(typ reflect.Type) interface{} {
	return c.collections[typ]
}

func (c *ComponentCollection) getCollections() map[reflect.Type]interface{} {
	return c.collections
}

func (c *ComponentCollection) removeAllByType(typ reflect.Type) {
	delete(c.collections, typ)
}

func (c *ComponentCollection) collectorRun() {
}
