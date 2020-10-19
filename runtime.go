package ecs

import (
	"sync"
	"time"
)

type Runtime struct {
	sync.Mutex
	//runtime config
	config *RuntimeConfig
	//system flow,all systems
	systemFlow *systemFlow
	//all components
	components *ComponentCollection
	//all entities
	entities *EntityCollection
	//runtime worker pool
	workPool *Pool
	//logger
	logger ILogger
}

func NewRuntime() *Runtime {
	//default config
	config := NewDefaultRuntimeConfig()
	rt := &Runtime{
		config:     config,
		systemFlow: nil,
		components: NewComponentCollection(),
		entities:   NewEntityCollection(),
	}
	rt.workPool = NewPool(rt, config.MaxPoolThread, config.MaxPoolJobQueue)
	//initialise system flow
	sf := newSystemFlow(rt)
	rt.systemFlow = sf
	//generate runtime
	return rt
}

//config the runtime
func (p *Runtime) SetConfig(config *RuntimeConfig) {
	p.config = config
}

//set logger
func (p *Runtime) SetLogger(logger ILogger) {
	p.logger = logger
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
	p.entities.add(entity)
}

//entity operate : delete
func (p *Runtime) DeleteEntity(entity *Entity) {
	p.entities.delete(entity)
}

//entity operate : delete
func (p *Runtime) DeleteEntityByID(id uint64) {
	p.entities.deleteByID(id)

}

func (p *Runtime) ComponentAttach(com IComponent) {
	p.components.TempComponentOperate(com, COLLECTION_OPERATE_ADD)
}

func (p *Runtime) ComponentRemove(com IComponent) {
	p.components.TempComponentOperate(com, COLLECTION_OPERATE_DELETE)
}

func (p *Runtime) GetComponentsNew() []*CollectionOperateInfo {
	return p.components.GetComponentsNew()
}

func (p *Runtime) GetAllComponents() []IComponent {
	return p.components.GetAllComponents()
}
