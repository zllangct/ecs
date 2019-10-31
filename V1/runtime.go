package main

import (
	"runtime"
	"sync"
	"time"
)

//error define
//var (
//	ErrInvalidRegisterTime = errors.New("can not register system this period")
//)

type EcsConfig struct {
	CpuNum          int           //使用的最大cpu数量
	FrameInterval   time.Duration //帧间隔
	MaxPoolThread   int           //线程池最大线程数量
	MaxPoolJobQueue int           //线程池最大任务队列长度
}

type Runtime struct {
	sync.Mutex

	//runtime config
	config *EcsConfig
	//system flow,all systems
	systemFlow *systemFlow
	//all components
	components *ComponentCollection
	//all entities
	entities *EntityCollection
	//runtime worker pool
	workPool *Pool
}

func NewRuntime() *Runtime {
	//default config
	config := &EcsConfig{
		CpuNum:        runtime.NumCPU(),
		FrameInterval: time.Millisecond * 33,
	}
	rt := &Runtime{
		config:     config,
		systemFlow: nil,
		components: NewComponentCollection(),
		entities:   NewEntityCollection(),
		workPool:   NewPool(config.MaxPoolThread, config.MaxPoolJobQueue),
	}
	//initialise system flow
	sf := newSystemFlow(rt)
	rt.systemFlow = sf
	//generate runtime
	return rt
}

//config the runtime
func (p *Runtime) SetConfig(config *EcsConfig) {
	p.config = config
}

//start ecs world
func (p *Runtime) Run() {
	//start the work pool
	p.workPool.Start()

	//main loop
	frameInterval := p.config.FrameInterval
	var ts time.Time
	var delta time.Duration
	for {
		ts = time.Now()
		p.systemFlow.run(delta)
		delta = time.Since(ts)
		if frameInterval-delta > 0 {
			time.Sleep(frameInterval - delta)
			delta = frameInterval
		}
	}
}

//register system
func (p *Runtime) Register(system ISystem) {
	p.systemFlow.register(system)
}

//entity operate : add
func (p *Runtime) AddEntity(entity *Entity) {
	p.entities.Add(entity)
}
//entity operate : delete
func (p *Runtime) DeleteEntity(entity *Entity) {
	p.entities.Delete(entity)
}
//entity operate : delete
func (p *Runtime) DeleteEntityByID(id uint64) {
	p.entities.DeleteByID(id)
}

func (p *Runtime) ComponentAttach(com IComponent) {
	p.components.TempComponentOperate(com,COLLECTION_OPERATE_ADD)
}
func (p *Runtime) ComponentRemove(com IComponent) {
	p.components.TempComponentOperate(com,COLLECTION_OPERATE_DELETE)
}

func (p *Runtime) GetComponentsNew() []*CollectionOperateInfo {
	return p.components.GetComponentsNew()
}

func (p *Runtime) GetAllComponents() []IComponent {
	return p.components.GetAllComponents()
}