package ecs

import "sync"

var Runtime = NewRuntime()

type ecsRuntime struct {
	//mutex
	mutex sync.Mutex
	//world config
	config *RuntimeConfig
	//world worker pool
	workPool *Pool
	//logger
	logger IInternalLogger
	//world collection
	world []*World
}

//TODO runtime global event system

func NewRuntime() *ecsRuntime {
	return &ecsRuntime{}
}

func (r *ecsRuntime) NewWorld() *World{
	r.mutex.Lock()
	defer r.mutex.Unlock()

	world := NewWorld(r)
	r.world = append(r.world, world)

	return world
}

// SetConfig config the world
func (r *ecsRuntime) SetConfig(config *RuntimeConfig) {
	r.config = config
}

// SetLogger set logger
func (w *World) SetLogger(logger IInternalLogger) {
	w.logger = logger
}

func (r *ecsRuntime) Run() {
	//default config
	config := NewDefaultRuntimeConfig()
	rt := &ecsRuntime{
		config:     config,
		logger:     NewStdLogger(),
	}
	rt.workPool = NewPool(rt, config.MaxPoolThread, config.MaxPoolJobQueue)

	r.logger.Info("start world success")
	//start the work pool
	r.workPool.Start()

	//run all independent world
	for _, world := range r.world {
		go world.Run()
	}
}

func (r *ecsRuntime) AddJob(handler func(JobContext, ...interface{}), args ...interface{}) {
	r.workPool.AddJob(handler, args...)
}

