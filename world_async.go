package ecs

import (
	"sync"
	"time"
)

type SyncWrapper struct {
	world *iWorldBase
}

func (g SyncWrapper) getWorld() iWorldBase {
	return *g.world
}

func (g SyncWrapper) NewEntity() *EntityInfo {
	return g.getWorld().newEntity()
}

type AsyncWorld struct {
	worldBase
	lock        sync.Mutex
	syncQueue   []func(wrapper SyncWrapper)
	wStop       chan struct{}
	stopHandler func(world *AsyncWorld)
}

func NewAsyncWorld(config *WorldConfig) *AsyncWorld {
	w := &AsyncWorld{
		wStop: make(chan struct{}),
	}
	w.worldBase.init(config)
	return w
}

func (w *AsyncWorld) Startup() {
	//main loop
	go func() {
		w.startup()

		frameInterval := w.config.FrameInterval
		w.setStatus(WorldStatusRunning)
		Log.Info("start world success")

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
			if w != nil {
				w.dispatch()
			}
			w.update()
			//world.Info(delta, frameInterval - delta)
			if d := frameInterval - w.delta; d > 0 {
				time.Sleep(d)
			}
		}
	}()
}

func (w *AsyncWorld) Optimize(t time.Duration, force bool) {}

func (w *AsyncWorld) Stop() {
	w.wStop <- struct{}{}
}

func (w *AsyncWorld) getWorld() iWorldBase {
	return w
}

func (w *AsyncWorld) dispatch() {
	w.lock.Lock()
	defer w.lock.Unlock()

	gaw := SyncWrapper{}
	ig := iWorldBase(w)
	gaw.world = &ig
	for _, fn := range w.syncQueue {
		err := TryAndReport(func() {
			fn(gaw)
		})
		if err != nil {
			Log.Error(err)
		}
	}
	w.syncQueue = make([]func(SyncWrapper), 0)

	*gaw.world = nil
	gaw.world = nil
}

func (w *AsyncWorld) Sync(fn func(g SyncWrapper)) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.syncQueue = append(w.syncQueue, fn)
}
