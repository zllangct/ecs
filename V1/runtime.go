package main

import (
	"container/list"
	"errors"
	"sync"
)

//error define
var (
	ErrInvalidRegisterTime = errors.New("can not register system this period")
)

type EcsConfig struct {
	CpuNum int
}

type ComponentName = string

type Runtime struct {
	sync.Mutex
	isInit        bool
	config        *EcsConfig
	systemFlow    *systemFlow
	components    map[ComponentName]interface{}
	componentsNew list.List
}

//config the runtime
func (p *Runtime) SetConfig(config *EcsConfig) {

}

//start ecs world
func (p *Runtime) Run() {
	p.systemFlow.Run(p)
}

//register system
func (p *Runtime) Register(system ISystem) {
	if p.isInit {
		panic(ErrInvalidRegisterTime)
	}
	p.systemFlow.register(system)
}
