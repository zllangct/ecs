package ecs

import "time"

type IWorld interface {
	iWorldBase
	Startup()
	Update()
	Optimize(t time.Duration, force bool)
}

type WorldBase struct {
	ecsWorld
}

func (w *WorldBase) Init(config *WorldConfig) {
	w.ecsWorld.init(config)
}

func (w *WorldBase) Startup() {
	w.startup()
}

func (w *WorldBase) Update() {
	w.update()
}

func (w *WorldBase) Optimize(t time.Duration, force bool) {}

type AsyncWorld struct {
	ecsWorld
	wStop       chan struct{}
	gate        IGate
	stopHandler func(world *AsyncWorld)
}

func NewAsyncWorld(config *WorldConfig) *AsyncWorld {
	w := &AsyncWorld{
		wStop: make(chan struct{}),
	}
	w.ecsWorld.init(config)
	return w
}

func (w *AsyncWorld) Startup() {
	w.startup()

	frameInterval := w.config.FrameInterval
	w.setStatus(WorldStatusRunning)
	Log.Info("start world success")

	//main loop
	for {
		select {
		case <-w.wStop:
			w.setStatus(WorldStatusStop)
			if w.stopHandler != nil {
				w.stopHandler(w)
			}
			w.systemFlow.stop()
			return
		default:
		}
		if w.gate != nil {
			w.gate.dispatch()
		}
		w.update()
		//world.Info(delta, frameInterval - delta)
		if d := frameInterval - w.delta; d > 0 {
			time.Sleep(d)
		}
	}
}

func (w *AsyncWorld) Update() {
	w.update()
}

func (w *AsyncWorld) Optimize(t time.Duration, force bool) {}

func (w *AsyncWorld) SetGate(gate IGate) IGate {
	w.gate = gate
	gate.resetData(&w.gate)
	w.gate.baseInit(&w.ecsWorld)
	return w.gate
}

func (w *AsyncWorld) GetGate() IGate {
	return w.gate
}

func (w *AsyncWorld) Stop() {
	w.wStop <- struct{}{}
}

type SyncWorld struct {
	ecsWorld
}

func NewSyncWorld(config *WorldConfig) *SyncWorld {
	w := &SyncWorld{}
	w.ecsWorld.init(config)
	return w
}

func (w *SyncWorld) Startup() {
	w.startup()
}

func (w *SyncWorld) Update() {
	w.update()
}

func (w *SyncWorld) Optimize(t time.Duration, force bool) {}
