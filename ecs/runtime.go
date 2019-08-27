package ecs

import (
	"reflect"
	"time"
)

type Runtime struct {
	systemGroup  map[reflect.Type]*SystemGroup
	root Entity
	UpdateInterval time.Duration
}

func (this *Runtime)Filter(system ISystem)  {
	requirements:=system.ComponentRequire()
	var systemGroup *SystemGroup
	for _,requirement := range requirements {
		typ:=requirement.Type()
		if sg,ok:=this.systemGroup[typ];ok {
			systemGroup = sg
			break
		}
	}
	if systemGroup == nil {
		systemGroup=&SystemGroup{}
	}
	for _,requirement := range requirements {
		typ:=requirement.Type()
		if _,ok:=this.systemGroup[typ];!ok {
			this.systemGroup[typ] = systemGroup
		}
	}
	systemGroup.Add(system)
}

func (this *Runtime)Run()  {
	for{
		this.update()
		time.Sleep(time.Millisecond * this.UpdateInterval)
	}
}

func (this Runtime)update()  {
	for _, sg := range this.systemGroup {
		go sg.Update()
	}
}