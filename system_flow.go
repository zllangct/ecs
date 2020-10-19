package ecs

import (
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

//default system execute period
const PERIOD_DEFAULT = PERIOD_UPDATE

// default suborder of system
type Order int32

const (
	ORDER_FRONT   Order = -1
	ORDER_APPEND  Order = 999999
	ORDER_DEFAULT Order = ORDER_APPEND
)

//extension of system group slice
type OrderSequence []*SystemGroup

//system execute flow
type systemFlow struct {
	sync.Mutex
	runtime      *Runtime
	systemPeriod map[SystemPeriod]OrderSequence
	periodList   []SystemPeriod
	wg           *sync.WaitGroup
}

func newSystemFlow(runtime *Runtime) *systemFlow {
	sf := &systemFlow{
		runtime: runtime,
		wg:      &sync.WaitGroup{},
	}
	sf.init()
	return sf
}

//initialize the system flow
func (p *systemFlow) init() {
	p.periodList = []SystemPeriod{
		PERIOD_PRE_START,
		PERIOD_START,
		PERIOD_POST_START,
		PERIOD_PRE_UPDATE,
		PERIOD_UPDATE,
		PERIOD_POST_UPDATE,
		PERIOD_PER_DESTROY,
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
					p.wg.Add(p.runtime.config.CpuNum)
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

						p.runtime.workPool.AddJob(func(ctx *JobContext, args ...interface{}) {
							fn := args[0].(func(event Event))
							delta := args[1].(time.Duration)
							wg := args[1].(*sync.WaitGroup)
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
	//do filter
	p.runtime.components.TempFlush()
	p.filterExecute()
}

func (p *systemFlow) filterExecute() {
	var sq OrderSequence
	comInfos := p.runtime.GetComponentsNew()
	for _, period := range p.periodList {
		sq = p.systemPeriod[period]
		for _, sl := range sq {
			ss := sl.all()
			if systemCount := len(ss); systemCount != 0 {
				p.wg.Add(systemCount)
				for i := 0; i < systemCount; i++ {
					p.runtime.workPool.AddJob(func(ctx *JobContext, args ...interface{}) {
						sys := args[0].(ISystem)
						wg := args[1].(*sync.WaitGroup)
						if !sys.GetBase().isPreFilter {
							cpts := ctx.Runtime.GetAllComponents()
							for _, com := range cpts {
								sys.Filter(com, COLLECTION_OPERATE_ADD)
							}
							sys.GetBase().isPreFilter = true
						}
						for _, comInfo := range comInfos {
							sys.Filter(comInfo.com, comInfo.op)
						}
						wg.Done()
					}, ss[i], p.wg)
				}
			}
			//waiting for all complete
			p.wg.Wait()
		}
	}
}

//register method only in runtime init or func init(){}
func (p *systemFlow) register(system ISystem) {
	system.Init(p.runtime)
	order := system.GetOrder()

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
		if sys, ok := system.(IEventInit); ok {
			err := Try(func() {
				sys.Initialize()
			})
			if err != nil && p.runtime.logger != nil {
				p.runtime.logger.Error(err)
			}
		}
	}
}
