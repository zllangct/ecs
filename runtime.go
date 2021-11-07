package ecs

import (
	"runtime"
	"sync"
)

var Runtime = newRuntime()

var Log Logger = NewStdLog()

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
	//world rtStatus
	rtStatus RuntimeStatus
	//world worker pool
	workPool *Pool
	//world collections
	world []*ecsWorld

	isInited bool
	rtStop   chan struct{}
}

func newRuntime() *ecsRuntime {
	return &ecsRuntime{}
}

func (r *ecsRuntime) Configure(config *RuntimeConfig) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.config = config
	if r.config.MaxPoolThread <= 0 {
		r.config.MaxPoolThread = uint32(runtime.NumCPU())
	}

	if r.config.MaxPoolJobQueue <= 0 {
		r.config.MaxPoolJobQueue = 20
	}

	r.workPool = NewPool(config.MaxPoolThread, config.MaxPoolJobQueue)

	r.isInited = true
}

func (r *ecsRuntime) newWorld(config *WorldConfig) *ecsWorld {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInited {
		panic("you must config the runtime first")
	}

	world := newWorld(r, config)
	r.world = append(r.world, newWorld(r, config))

	return world
}

func (r *ecsRuntime) destroyWorld(world *ecsWorld) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := 0; i < len(r.world); i++ {
		if r.world[i].id == world.id {
			r.world = append(r.world[:i], r.world[i+1:]...)
			return
		}
	}

	world.stop()
}

// SetLogger set logger
func (r *ecsRuntime) setLogger(logger Logger) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInited {
		panic("you must config the runtime first")
	}

	r.config.Logger = logger
	Log = r.config.Logger
}

func (r *ecsRuntime) logger() Logger {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInited {
		panic("you must config the runtime first")
	}

	return r.config.Logger
}

func (r *ecsRuntime) status() RuntimeStatus {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	return r.rtStatus
}

func (r *ecsRuntime) run() {
	//default config
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInited {
		panic("you must config the runtime first")
	}

	if r.rtStatus == StatusInit {
		//start the work pool
		r.workPool.Start()
		r.rtStatus = StatusRunning
	}
}

func (r *ecsRuntime) stop() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.isInited {
		panic("you must config the runtime first")
	}

	for _, world := range r.world {
		if status := world.GetStatus(); status != StatusStop {
			world.stop()
		}
	}

	r.rtStop <- struct{}{}
}

func (r *ecsRuntime) addJob(job func(), hashKey ...uint32) {
	if !r.isInited {
		panic("you must config the runtime first")
	}
	r.workPool.Add(job, hashKey...)
}
