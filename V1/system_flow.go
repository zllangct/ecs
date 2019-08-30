package main

//system execute order, 32bit + 32bit: period + suborder
type SystemOrder  uint64

//system execute period:start->pre_update->update->pre_destroy->destroy
type SystemPeriod  uint32
const(
	SYSTEM_PERIOD_START SystemPeriod = iota
	SYSTEM_PERIOD_PRE_UPDATE
	SYSTEM_PERIOD_UPDATE
	SYSTEM_PERIOD_PRE_DESTROY
	SYSTEM_PERIOD_DESTROY
)
//default system execute period
const SYSTEM_PERIOD_DEFAULT  = SYSTEM_PERIOD_PRE_UPDATE

//extension of system group slice
type SystemList []SystemGroup

func (p *SystemList)Insert(system ISystem,index int)  {
	for i, v := range *p {
		if index == i {
			v.Insert(system)
			break
		}else if index < i {
			sg:=SystemGroup{}
			sg.Insert(system)
			temp := append(SystemList{},(*p)[i-1:]...)
			*p = append(append((*p)[:i-1], sg), temp...)
			break
		}
	}

	sg:=SystemGroup{}
	sg.Insert(system)
	*p = append(*p, sg)
}


type systemFlow struct {
	systemList []SystemList
}

func (p *systemFlow)init()  {
	p.systemList = make([]SystemList,5,5)
	p.systemList[SYSTEM_PERIOD_START] = SystemList{}
	p.systemList[SYSTEM_PERIOD_PRE_UPDATE] = SystemList{}
	p.systemList[SYSTEM_PERIOD_UPDATE] = SystemList{}
	p.systemList[SYSTEM_PERIOD_PRE_DESTROY] = SystemList{}
	p.systemList[SYSTEM_PERIOD_DESTROY] = SystemList{}
}

//register method only in runtime init or func init(){}
func (p *systemFlow)register(system ISystem)  {
	period,order:= system.GetOrder()
	p.systemList[period].Insert(system,order)
}
