package ecs

import (
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
