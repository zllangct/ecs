package main

//system execute order, 32bit + 32bit: period + suborder
type SystemOrder  uint64

//system execute period:start->pre_update->update->pre_destroy->destroy
type SystemPeriod  uint32
const(
	PERIOD_START SystemPeriod = iota
	PERIOD_PRE_UPDATE
	PERIOD_UPDATE
	PERIOD_POST_DESTROY
	PERIOD_DESTROY
)
//default system execute period
const PERIOD_DEFAULT = PERIOD_PRE_UPDATE

// default suborder of system
const SUBORDER_DEFAULT = -1

//extension of system group slice
type SystemList []*SystemGroup

//insert system to system group
func (p *SystemList)Insert(system ISystem, order int)  {
	for i, v := range *p {
		if order == v.order {
			v.insert(system)
			break
		}else if order < v.order {
			sg:= NewSystemGroup()
			sg.order = order
			sg.insert(system)
			temp := append(SystemList{},(*p)[i-1:]...)
			*p = append(append((*p)[:i-1], sg), temp...)
			break
		}
	}
}

//system execute flow
type systemFlow struct {
	systemPeriod []SystemList
}

//initialize the system flow
func (p *systemFlow)init()  {
	p.systemPeriod = make([]SystemList,5,5)
	p.systemPeriod[PERIOD_START] = SystemList{}
	p.systemPeriod[PERIOD_PRE_UPDATE] = SystemList{}
	p.systemPeriod[PERIOD_UPDATE] = SystemList{}
	p.systemPeriod[PERIOD_POST_DESTROY] = SystemList{}
	p.systemPeriod[PERIOD_DESTROY] = SystemList{}
}

func (p * systemFlow)Run(runtime *Runtime)  {

}

//register method only in runtime init or func init(){}
func (p *systemFlow)register(system ISystem)  {
	period,order:= system.GetBase().GetOrder()
	p.systemPeriod[period].Insert(system,order)
}
