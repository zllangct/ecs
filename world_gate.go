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
	gate *IGate
}

func (g *GateInitializer) EventRegister(event CustomEventName, fn CustomEventHandler) {
	(*g.gate).eventRegister(event, fn)
}

type IGateInitializer interface {
	Init(initializer GateInitializer)
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
	syncQueue     []func(UtilityGetter)
	isInitialized bool
}

func (p Gate[T]) __GateIdentification() {}

func (p *Gate[T]) baseInit(w *ecsWorld) {
	p.world = w
	p.events = make(map[CustomEventName]CustomEventHandler)
	p.eventQueue = make([]CustomEvent, 0)
	p.syncQueue = make([]func(UtilityGetter), 0)
	ins := (*T)(unsafe.Pointer(p))
	if i, ok := any(ins).(IGateInitializer); ok {
		gi := GateInitializer{}
		*gi.gate = p
		i.Init(gi)
		*gi.gate = nil
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

	ug := UtilityGetter{}
	iw := iWorldBase(p.world)
	ug.world = &iw

	for _, e := range p.eventQueue {
		if fn, ok := p.events[e.Event]; ok {
			err := TryAndReport(func() {
				fn(ug, e.Args)
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
		fn(ug)
	}
	p.syncQueue = make([]func(UtilityGetter), 0)

	*ug.world = nil
	ug.world = nil
}

func (p *Gate[T]) Sync(fn func(getter UtilityGetter)) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.syncQueue = append(p.syncQueue, fn)
}

func GetUtility[T SystemObject, U UtilityObject](getter UtilityGetter) (*U, bool) {
	w := getter.getWorld()
	if w == nil {
		return nil, false
	}
	sys, ok := w.getSystem(TypeOf[T]())
	if !ok {
		return nil, false
	}
	u := (*U)(sys.GetUtility().getPointer())
	if u == nil {
		return nil, false
	}
	return u, true
}
