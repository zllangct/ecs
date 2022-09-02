package ecs

import (
	"sync"
	"unsafe"
)

type GateObject interface {
	__GateIdentification()
}

type GatePointer[T GateObject] interface {
	IGate
	*T
}

type GateInitializer struct {
	gate IGate
}

func (g *GateInitializer) EventRegister(event CustomEventName, fn CustomEventHandler) {
	g.gate.eventRegister(event, fn)
}

type IGateInitializer interface {
	Init(initializer GateInitializer)
}

type GateApi struct {
	world *ecsWorld
}

func (g *GateApi) getWorld() IWorld {
	return g.world
}

type IGate interface {
	baseInit(w *ecsWorld)
	dispatch()
	eventRegister(event CustomEventName, fn CustomEventHandler)
	resetData(src *IGate)
	getWorld() *ecsWorld
	Emit(event CustomEventName, args ...interface{})
}

type Gate[T GateObject] struct {
	lock          sync.Mutex
	world         *ecsWorld
	events        map[CustomEventName]CustomEventHandler
	eventQueue    []CustomEvent
	syncQueue     []func(*GateApi)
	isInitialized bool
}

func (p Gate[T]) __GateIdentification() {}

func (p *Gate[T]) baseInit(w *ecsWorld) {
	p.world = w
	p.events = make(map[CustomEventName]CustomEventHandler)
	p.eventQueue = make([]CustomEvent, 0)
	p.syncQueue = make([]func(*GateApi), 0)
	ins := (*T)(unsafe.Pointer(p))
	if i, ok := any(ins).(IGateInitializer); ok {
		i.Init(GateInitializer{
			gate: p,
		})
	}
	p.isInitialized = true
}

func (p *Gate[T]) resetData(src *IGate) {
	cp := *(*T)(unsafe.Pointer(p))
	(*iface)(unsafe.Pointer(src)).data = unsafe.Pointer(&cp)
}

func (p *Gate[T]) getWorld() *ecsWorld {
	return p.world
}

func (p *Gate[T]) eventRegister(event CustomEventName, fn CustomEventHandler) {
	p.events[event] = fn
}

func (p *Gate[T]) Emit(event CustomEventName, args ...interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.isInitialized {
		return
	}

	p.eventQueue = append(p.eventQueue, CustomEvent{Event: event, Args: args})
}

func (p *Gate[T]) dispatch() {
	p.lock.Lock()
	defer p.lock.Unlock()

	api := &GateApi{world: p.world}
	for _, e := range p.eventQueue {
		if fn, ok := p.events[e.Event]; ok {
			err := TryAndReport(func() {
				fn(api, e.Args)
			})
			if err != nil {
				Log.Error(err)
			}
		} else {
			Log.Errorf("event not found: %p", e.Event)
		}
	}
	p.eventQueue = make([]CustomEvent, 0)

	for _, fn := range p.syncQueue {
		fn(api)
	}
	p.syncQueue = make([]func(*GateApi), 0)

	api.world = nil
}

func (p *Gate[T]) Sync(fn func(api *GateApi)) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.syncQueue = append(p.syncQueue, fn)
}

func GetUtilityByGate[T SystemObject, U UtilityObject](g *GateApi) (*U, bool) {
	sys, ok := g.getWorld().GetSystem(TypeOf[T]())
	if !ok {
		return nil, false
	}
	u := (*U)(sys.GetUtility().getPointer())
	if u == nil {
		return nil, false
	}
	return u, true
}
