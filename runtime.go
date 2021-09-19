package ecs

import (
	"runtime"
	"sync"
)

var Runtime = NewRuntime(NewDefaultRuntimeConfig())

var Log = Runtime.Logger()

const (
	StatusInit = iota
	StatusRunning
	StatusPause
	StatusStop
)

type RuntimeStatus int

type ecsRuntime struct {
	//mutex
	mutex sync.Mutex
	//config
	config *RuntimeConfig
	//world status
	status RuntimeStatus
	//world worker pool
	workPool *Pool
	//world collections
	world []*World

	stop chan struct{}
}

func NewRuntime(config *RuntimeConfig) *ecsRuntime {

	r := &ecsRuntime{
		config: config,
	}

	if r.config.MaxPoolThread <= 0 {
		r.config.MaxPoolThread = uint32(runtime.NumCPU())
	}

	if r.config.MaxPoolJobQueue <= 0 {
		r.config.MaxPoolJobQueue = 20
	}

	r.workPool = NewPool(config.MaxPoolThread, config.MaxPoolJobQueue)

	return r
}

func (r *ecsRuntime) NewWorld(config *WorldConfig) *World {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	world := NewWorld(r, config)
	r.world = append(r.world, NewWorld(r, config))

	return world
}

// SetLogger set logger
func (r *ecsRuntime) SetLogger(logger IInternalLogger) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.config.Logger = logger
	Log = r.config.Logger
}

func (r *ecsRuntime) Logger() IInternalLogger {
	return r.config.Logger
}

func (r *ecsRuntime) Status() RuntimeStatus {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.status
}

func (r *ecsRuntime) Run() {
	r.run()
}

func (r *ecsRuntime) run() {
	//default config
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.status == StatusInit {
		//start the work pool
		r.workPool.Start()
		r.status = StatusRunning
	}
}

func (r *ecsRuntime) Stop() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, world := range r.world {
		if status := world.GetStatus(); status != StatusStop {
			world.Stop()
		}
	}

	r.stop <- struct{}{}
}

func (r *ecsRuntime) AddJob(job func(), hashKey ...uint32) {
	r.workPool.Add(job, hashKey...)
}
