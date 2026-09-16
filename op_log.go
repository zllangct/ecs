package ecs

import (
	"sync"
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

// OpLog 组件操作日志。
// 写入路径分双轨：
//   - system 执行期间（经 SystemContext.AddComponents）写入所属 system 的独立队列，
//     单写者、无锁；外层 map 在注册期建好，执行期只读；
//   - 主线程/EntityInfo.Add 路径写入默认队列，Mutex 保护。
//
// flush（getOpTasks）只在帧同步点由主线程调用，与 system 执行互不重叠。
type OpLog struct {
	world *World

	mu   sync.Mutex
	main map[ComponentIntType]*opTaskList
	sys  map[uint64]map[ComponentIntType]*opTaskList
}

func NewOpLog(world *World, k int) *OpLog {
	return &OpLog{
		world: world,
		main:  make(map[ComponentIntType]*opTaskList),
		sys:   make(map[uint64]map[ComponentIntType]*opTaskList),
	}
}

// registerSystem 为 system 分配独立操作队列，注册期主线程调用
func (c *OpLog) registerSystem(id uint64) {
	if _, ok := c.sys[id]; !ok {
		c.sys[id] = make(map[ComponentIntType]*opTaskList)
	}
}

// operate 默认队列入口（主线程 / EntityInfo.Add）
func (c *OpLog) operate(op Operate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	appendOp(c.main, op)
}

// operateForSystem system 执行期间入口，仅所属 system 的 goroutine 写入，无锁
func (c *OpLog) operateForSystem(id uint64, op Operate) {
	if q, ok := c.sys[id]; ok {
		appendOp(q, op)
		return
	}
	c.operate(op)
}

func appendOp(b map[ComponentIntType]*opTaskList, op Operate) {
	typ := GetIntTypeByComp(op.Comp)
	newOpt := opTaskPool.Get()
	newOpt.target = op.Entity
	newOpt.com = op.Comp
	newOpt.op = op.Op

	tl, ok := b[typ]
	if !ok {
		tl = &opTaskList{}
		b[typ] = tl
	}
	tl.Append(newOpt)
}

// getOpTasks 返回两阶段任务：
//   - 第一相：组件集合增删 + compound 维护（可并发）；
//   - 第二相：组表成员资格 sync（依赖第一相完成后的 CSet/compound 终态，必须滞后执行）。
func (c *OpLog) getOpTasks() ([]func(), []func(), func()) {
	combination := make(map[ComponentIntType]*opTaskList)

	merge := func(src map[ComponentIntType]*opTaskList) {
		for it, list := range src {
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
	}

	c.mu.Lock()
	merge(c.main)
	c.mu.Unlock()
	// 帧同步点调用，system 均未运行，sys 队列无并发写
	for _, q := range c.sys {
		merge(q)
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
			c.world.compVersion++
			if taskList.head.com.IsDisposable() {
				c.world.disposableTypes = append(c.world.disposableTypes, intType)
			}
			if taskList.head.com.IsNomadic() {
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

	// 第二相：组表成员资格 sync。收集本帧涉及组内类型的实体，按组聚合。
	touched := map[*Archetype]map[EntityIndex]struct{}{}
	for it, list := range combination {
		a, ok := c.world.archetypes.owner(it)
		if !ok {
			continue
		}
		if touched[a] == nil {
			touched[a] = map[EntityIndex]struct{}{}
		}
		for task := list.head; task != nil; task = task.next {
			touched[a][task.target.Index()] = struct{}{}
		}
	}
	var syncTasks []func()
	for a, entities := range touched {
		a, entities := a, entities
		syncTasks = append(syncTasks, func() {
			c.world.syncGroup(a, entities)
		})
	}

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
	return tasks, syncTasks, clean
}

func (c *OpLog) opExecute(taskList *opTaskList, set ComponentSet) {
	for task := taskList.head; task != nil; task = task.next {
		switch task.op {
		case ComponentOperateAdd:
			// 实体已在本次帧同步点销毁时丢弃 Add（destroy 先于 op flush 执行），
			// 否则会在 CSet 留下孤儿数据，index 复用后僵尸复活。
			// nomadic 组件无实体属主（target 为零值），豁免校验。
			if !task.com.IsNomadic() {
				if _, ok := c.world.getEntityInfo(task.target); !ok {
					continue
				}
			}
			set.Add(task.target, task.com)
		case ComponentOperateDelete:
			set.Remove(task.target)
		case ComponentOperateDeleteAll:
			set.Reset()
		}
	}
}
