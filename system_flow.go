package ecs

import (
	"reflect"
	"sync"
	"time"
)

//system execute Order, 32bit + 32bit: period + suborder
type SystemOrder uint64

//system execute period:start->pre_update->update->pre_destroy->destroy
type SystemPeriod uint32

const (
	PERIOD_PRE_START SystemPeriod = iota
	PERIOD_START
	PERIOD_POST_START
	PERIOD_PRE_UPDATE
	PERIOD_UPDATE
	PERIOD_POST_UPDATE
	PERIOD_PER_DESTROY
	PERIOD_DESTROY
	PERIOD_POST_DESTROY
)

// Order default suborder of system
type Order int32

const (
	ORDER_FRONT   Order = -1
	ORDER_APPEND  Order = 999999
	ORDER_DEFAULT Order = ORDER_APPEND
)

type TempTask struct {
	lock *sync.Mutex
	m  map[reflect.Type][]ComponentOptResult
	wg *sync.WaitGroup
	fn func()(reflect.Type, []ComponentOptResult)
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
		PERIOD_START,
		PERIOD_POST_START,
		PERIOD_UPDATE,
		PERIOD_POST_UPDATE,
		PERIOD_DESTROY,
		PERIOD_POST_DESTROY,
	}
	p.systemPeriod = make(map[SystemPeriod]OrderSequence)
	for _, value := range p.periodList {
		p.systemPeriod[value] = OrderSequence{}
		sgFront := NewSystemGroup()
		sgFront.order = ORDER_FRONT
		sgAppend := NewSystemGroup()
		sgAppend.order = ORDER_APPEND
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
						case PERIOD_START:
							system, ok := ss[i].(IEventStart)
							fn = system.Start
							imp = ok
						case PERIOD_POST_START:
							system, ok := ss[i].(IEventPostStart)
							fn = system.PostStart
							imp = ok
						case PERIOD_UPDATE:
							system, ok := ss[i].(IEventUpdate)
							fn = system.Update
							imp = ok
						case PERIOD_POST_UPDATE:
							system, ok := ss[i].(IEventPostUpdate)
							fn = system.PostUpdate
							imp = ok
						case PERIOD_DESTROY:
							system, ok := ss[i].(IEventDestroy)
							fn = system.Destroy
							imp = ok
						case PERIOD_POST_DESTROY:
							system, ok := ss[i].(IEventPostDestroy)
							fn = system.PostDestroy
							imp = ok
						}

						if !imp {
							continue
						}

						p.wg.Add(1)
						Runtime.AddJob(func(ctx JobContext, args ...interface{}) {
							fn := args[0].(func(event Event))
							delta := args[1].(time.Duration)
							wg := args[2].(*sync.WaitGroup)
							fn(Event{Delta: delta})
							wg.Done()
						}, fn, delta, p.wg)
					}
				}
			}
			//waiting for all complete
			p.wg.Wait()
		}
	}

	p.wg.Wait()

	tasks := p.world.components.GetTempTasks()
	newList := map[reflect.Type][]ComponentOptResult{}
	l := sync.Mutex{}
	p.wg.Add(len(tasks))
	for _, task := range tasks{
		Runtime.AddJob(func(context JobContext, args ...interface{}) {
			Runtime.logger.Info("temp task execute")
			t := args[0].(TempTask)
			typ, rn := t.fn()

			t.lock.Lock()
			t.m[typ] = rn
			t.lock.Unlock()

			t.wg.Done()
		}, TempTask{
			fn: task,
			wg: p.wg,
			m: newList,
			lock: &l,
		})
	}
	p.wg.Wait()

	p.world.components.TempTasksDone(newList)
}

//register method only in world init or func init(){}
func (p *systemFlow) register(system ISystem) {
	//init function call
	err := Try(func() {
		system.baseInit(p.world, system)
	})
	if err != nil && p.world.logger != nil {
		p.world.logger.Error(err)
		return
	}

	order := system.Order()
	for _, period := range p.periodList {
		imp := false
		switch period {
		case PERIOD_START:
			_, imp = system.(IEventStart)
		case PERIOD_POST_START:
			_, imp = system.(IEventPostStart)
		case PERIOD_UPDATE:
			_, imp = system.(IEventUpdate)
		case PERIOD_POST_UPDATE:
			_, imp = system.(IEventPostUpdate)
		case PERIOD_DESTROY:
			_, imp = system.(IEventDestroy)
		case PERIOD_POST_DESTROY:
			_, imp = system.(IEventPostDestroy)
		}

		if !imp {
			continue
		}

		sl := p.systemPeriod[period]
		if order == ORDER_FRONT {
			p.systemPeriod[period][0].insert(system)
		} else if order == ORDER_APPEND {
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
