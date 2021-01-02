package ecs

import (
	"reflect"
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
	logger IInternalLogger
}

func NewRuntime() *Runtime {
	//default config
	config := NewDefaultRuntimeConfig()
	rt := &Runtime{
		config:     config,
		systemFlow: nil,
		components: NewComponentCollection(config.HashCount),
		entities:   NewEntityCollection(config.HashCount),
		logger:     NewStdLogger(),
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
func (p *Runtime) SetLogger(logger IInternalLogger) {
	p.logger = logger
}

//start ecs world
func (p *Runtime) Run() {
	p.logger.Info("start runtime success")
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

func (p *Runtime) ComponentAttach(target *Entity, com IComponent) {
	p.components.TempComponentOperate(target, com, COLLECTION_OPERATE_ADD)
}

func (p *Runtime) ComponentRemove(target *Entity, com IComponent) {
	p.components.TempComponentOperate(target, com, COLLECTION_OPERATE_DELETE)
}

func (p *Runtime) GetAllComponents() ComponentCollectionIter {
	return p.components.GetAllComponents()
}

func (p *Runtime) Error(v ...interface{}) {
	if p.logger != nil {
		p.logger.Error(v...)
	}
}

func (p *Runtime) getNewComponentsAll() []CollectionOperateInfo {
	return p.components.GetNewComponentsAll()
}

func (p *Runtime) getNewComponents(op CollectionOperate, typ reflect.Type) []CollectionOperateInfo {
	return p.components.GetNewComponents(op, typ)
}
