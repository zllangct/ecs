package ecs

import (
	"reflect"
	"sync"
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
	lock         sync.Mutex
	world        *ecsWorld
	systemPeriod map[Period]OrderSequence
	periodList   []Period
	systems      sync.Map
	wg           *sync.WaitGroup
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
	p.reset()
}

func (p *systemFlow) reset() {
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

func (p *systemFlow) run(event Event) {
	p.lock.Lock()
	defer p.lock.Unlock()

	removeList := map[int64]ISystem{}

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
								system, ok := ss[i].(StartReceiver)
								fn = system.Start
								imp = ok
							}
							sys.setState(SystemStateUpdate)
						} else if state == SystemStateUpdate {
							switch period {
							case PeriodPreUpdate:
								system, ok := sys.(PreUpdateReceiver)
								fn = system.PreUpdate
								imp = ok
							case PeriodUpdate:
								system, ok := ss[i].(UpdateReceiver)
								fn = system.Update
								imp = ok
							case PeriodPostUpdate:
								system, ok := ss[i].(PostUpdateReceiver)
								fn = system.PostUpdate
								imp = ok
							}
						} else if state == SystemStateDestroy {
							switch period {
							case PeriodDestroy:
								system, ok := ss[i].(DestroyReceiver)
								fn = system.Destroy
								imp = ok
								sys.setState(SystemStateDestroyed)
								removeList[sys.ID()] = sys
							}
						}

						if !imp {
							continue
						}

						p.wg.Add(1)
						wg := p.wg
						Runtime.addJob(func() {
							defer wg.Done()
							fn(event)
						})
					}
				}
				p.wg.Wait()
			}
		}
	}

	buckets := p.world.entities.getBuckets()
	for _, bucket := range buckets {
		b := bucket
		wg := p.wg
		wg.Add(1)
		Runtime.addJob(func() {
			b.Range(func(key Entity, value *EntityInfo) bool {
				value.clearDisposable()
				return true
			})
			wg.Done()
		})
	}
	p.world.components.clearDisposable()
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

	tasks := p.world.components.getTempTasks()
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
	p.world.components.tempTasksDone(newList)

	//do something clean
	for _, system := range removeList {
		p.unregister(system)
	}
}

//register method only in world init or func init(){}
func (p *systemFlow) register(system ISystem) {
	p.lock.Lock()
	defer p.lock.Unlock()

	//init function call
	system.baseInit(p.world, system)

	order := system.Order()
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, reset order to %d", OrderAppend+1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.periodList {
		imp := false
		switch period {
		case PeriodStart:
			_, imp = system.(StartReceiver)
		case PeriodPreUpdate:
			_, imp = system.(PreUpdateReceiver)
		case PeriodUpdate:
			_, imp = system.(UpdateReceiver)
		case PeriodPostUpdate:
			_, imp = system.(PostUpdateReceiver)
		case PeriodDestroy:
			_, imp = system.(DestroyReceiver)
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

func (p *systemFlow) unregister(system ISystem) {
	p.lock.Lock()
	defer p.lock.Unlock()

	order := system.Order()
	if order > OrderAppend {
		Log.Errorf("system order must less then %d, reset order to %d", OrderAppend+1, OrderAppend)
		order = OrderAppend
	}

	for _, period := range p.periodList {
		imp := false
		switch period {
		case PeriodStart:
			_, imp = system.(StartReceiver)
		case PeriodPreUpdate:
			_, imp = system.(PreUpdateReceiver)
		case PeriodUpdate:
			_, imp = system.(UpdateReceiver)
		case PeriodPostUpdate:
			_, imp = system.(PostUpdateReceiver)
		case PeriodDestroy:
			_, imp = system.(DestroyReceiver)
		}

		if !imp {
			continue
		}

		sl := p.systemPeriod[period]
		for _, group := range sl {
			if group.has(system) {
				group.remove(system)
			}
		}
	}
}

func (p *systemFlow) stop() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.reset()
}
