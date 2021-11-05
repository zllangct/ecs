package ecs

import (
	"reflect"
	"sync"
	"time"
)

// SystemOrder system execute Order, 32bit + 32bit: period + suborder
type SystemOrder uint64

// SystemPeriod system execute period:start->pre_update->update->pre_destroy->destroy
type SystemPeriod uint32

const (
	PeriodPreStart SystemPeriod = iota
	PeriodCustomEvent
	PeriodStart
	PeriodPostStart
	PeriodPreUpdate
	PeriodUpdate
	PeriodPostUpdate
	PeriodPerDestroy
	PeriodDestroy
	PeriodPostDestroy
)

// Order default suborder of system
type Order int32

const (
	OrderFront   Order = -1
	OrderAppend  Order = 999999
	OrderDefault Order = OrderAppend
)

type TempTask struct {
	lock *sync.Mutex
	m    map[reflect.Type][]OperateInfo
	wg   *sync.WaitGroup
	fn   func() (reflect.Type, []OperateInfo)
}

// OrderSequence extension of system group slice
type OrderSequence []*SystemGroup

//system execute flow
type systemFlow struct {
	sync.Mutex
	world        *World
	systemPeriod map[SystemPeriod]OrderSequence
	periodList   []SystemPeriod
	wg           *sync.WaitGroup
}

func newSystemFlow(runtime *World) *systemFlow {
	sf := &systemFlow{
		world: runtime,
		wg:    &sync.WaitGroup{},
	}
	sf.init()
	return sf
}

//initialize the system flow
func (p *systemFlow) init() {
	p.periodList = []SystemPeriod{
		PeriodCustomEvent,
		PeriodStart,
		PeriodPostStart,
		PeriodUpdate,
		PeriodPostUpdate,
		PeriodDestroy,
		PeriodPostDestroy,
	}
	p.systemPeriod = make(map[SystemPeriod]OrderSequence)
	for _, value := range p.periodList {
		p.systemPeriod[value] = OrderSequence{}
		sgFront := NewSystemGroup()
		sgFront.order = OrderFront
		sgAppend := NewSystemGroup()
		sgAppend.order = OrderAppend
		p.systemPeriod[value] = append(p.systemPeriod[value], sgFront, sgAppend)
	}
}

func (p *systemFlow) run(delta time.Duration) {
	var sq OrderSequence
	for _, period := range p.periodList {
		sq = p.systemPeriod[period]
		for _, sl := range sq {
			sl.reset()
			for ss := sl.next(); len(ss) > 0; ss = sl.next() {
				if systemCount := len(ss); systemCount != 0 {
					for i := 0; i < systemCount; i++ {
						imp := false
						var fn func(event Event)
						switch period {
						case PeriodCustomEvent:
							ss[i].eventDispatch()
						case PeriodStart:
							system, ok := ss[i].(IEventStart)
							fn = system.Start
							imp = ok
						case PeriodPostStart:
							system, ok := ss[i].(IEventPostStart)
							fn = system.PostStart
							imp = ok
						case PeriodUpdate:
							system, ok := ss[i].(IEventUpdate)
							fn = system.Update
							imp = ok
						case PeriodPostUpdate:
							system, ok := ss[i].(IEventPostUpdate)
							fn = system.PostUpdate
							imp = ok
						case PeriodDestroy:
							system, ok := ss[i].(IEventDestroy)
							fn = system.Destroy
							imp = ok
						case PeriodPostDestroy:
							system, ok := ss[i].(IEventPostDestroy)
							fn = system.PostDestroy
							imp = ok
						}

						if !imp {
							continue
						}

						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							fn(Event{Delta: delta})
							wg.Done()
						})
					}
				}
			}
			//waiting for all complete
			p.wg.Wait()
		}
	}

	p.wg.Wait()

	tasks := p.world.components.GetTempTasks()
	//Log.Info("temp task count:", len(tasks))
	newList := map[reflect.Type][]OperateInfo{}
	lock := sync.Mutex{}
	p.wg.Add(len(tasks))
	for _, task := range tasks {
		wg := p.wg
		fn := task
		Runtime.addJob(func() {
			typ, rn := fn()

			lock.Lock()
			newList[typ] = rn
			lock.Unlock()

			wg.Done()
		})
	}
	p.wg.Wait()

	//Log.Info("new component this frame:", len(newList))
	p.world.components.TempTasksDone(newList)
}

//register method only in world init or func init(){}
func (p *systemFlow) register(system ISystem) {
	//init function call
	Try(func() {
		system.baseInit(p.world, system)
	})

	order := system.Order()
	for _, period := range p.periodList {
		imp := false
		switch period {
		case PeriodStart:
			_, imp = system.(IEventStart)
		case PeriodPostStart:
			_, imp = system.(IEventPostStart)
		case PeriodUpdate:
			_, imp = system.(IEventUpdate)
		case PeriodPostUpdate:
			_, imp = system.(IEventPostUpdate)
		case PeriodDestroy:
			_, imp = system.(IEventDestroy)
		case PeriodPostDestroy:
			_, imp = system.(IEventPostDestroy)
		}

		if !imp {
			continue
		}

		sl := p.systemPeriod[period]
		if order == OrderFront {
			p.systemPeriod[period][0].insert(system)
		} else if order == OrderAppend {
			p.systemPeriod[period][len(sl)-1].insert(system)
		} else {
			for i, v := range sl {
				if order == v.order {
					v.insert(system)
					break
				} else if order < v.order {
					sg := NewSystemGroup()
					sg.order = order
					sg.insert(system)
					temp := append(OrderSequence{}, sl[i-1:]...)
					p.systemPeriod[period] = append(append(sl[:i-1], sg), temp...)
					break
				}
			}
		}
	}
}
