package ecs

import "time"

type WorldLauncher struct {
	world       *ecsWorld
	wStop       chan struct{}
	gate        IGate
	stopHandler func(world *ecsWorld)
}

func NewAsyncWorldLauncher(w *ecsWorld) *WorldLauncher {
	return &WorldLauncher{
		world: w,
		wStop: make(chan struct{}),
	}
}

func (w *WorldLauncher) SetGate(gate IGate) IGate {
	w.gate = gate
	gate.resetData(&w.gate)
	w.gate.baseInit(w.world)
	return w.gate
}

func (w *WorldLauncher) GetGate() IGate {
	return w.gate
}

func (w *WorldLauncher) Run() {
	if w.world.GetStatus() != WorldStatusInit {
		panic("this world is already running.")
	}
	frameInterval := w.world.config.FrameInterval
	w.world.setStatus(WorldStatusRunning)
	Log.Info("start world success")

	//main loop
	for {
		select {
		case <-w.wStop:
			w.world.setStatus(WorldStatusStop)
			if w.stopHandler != nil {
				w.stopHandler(w.world)
			}
			w.world.systemFlow.stop()
			return
		default:
		}
		if w.gate != nil {
			w.gate.dispatch()
		}
		w.world.Update()
		//world.Info(delta, frameInterval - delta)
		if d := frameInterval - w.world.delta; d > 0 {
			time.Sleep(d)
		}
	}
}

func (w *WorldLauncher) Stop() {
	w.wStop <- struct{}{}
}
