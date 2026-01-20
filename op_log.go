package ecs

import (
	"sync"
	"unsafe"
)

type ComponentOperate uint8

const (
	ComponentOperateNone      ComponentOperate = iota
	ComponentOperateAdd                        //add component operation
	ComponentOperateDelete                     //delete component operation
	ComponentOperateDeleteAll                  //delete component by type operation
)

type Operate struct {
	Entity Entity
	Comp   Component
	Op     ComponentOperate
}

type OpLog struct {
	world  *world
	bucket int64
	locks  []sync.RWMutex
	log    []map[ComponentIntType]*opTaskList
}

func NewOpLog(world *world, k int) *OpLog {
	cc := &OpLog{
		world: world,
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
	cc.log = make([]map[ComponentIntType]*opTaskList, cc.bucket+1)
	cc.initOptTemp()

	return cc
}

func (c *OpLog) initOptTemp() {
	for index := range c.log {
		c.locks[index].Lock()
		c.log[index] = make(map[ComponentIntType]*opTaskList)
		c.locks[index].Unlock()
	}
}

func (c *OpLog) operate(op Operate) {
	var hash int64
	hash = int64((uintptr)(unsafe.Pointer(&hash))) & c.bucket

	typ := GetIntTypeByComp(op.Comp)
	newOpt := opTaskPool.Get()
	newOpt.target = op.Entity
	newOpt.com = op.Comp
	newOpt.op = op.Op

	b := c.log[hash]

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	tl, ok := b[typ]
	if !ok {
		tl = &opTaskList{}
		b[typ] = tl
	}

	tl.Append(newOpt)
}

func (c *OpLog) getOpTasks() ([]func(), func()) {
	combination := make(map[ComponentIntType]*opTaskList)

	for i := 0; i < len(c.log); i++ {
		c.locks[i].RLock()
		for it, list := range c.log[i] {
			if list.Len() == 0 {
				continue
			}
			if _, ok := combination[it]; ok {
				combination[it].Combine(list)
			} else {
				combination[it] = list.Clone()
			}
			list.Reset()
		}

		c.locks[i].RUnlock()
	}

	var tasks []func()
	for intType, list := range combination {
		taskList := list
		if taskList.Len() == 0 {
			continue
		}
		set := c.world.components[intType]
		if set == nil {
			set = taskList.head.com.NewComponentSet()
			c.world.components[intType] = set
			if taskList.head.com.IsDisposable() {
				c.world.disposableTypes = append(c.world.disposableTypes, intType)
			}
			if taskList.tail.com.IsNomadic() {
				c.world.nomadicTypes = append(c.world.nomadicTypes, intType)
			}
		}

		fn := func() {
			c.opExecute(taskList, set)
		}
		tasks = append(tasks, fn)
	}

	fn := func() {
		for it, list := range combination {
			for task := list.head; task != nil; task = task.next {
				if task.op == ComponentOperateDelete {
					continue
				}
				info, ok := c.world.getEntityInfo(task.target)
				if ok {
					switch task.op {
					case ComponentOperateAdd:
						if !task.com.IsNomadic() {
							info.compound.Add(it)
						}
					case ComponentOperateDelete:
						if !task.com.IsNomadic() {
							info.compound.Remove(it)
						}
					}
				}
			}
		}
	}
	tasks = append(tasks, fn)

	clean := func() {
		for _, list := range combination {
			next := list.head
			for next != nil {
				task := next
				next = next.next
				opTaskPool.Put(task)
			}
			list.Reset()
		}
	}
	return tasks, clean
}

func (c *OpLog) opExecute(taskList *opTaskList, set ComponentSet) {
	for task := taskList.head; task != nil; task = task.next {
		switch task.op {
		case ComponentOperateAdd:
			set.Add(task.target, task.com)
		case ComponentOperateDelete:
			set.Remove(task.target)
		case ComponentOperateDeleteAll:
			set.Reset()
		}
	}
}
