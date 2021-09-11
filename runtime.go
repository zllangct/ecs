package ecs

import "sync"

var Runtime = NewRuntime()

const(
	STATUS_INIT = iota
	STATUS_RUNNING
	STATUS_PAUSE
	STATUS_STOP
)

type RuntimeStatus int

type ecsRuntime struct {
	//mutex
	mutex sync.Mutex
	//runtime status
	status RuntimeStatus
	//world config
	config *RuntimeConfig
	//world worker pool
	workPool *Pool
	//logger
	logger IInternalLogger
	//world collections
	world []*World

	stop chan struct{}
}


//TODO runtime global event system

func NewRuntime() *ecsRuntime {
	return &ecsRuntime{}
}

func (r *ecsRuntime) NewWorld() *World{
	r.mutex.Lock()
	defer r.mutex.Unlock()

	world := NewWorld(r)
	r.world = append(r.world, NewWorld(r))

	return world
}

// SetConfig config the world
func (r *ecsRuntime) SetConfig(config *RuntimeConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.config = config
}

// SetLogger set logger
func (r *ecsRuntime) SetLogger(logger IInternalLogger) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.logger = logger
}

func (r *ecsRuntime) Run() {
	//default config
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.status == STATUS_INIT {
		config := NewDefaultRuntimeConfig()
		rt := &ecsRuntime{
			config:     config,
			logger:     NewStdLogger(),
		}
		rt.workPool = NewPool(rt, config.MaxPoolThread, config.MaxPoolJobQueue)

		r.logger.Info("start world success")
		//start the work pool
		r.workPool.Start()

		r.status = STATUS_RUNNING
	}

	//run all independent world
	for _, world := range r.world {
		if world.GetStatus() == STATUS_INIT {
			go world.Run()
		}
	}
}

func (r *ecsRuntime) Stop()  {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, world := range r.world {
		if status := world.GetStatus(); status != STATUS_STOP {
			world.Stop()
		}
	}

	r.stop<- struct{}{}
}

func (r *ecsRuntime) AddJob(handler func(JobContext, ...interface{}), args ...interface{}) {
	r.workPool.AddJob(handler, args...)
}

