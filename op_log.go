package ecs

import (
	"cmp"
	"slices"
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

// Each writer retains its type buckets across frames. Only nonempty buckets
// are linked, so a flush never visits a type merely because it was used before.
type opBucket struct {
	opTaskList
	typ  ComponentIntType
	next *opBucket
}

type opQueue struct {
	systemID uint64
	buckets  map[ComponentIntType]*opBucket
	active   *opBucket
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

	mu       sync.Mutex
	main     opQueue
	sys      map[uint64]*opQueue
	sysOrder []*opQueue
}

func NewOpLog(world *World, k int) *OpLog {
	return &OpLog{
		world: world,
		main:  opQueue{buckets: make(map[ComponentIntType]*opBucket)},
		sys:   make(map[uint64]*opQueue),
	}
}

// registerSystem 为 system 分配独立操作队列，注册期主线程调用，并记录固定合并顺序
func (c *OpLog) registerSystem(id uint64) {
	if _, ok := c.sys[id]; !ok {
		q := &opQueue{systemID: id, buckets: make(map[ComponentIntType]*opBucket)}
		c.sys[id] = q
		c.sysOrder = append(c.sysOrder, q)
	}
}

// HasPending 是否仍有尚未 flush 到组件存储的操作（Marshal 只读导出前校验用）
func (c *OpLog) HasPending() bool {
	c.mu.Lock()
	pending := c.main.active != nil
	c.mu.Unlock()
	if pending {
		return true
	}
	for _, q := range c.sysOrder {
		if q.active != nil {
			return true
		}
	}
	return false
}

// operate 默认队列入口（主线程 / EntityInfo.Add）
func (c *OpLog) operate(op Operate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.main.append(op)
}

// operateForSystem system 执行期间入口，仅所属 system 的 goroutine 写入，无锁
func (c *OpLog) operateForSystem(id uint64, op Operate) {
	if q, ok := c.sys[id]; ok {
		q.append(op)
		return
	}
	c.operate(op)
}

func (q *opQueue) append(op Operate) {
	typ := GetIntTypeByComp(op.Comp)
	newOpt := opTaskPool.Get()
	newOpt.target = op.Entity
	newOpt.com = op.Comp
	newOpt.op = op.Op

	bucket := q.buckets[typ]
	if bucket == nil {
		bucket = &opBucket{typ: typ}
		q.buckets[typ] = bucket
	}
	if bucket.head == nil {
		bucket.next = q.active
		q.active = bucket
	}
	bucket.Append(newOpt)
}

// getOpTasks 返回两阶段任务：
//   - 第一相：组件集合增删 + compound 维护（可并发）；
//   - 第二相：组表成员资格 sync（依赖第一相完成后的 CSet/compound 终态，必须滞后执行）。
func (c *OpLog) getOpTasks() ([]func(), []func(), func()) {
	var combination map[ComponentIntType]*opTaskList

	merge := func(src *opQueue) {
		for src.active != nil {
			bucket := src.active
			src.active = bucket.next
			bucket.next = nil
			if dst := combination[bucket.typ]; dst != nil {
				dst.Combine(&bucket.opTaskList)
			} else {
				if combination == nil {
					combination = make(map[ComponentIntType]*opTaskList)
				}
				combination[bucket.typ] = bucket.Clone()
			}
			bucket.Reset()
		}
	}

	c.mu.Lock()
	merge(&c.main)
	c.mu.Unlock()
	// 帧同步点调用，system 均未运行，按系统注册顺序确定性合并 sys 队列
	for _, q := range c.sysOrder {
		if q.active != nil {
			merge(q)
		}
	}
	if len(combination) == 0 {
		return nil, nil, func() {}
	}

	// 按 ComponentIntType 升序排列，保证任务构造成果与 compVersion 递增确定性
	orderedTypes := make([]ComponentIntType, 0, len(combination))
	for intType, list := range combination {
		if list.Len() > 0 {
			orderedTypes = append(orderedTypes, intType)
		}
	}
	slices.Sort(orderedTypes)

	tasks := make([]func(), 0, len(orderedTypes)+1)
	for _, intType := range orderedTypes {
		taskList := combination[intType]
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
		for _, it := range orderedTypes {
			list := combination[it]
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

	// 第二相：组表成员资格 sync。收集本帧涉及组内类型的实体，按组聚合（按组首类型升序确定性排列）。
	var touchedGroups []*Archetype
	for _, it := range orderedTypes {
		list := combination[it]
		a, ok := c.world.archetypes.owner(it)
		if !ok {
			continue
		}
		if len(a.pendingRows) == 0 {
			touchedGroups = append(touchedGroups, a)
		}
		for task := list.head; task != nil; task = task.next {
			index := task.target.Index()
			if n := len(a.pendingRows); n == 0 || a.pendingRows[n-1] != index {
				a.pendingRows = append(a.pendingRows, index)
			}
		}
	}
	slices.SortFunc(touchedGroups, func(left, right *Archetype) int {
		return cmp.Compare(left.types[0], right.types[0])
	})
	syncTasks := make([]func(), 0, len(touchedGroups))
	for _, a := range touchedGroups {
		syncTasks = append(syncTasks, func() {
			slices.Sort(a.pendingRows)
			c.world.syncGroupOrdered(a, slices.Compact(a.pendingRows))
		})
	}

	clean := func() {
		for _, a := range touchedGroups {
			a.pendingRows = a.pendingRows[:0]
		}
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
			// 校验 reuse 代数，避免旧句柄误删新实体组件
			if !task.com.IsNomadic() {
				if _, ok := c.world.getEntityInfo(task.target); !ok {
					continue
				}
			}
			set.Remove(task.target)
		case ComponentOperateDeleteAll:
			set.Reset()
		}
	}
}
