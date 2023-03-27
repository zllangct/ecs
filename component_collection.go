package ecs

import (
	"reflect"
	"sync"
	"unsafe"
)

type CollectionOperate uint8

const (
	CollectionOperateNone      CollectionOperate = iota
	CollectionOperateAdd                         //add component operation
	CollectionOperateDelete                      //delete component operation
	CollectionOperateDeleteAll                   //delete component by type operation
)

type IComponentCollection interface {
	operate(op CollectionOperate, entity Entity, component IComponent)
	deleteOperate(op CollectionOperate, entity Entity, it uint16)
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

	cc.bucket = 0

	cc.locks = make([]sync.RWMutex, cc.bucket+1)
	for i := int64(0); i < cc.bucket+1; i++ {
		cc.locks[i] = sync.RWMutex{}
	}
	cc.opLog = make([]map[reflect.Type]*opTaskList, cc.bucket+1)
	cc.initOptTemp()

	return cc
}

func (c *ComponentCollection) initOptTemp() {
	for index := range c.opLog {
		c.locks[index].Lock()
		c.opLog[index] = make(map[reflect.Type]*opTaskList)
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

func (c *ComponentCollection) deleteOperate(op CollectionOperate, entity Entity, it uint16) {
	var hash int64
	meta := c.world.componentMeta.GetComponentMetaInfoByIntType(it)
	if meta.componentType&ComponentTypeFreeMask > 0 {
		hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket
	} else {
		hash = int64(entity) & c.bucket
	}

	typ := meta.typ
	newOpt := opTaskPool.Get()
	newOpt.target = entity
	newOpt.com = nil
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
	disposable := c.world.componentMeta.GetDisposableTypes()
	for r, _ := range disposable {
		meta := c.world.getComponentMetaInfoByType(r)
		if meta.componentType&ComponentTypeFreeMask > 0 {
			continue
		}
		set := c.collections.Get(meta.it)
		if set == nil {
			return
		}
		(*set).Range(func(com IComponent) bool {
			info, ok := c.world.entities.GetEntityInfo(com.Owner())
			if ok {
				info.removeFromCompound(meta.it)
			}
			return true
		})

		(*set).Clear()
	}
}

func (c *ComponentCollection) clearFree() {
	free := c.world.componentMeta.GetFreeTypes()
	for r, _ := range free {
		meta := c.world.getComponentMetaInfoByType(r)
		set := c.collections.Get(meta.it)
		if set == nil {
			return
		}
		(*set).Clear()
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
				combination[typ] = list.Clone()
			}
			list.Reset()
		}

		c.locks[i].RUnlock()
	}

	var tasks []func()
	for typ, list := range combination {
		taskList := list
		if taskList.Len() == 0 {
			continue
		}
		meta := c.world.getComponentMetaInfoByType(typ)
		setp := c.collections.Get(meta.it)
		if setp == nil {
			c.checkSet(taskList.head.com)
			setp = c.collections.Get(meta.it)
		}

		fn := func() {
			c.opExecute(taskList, *setp)
		}
		tasks = append(tasks, fn)
	}

	fn := func() {
		for typ, list := range combination {
			meta := c.world.getComponentMetaInfoByType(typ)
			for task := list.head; task != nil; task = task.next {
				if task.op == CollectionOperateDelete {
					continue
				}
				info, ok := c.world.getEntityInfo(task.target)
				if ok {
					switch task.op {
					case CollectionOperateAdd:
						switch task.com.getComponentType() {
						case ComponentTypeNormal, ComponentTypeDisposable:
							info.addToCompound(meta.it)
						}
					case CollectionOperateDelete:
						switch task.com.getComponentType() {
						case ComponentTypeNormal, ComponentTypeDisposable:
							info.removeFromCompound(meta.it)
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
			task.com.addToCollection(task.com.getComponentType(), collection.pointer())
		case CollectionOperateDelete:
			if meta.componentType&ComponentTypeFreeMask == 0 {
				collection.Remove(task.target)
			}
		case CollectionOperateDeleteAll:
			collection.Clear()
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
	meta := c.world.getComponentMetaInfoByType(typ)
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
	meta := c.world.getComponentMetaInfoByType(typ)
	isExist := c.collections.Exist(meta.it)
	if !isExist {
		set := com.newCollection(meta)
		c.collections.Add(set.GetElementMeta().it, &set)
	}
}
