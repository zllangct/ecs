package main

import "sync"

type EcsConfig struct {

}

type Runtime struct {
	sync.RWMutex
 	config *EcsConfig
}

//config the runtime
func (p *Runtime)SetConfig(config *EcsConfig)  {

}

//start ecs world
func (p *Runtime)Run()  {

}