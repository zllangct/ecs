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
	operate(op CollectionOperate, entity Entity, component IComponent)
	deleteOperate(entity Entity, it uint16)
	getTempTasks() []func()
	clearDisposable()
	getComponentSet(typ reflect.Type) IComponentSet
	getComponentSetByIntType(typ uint16) IComponentSet
	getCollections() *SparseArray[uint16, IComponentSet]
	checkSet(com IComponent)
}

type ComponentCollection struct {
	collections *SparseArray[uint16, IComponentSet]
	world       *ecsWorld
	bucket      int64
	locks       []sync.RWMutex
	opLog       []map[reflect.Type]*opTaskList
	once        []map[reflect.Type]struct{}
}

func NewComponentCollection(world *ecsWorld, k int) *ComponentCollection {
	cc := &ComponentCollection{
		world:       world,
		collections: NewSparseArray[uint16, IComponentSet](),
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

func (c *ComponentCollection) operate(op CollectionOperate, entity Entity, component IComponent) {
	var hash int64
	switch component.getComponentType() {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
	case ComponentTypeNormal, ComponentTypeDisposable:
		hash = int64(entity) & c.bucket
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

func (c *ComponentCollection) deleteOperate(entity Entity, it uint16) {
	var hash int64
	meta := c.world.componentMeta.GetComponentMetaInfo(it)
	switch meta.componentType {
	case ComponentTypeFree, ComponentTypeFreeDisposable:
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
	case ComponentTypeNormal, ComponentTypeDisposable:
		hash = int64(entity) & c.bucket
	}

	typ := meta.typ
	newOpt := opTaskPool.Get()
	newOpt.target = entity
	newOpt.com = nil
	newOpt.op = CollectionOperateDelete

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
		e := com.Owner()
		hash = int64(e) & c.bucket
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
		taskList := list
		meta := c.world.GetComponentMetaInfo(typ)
		set := *(c.collections.Get(meta.it))
		if set == nil {
			c.checkSet(taskList.head.com)
			set = *(c.collections.Get(meta.it))
		}

		fn := func() {
			c.opExecute(taskList, set)
		}
		tasks = append(tasks, fn)
	}

	fn := func() {
		for typ, list := range combination {
			meta := c.world.GetComponentMetaInfo(typ)
			for task := list.head; task != nil; task = task.next {
				info, ok := c.world.GetEntityInfo(task.target)
				if ok {
					switch task.op {
					case CollectionOperateAdd:
						switch task.com.getComponentType() {
						case ComponentTypeNormal, ComponentTypeDisposable:
							info.addComponentInternal(meta.it)
						}
					case CollectionOperateDelete:
						switch task.com.getComponentType() {
						case ComponentTypeNormal, ComponentTypeDisposable:
							info.removeComponentInternal(meta.it)
						}
					}
				}
			}
		}
	}
	tasks = append(tasks, fn)
	return tasks
}

func (c *ComponentCollection) opExecute(taskList *opTaskList, collection IComponentSet) {
	meta := collection.GetElementMeta()
	for task := taskList.head; task != nil; task = task.next {
		switch task.op {
		case CollectionOperateAdd:
			task.com.setIntType(meta.it)
			task.com.setOwner(task.target)
			task.com.addToCollection(collection.pointer())
			switch task.com.getComponentType() {
			case ComponentTypeDisposable:
				c.disposableTemp(task.com, meta.typ)
			case ComponentTypeFreeDisposable:
				c.disposableTemp(task.com, meta.typ)
			}
		case CollectionOperateDelete:
			collection.Remove(task.target)
			//switch task.com.getComponentType() {
			//case ComponentTypeNormal, ComponentTypeDisposable:
			//}
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

func (c *ComponentCollection) getComponentSet(typ reflect.Type) IComponentSet {
	meta := c.world.GetComponentMetaInfo(typ)
	return *(c.collections.Get(meta.it))
}

func (c *ComponentCollection) getComponentSetByIntType(it uint16) IComponentSet {
	return *(c.collections.Get(it))
}

func (c *ComponentCollection) getCollections() *SparseArray[uint16, IComponentSet] {
	return c.collections
}

func (c *ComponentCollection) checkSet(com IComponent) {
	typ := com.Type()
	meta := c.world.GetComponentMetaInfo(typ)
	isExist := c.collections.Exist(meta.it)
	if !isExist {
		set := com.newCollection(meta)
		c.collections.Add(set.GetElementMeta().it, &set)
	}
}

func (c *ComponentCollection) removeAllByType(typ reflect.Type) {
	meta := c.world.GetComponentMetaInfo(typ)
	set := c.collections.Get(meta.it)
	if set == nil {
		return
	}
	(*set).Clear()
}
