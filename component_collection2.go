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

func (o *opTask) Reset() {
	o.target = nil
	o.com = nil
	o.op = CollectionOperateNone
	o.next = nil
}

type opTaskList struct {
	comType ComponentType
	len     int
	head    *opTask
	tail    *opTask
}

func (o *opTaskList) Len() int {
	return o.len
}

func (o *opTaskList) Combine(list *opTaskList) {
	if o.head == nil {
		o.head = list.head
		o.tail = list.tail
	} else {
		o.tail.next = list.head
		o.tail = list.tail
	}
	o.len += list.len
}

func (o *opTaskList) Append(task *opTask) {
	if o.head == nil {
		o.head = task
		o.tail = task
	} else {
		o.tail.next = task
		o.tail = task
	}
	o.len++
}

func (o *opTaskList) Reset() {
	o.len = 0
	o.head = nil
	o.tail = nil
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
	t.Reset()
	p.pool.Put(t)
}

type ComponentCollection2 struct {
	collections map[reflect.Type]interface{}

	opLog map[reflect.Type]*opTaskList
	once  map[reflect.Type]struct{}

	temp   chan *opTask
	toggle chan chan map[reflect.Type]*opTaskList
}

func newComponentCollection2() *ComponentCollection2 {
	c := &ComponentCollection2{
		collections: make(map[reflect.Type]interface{}),
		temp:        make(chan *opTask, 64),
		toggle:      make(chan chan map[reflect.Type]*opTaskList),
		once:        map[reflect.Type]struct{}{},
		opLog:       map[reflect.Type]*opTaskList{},
	}
	c.collectorRun()
	return c
}

func (c *ComponentCollection2) operate(op CollectionOperate, e *EntityInfo, com IComponent) {
	task := opTaskPool.Get()
	task.target = e
	task.com = com
	task.op = op
	c.temp <- task
}

func (c *ComponentCollection2) opLogCollector() {
	var tasks chan map[reflect.Type]*opTaskList
	b := false
	for {
		select {
		case op := <-c.temp:
			c.addOpLog(op)
		case tasks = <-c.toggle:
			// 保证temp管道中的OpLog全部被处理
			for len(c.temp) > 0 {
				op := <-c.temp
				c.addOpLog(op)
			}
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
		collection, ok := c.collections[typ]
		if !ok {
			c.collections[typ] = list.head.com.newCollection()
			collection = c.collections[typ]
		}
		taskList := list
		//Log.Infof("com:%s, collection:%s", typ.String(), reflect.TypeOf(collection).String())
		tasks = append(tasks, func() {
			c.opExecute(taskList, collection)
		})
	}
	return tasks
}

func (c *ComponentCollection2) pauseAndGetTasks() map[reflect.Type]*opTaskList {
	ret := make(chan map[reflect.Type]*opTaskList)
	c.toggle <- ret
	return <-ret
}

func (c *ComponentCollection2) collectorRun() {
	go c.opLogCollector()
}

func (c *ComponentCollection2) addOpLog(task *opTask) {
	typ := task.com.Type()
	tl, ok := c.opLog[typ]
	if !ok {
		tl = &opTaskList{}
		c.opLog[typ] = tl
	}

	tl.Append(task)

	ct := task.com.getComponentType()
	if task.op == CollectionOperateAdd &&
		(ct == ComponentTypeDisposable || ct == ComponentTypeFreeDisposable) {
		c.disposableTemp(task.com.Type())
	}
}

func (c *ComponentCollection2) opExecute(taskList *opTaskList, collection any) {
	//Log.Infof("opExecute com:%s, collection:%s", taskList.head.com.Type().String(), reflect.TypeOf(collection).String())
	for task := taskList.head; task != nil; task = task.next {
		t := task.com.Type()
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
	next := taskList.head
	for next != nil {
		task := next
		next = next.next
		opTaskPool.Put(task)
	}
	taskList.Reset()
}

func (c *ComponentCollection2) disposableTemp(typ reflect.Type) {
	if _, ok := c.once[typ]; !ok {
		c.once[typ] = struct{}{}
	}
}
func (c *ComponentCollection2) clearDisposable() {
	for typ, _ := range c.once {
		c.removeAllByType(typ)
	}
}

func (c *ComponentCollection2) getCollection(typ reflect.Type) interface{} {
	return c.collections[typ]
}

func (c *ComponentCollection2) getCollections() map[reflect.Type]interface{} {
	return c.collections
}

func (c *ComponentCollection2) removeAllByType(typ reflect.Type) {
	delete(c.collections, typ)
}
