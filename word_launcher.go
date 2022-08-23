package ecs

import "time"

type WorldLauncher struct {
	world       *ecsWorld
	wStop       chan struct{}
	stopHandler func(world *ecsWorld)
}

func DefaultAsyncWorldLauncher(w *ecsWorld) *WorldLauncher {
	return &WorldLauncher{
		world: w,
		wStop: make(chan struct{}),
	}
}

// Run start ecs world
func (w *WorldLauncher) Run() {
	go w.run()
}

func (w *WorldLauncher) run() {
	w.world.mutex.Lock()
	if w.world.status != WorldStatusInit {
		Log.Error("this world is already running.")
		return
	}
	frameInterval := w.world.config.DefaultFrameInterval
	w.world.status = WorldStatusRunning
	w.world.mutex.Unlock()

	Log.Info("start world success")

	defer func() {
		w.world.mutex.Lock()
		w.world.status = WorldStatusStop
		w.world.mutex.Unlock()
	}()

	//main loop
	for {
		select {
		case <-w.wStop:
			w.world.mutex.Lock()
			if w.stopHandler != nil {
				w.stopHandler(w.world)
			}
			w.world.systemFlow.stop()
			w.world.mutex.Unlock()
			return
		default:
		}
		w.world.Update()
		//w.Info(delta, frameInterval - delta)
		if d := frameInterval - w.world.delta; d > 0 {
			time.Sleep(d)
		}
	}
}

func (w *WorldLauncher) stop() {
	w.wStop <- struct{}{}
}
