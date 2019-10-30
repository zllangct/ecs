package main

import (
	"sync"
	"time"
)

//system execute Order, 32bit + 32bit: period + suborder
type SystemOrder  uint64

//system execute period:start->pre_update->update->pre_destroy->destroy
type SystemPeriod  uint32
const(
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
const(
	ORDER_FRONT  Order = -1
	ORDER_APPEND Order = 999999
	ORDER_DEFAULT  Order = ORDER_APPEND
)

//extension of system group slice
type OrderSequence []*SystemGroup

//system execute flow
type systemFlow struct {
	runtime *Runtime
	systemPeriod map[SystemPeriod]OrderSequence
	periodList []SystemPeriod
}

func newSystemFlow(runtime *Runtime) *systemFlow {
	sf:= &systemFlow{
		runtime: runtime,
	}
	sf.init()
	return sf
}

//initialize the system flow
func (p *systemFlow)init()  {
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
		sgFront:= NewSystemGroup()
		sgFront.order = ORDER_FRONT
		sgAppend:= NewSystemGroup()
		sgAppend.order = ORDER_APPEND
		p.systemPeriod[value] = append(p.systemPeriod[value], sgFront, sgAppend)
	}
}

func (p * systemFlow)run(delta time.Duration)  {
	var sq OrderSequence
	for _,period := range p.periodList {
		sq = p.systemPeriod[period]
		for _, sl := range sq {
			sl.iterInit()
			wg:= sync.WaitGroup{}
			for ss:=sl.pop(); len(ss) >0 ;ss=sl.pop() {
				//work balance
				if len(ss) != 0 {
					interval := len(ss) / p.runtime.config.CpuNum
					remainder := len(ss) % p.runtime.config.CpuNum
					offset := 0
					if interval != 0{
						wg.Add(p.runtime.config.CpuNum)
						for i := 0; i < p.runtime.config.CpuNum; i++ {
							p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
								for _, sys := range args[0].([]ISystem) {
									sys.SystemUpdate(delta)
								}
								wg.Done()
							}, []interface{}{ss[offset : offset+interval]})
							offset += interval
						}
					}
					wg.Add(remainder)
					for i := 0; i < remainder; i++ {
						p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
							args[0].(ISystem).SystemUpdate(delta)
							wg.Done()
						}, []interface{}{ss[offset]})
						offset += 1
					}
				}
				//waiting for all complete
				wg.Wait()
			}

		}
		//filter execute in post destroy period
		if period == PERIOD_POST_DESTROY {
			p.runtime.components.TempFlush()
			p.FilterExecute()
		}
	}
}

func (p *systemFlow) FilterExecute()  {
	var sq OrderSequence
	for _,period := range p.periodList {
		sq = p.systemPeriod[period]
		for _, sl := range sq {
			wg:= sync.WaitGroup{}
			ss:=sl.all()
			//work balance
			if len(ss) != 0 {
				interval := len(ss) / p.runtime.config.CpuNum
				remainder := len(ss) % p.runtime.config.CpuNum
				offset := 0
				if interval != 0{
					wg.Add(p.runtime.config.CpuNum)
					for i := 0; i < p.runtime.config.CpuNum; i++ {
						p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
							for _, sys := range args[0].([]ISystem) {
								sys.Filter()
							}
							wg.Done()
						}, []interface{}{ss[offset : offset+interval]})
						offset += interval
					}
				}
				wg.Add(remainder)
				for i := 0; i < remainder; i++ {
					p.runtime.workPool.AddJob(func(ctx []interface{}, args ...interface{}) {
						args[0].(ISystem).Filter()
						wg.Done()
					}, []interface{}{ss[offset]})
					offset += 1
				}
			}
			//waiting for all complete
			wg.Wait()
		}
	}
}

//register method only in runtime init or func init(){}
func (p *systemFlow)register(system ISystem)  {
	system.Init(p.runtime)
	period,order:= system.GetOrder()
	sl:= p.systemPeriod[period]
	if order == ORDER_FRONT {
		sl[0].insert(system)
	}else if order == ORDER_APPEND {
		sl[len(sl)-1].insert(system)
	}else{
		for i, v := range sl {
			if order == v.order {
				v.insert(system)
				break
			}else if order < v.order {
				sg:= NewSystemGroup()
				sg.order = order
				sg.insert(system)
				temp := append(OrderSequence{},sl[i-1:]...)
				sl = append(append(sl[:i-1], sg), temp...)
				break
			}
		}
	}
}
