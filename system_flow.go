package ecs

import (
	"reflect"
	"sync"
	"time"
)

const (
	PeriodStart Period = iota
	PeriodPreUpdate
	PeriodUpdate
	PeriodPostUpdate
	PeriodDestroy
)

// Period system execute period:start->pre_update->update->pre_destroy->destroy
type Period uint32

// Order default suborder of system
type Order int32

const (
	OrderFront   Order = -1
	OrderInvalid Order = 0
	OrderAppend  Order = 99999999
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
	world        *ecsWorld
	systemPeriod map[Period]OrderSequence
	periodList   []Period
	wg            *sync.WaitGroup
}

func newSystemFlow(runtime *ecsWorld) *systemFlow {
	sf := &systemFlow{
		world: runtime,
		wg:    &sync.WaitGroup{},
	}
	sf.init()
	return sf
}

//initialize the system flow
func (p *systemFlow) init() {
	p.periodList = []Period{
		PeriodStart,
		PeriodPreUpdate,
		PeriodUpdate,
		PeriodPostUpdate,
		PeriodDestroy,
	}
	p.systemPeriod = make(map[Period]OrderSequence)
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
						sys := ss[i]
						imp := false
						var fn func(event Event)
						state := ss[i].getState()
						if state == SystemStateInit {
							switch period {
							case PeriodStart:
								system, ok := ss[i].(IEventStart)
								fn = system.Start
								imp = ok
							}
							sys.setState(SystemStateUpdate)
						} else if state == SystemStateUpdate {
							switch period {
							case PeriodPreUpdate:
								system, ok := sys.(IEventPreUpdate)
								fn = system.PreUpdate
								imp = ok
							case PeriodUpdate:
								system, ok := ss[i].(IEventUpdate)
								fn = system.Update
								imp = ok
							case PeriodPostUpdate:
								system, ok := ss[i].(IEventPostUpdate)
								fn = system.PostUpdate
								imp = ok
							}
						} else if state == SystemStateDestroy {
							switch period {
							case PeriodDestroy:
								system, ok := ss[i].(IEventDestroy)
								fn = system.Destroy
								imp = ok
								sys.setState(SystemStateInvalid)
							}
						}

						if !imp {
							continue
						}

						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							defer wg.Done()
							fn(Event{Delta: delta})
						})
					}
				}
				p.wg.Wait()
			}
		}
	}

	p.wg.Wait()

	for _, period := range p.periodList {
		sq = p.systemPeriod[period]
		for _, sl := range sq {
			sl.reset()
			for ss := sl.next(); len(ss) > 0; ss = sl.next() {
				if systemCount := len(ss); systemCount != 0 {
					for i := 0; i < systemCount; i++ {
						fn := ss[i].eventDispatch
						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							defer wg.Done()
							fn()
						})
					}
				}
				p.wg.Wait()
			}
		}
	}

	p.wg.Wait()

	p.world.components.ClearDisposable()

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
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, reset order to %d", OrderAppend + 1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.periodList {
		imp := false
		switch period {
		case PeriodStart:
			_, imp = system.(IEventStart)
		case PeriodPreUpdate:
			_, imp = system.(IEventPreUpdate)
		case PeriodUpdate:
			_, imp = system.(IEventUpdate)
		case PeriodPostUpdate:
			_, imp = system.(IEventPostUpdate)
		case PeriodDestroy:
			_, imp = system.(IEventDestroy)
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