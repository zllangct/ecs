package ecs

import (
	"reflect"
	"sync"
)

type opTask struct {
	target *EntityInfo
	com    IComponent
	op     CollectionOperate
	next   *opTask
}

type opTaskList struct {
	comType ComponentType
	head    *opTask
	tail    *opTask
}

var opTaskPool = newTaskPool()

type taskPool struct {
	pool sync.Pool
}

func newTaskPool() *taskPool {
	return &taskPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(opTask)
			},
		},
	}
}

func (p *taskPool) Get() *opTask {
	v := p.pool.Get()
	if v == nil {
		return &opTask{}
	}
	return v.(*opTask)
}

func (p *taskPool) Put(t *opTask) {
	p.pool.Put(t)
}

type ComponentCollection2 struct {
	collections map[reflect.Type]interface{}

	opLog map[reflect.Type]opTaskList
	once  map[reflect.Type]struct{}

	temp   chan *opTask
	toggle chan chan map[reflect.Type]opTaskList
}

func newComponentCollection2() *ComponentCollection2 {
	return &ComponentCollection2{
		collections: make(map[reflect.Type]interface{}),
		temp:        make(chan *opTask, 64),
		toggle:      make(chan chan map[reflect.Type]opTaskList),
		once:        map[reflect.Type]struct{}{},
		opLog:       map[reflect.Type]opTaskList{},
	}
}

func (c *ComponentCollection2) operate(op CollectionOperate, e *EntityInfo, com IComponent) {
	task := opTaskPool.Get()
	task.target = e
	task.com = com
	task.op = op
	c.temp <- task
}

func (c *ComponentCollection2) opLogCollector() {
	var tasks chan map[reflect.Type]opTaskList
	b := false
	for {
		select {
		case op := <-c.temp:
			c.addOpLog(op)
		case tasks = <-c.toggle:
			b = true
			break
		}
		if b {
			break
		}
	}
	tasks <- c.opLog
}

func (c *ComponentCollection2) getTempTasks() []func() {
	var tasks []func()
	ret := c.pauseAndGetTasks()
	for typ, list := range ret {
		tasks = append(tasks, func() {
			c.opExecute(list, c.collections[typ])
		})
	}
	return tasks
}

func (c *ComponentCollection2) pauseAndGetTasks() map[reflect.Type]opTaskList {
	ret := make(chan map[reflect.Type]opTaskList)
	c.toggle <- ret
	return <-ret
}

func (c *ComponentCollection2) run() {
	go c.opLogCollector()
}

func (c *ComponentCollection2) addOpLog(task *opTask) {
	typ := task.com.Type()
	if _, ok := c.opLog[typ]; ok {
		c.opLog[typ].tail.next = task
	} else {
		c.opLog[typ] = opTaskList{
			comType: task.com.getComponentType(),
			head:    task,
			tail:    task,
		}
	}

	ct := task.com.getComponentType()
	if task.op == CollectionOperateAdd &&
		(ct == ComponentTypeDisposable || ct == ComponentTypeFreeDisposable) {
		c.disposableTemp(task.com.Type())
	}
}

func (c *ComponentCollection2) opExecute(taskList opTaskList, collection any) reflect.Type {
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
				task.target.componentAdded(t, ret)
			case ComponentTypeFreeDisposable:
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
	for task := taskList.head; task != nil; task = task.next {
		opTaskPool.Put(task)
	}
	return t
}

func (c *ComponentCollection2) disposableTemp(typ reflect.Type) {
	if _, ok := c.once[typ]; !ok {
		c.once[typ] = struct{}{}
	}
}
func (c *ComponentCollection2) clearDisposable() {
	for typ, _ := range c.once {
		delete(c.collections, typ)
	}
}
