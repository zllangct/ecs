package ecs

import (
	"sync"
	"time"
)

type SyncWrapper struct {
	world *IWorld
}

func (g SyncWrapper) getWorld() IWorld {
	return *g.world
}

func (g SyncWrapper) NewEntity() Entity {
	return g.getWorld().newEntity().Entity()
}

func (g SyncWrapper) DestroyEntity(entity Entity) {
	info, ok := (*g.world).getEntityInfo(entity)
	if !ok {
		return
	}
	info.Destroy(*g.world)
}

func (g SyncWrapper) Add(entity Entity, components ...IComponent) {
	info, ok := (*g.world).getEntityInfo(entity)
	if !ok {
		return
	}
	info.Add(*g.world, components...)
}

func (g SyncWrapper) Remove(entity Entity, components ...IComponent) {
	info, ok := (*g.world).getEntityInfo(entity)
	if !ok {
		return
	}
	info.Remove(*g.world, components...)
}

type syncTask struct {
	wait chan struct{}
	fn   func(wrapper SyncWrapper) error
}

type AsyncWorld struct {
	ecsWorld
	lock        sync.Mutex
	syncQueue   []syncTask
	wStop       chan struct{}
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
				w.workPool.Release()
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

func (w *AsyncWorld) Stop() {
	w.wStop <- struct{}{}
}

func (w *AsyncWorld) dispatch() {
	w.lock.Lock()
	defer w.lock.Unlock()

	gaw := SyncWrapper{}
	ig := IWorld(w)
	gaw.world = &ig
	for _, task := range w.syncQueue {
		err := TryAndReport(func() error {
			return task.fn(gaw)
		})
		if err != nil {
			Log.Error(err)
		}
		if task.wait != nil {
			task.wait <- struct{}{}
		}
	}
	w.syncQueue = make([]syncTask, 0)

	*gaw.world = nil
	gaw.world = nil
}

func (w *AsyncWorld) Sync(fn func(g SyncWrapper) error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.syncQueue = append(w.syncQueue, syncTask{
		wait: nil,
		fn:   fn,
	})
}

func (w *AsyncWorld) Wait(fn func(g SyncWrapper) error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	wait := make(chan struct{})
	w.syncQueue = append(w.syncQueue, syncTask{
		wait: wait,
		fn:   fn,
	})
	<-wait
}
