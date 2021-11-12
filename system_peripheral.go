package ecs

import (
	"context"
	"sync"
)

type IPeripheralSystemTemplate interface {
	sdf7xxh2l23h4h56g7g4g3lh43()
	toPeripheralSystem(psys interface{}) IPeripheralSystem
}

type IPeripheralSystem interface {
	init()
	Stop()
	EventRegister(event CustomEventName, fn CustomEventHandler)
	Emit(event CustomEventName, args ...interface{})
}

type PeripheralSystem[T any] struct {
	lock          sync.RWMutex
	events        map[CustomEventName]CustomEventHandler
	eventQueue    chan CustomEvent
	cancel        context.CancelFunc
	isInitialized bool
}

func (p PeripheralSystem[T]) sdf7xxh2l23h4h56g7g4g3lh43() {}

func (p PeripheralSystem[T]) toPeripheralSystem(psys interface{}) IPeripheralSystem {
	return psys.(IPeripheralSystem)
}

func (p *PeripheralSystem[T]) init() {
	p.events = make(map[CustomEventName]CustomEventHandler)
	p.eventQueue = make(chan CustomEvent, 10)
	var ctx context.Context
	ctx, p.cancel = context.WithCancel(context.Background())
	go p.eventDispatch(ctx)
	p.isInitialized = true
}

func (p *PeripheralSystem[T]) Stop() {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if !p.isInitialized {
		Log.Error("must create peripheral system by ecs.NewPeripheralSystem[IPeripheralSystem]()")
		return
	}
	if p.cancel != nil {
		p.cancel()
	}
}

func (p *PeripheralSystem[T]) EventRegister(event CustomEventName, fn CustomEventHandler) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if !p.isInitialized {
		Log.Error("must create peripheral system by ecs.NewPeripheralSystem[IPeripheralSystem]()")
		return
	}

	p.events[event] = fn
}

func (p *PeripheralSystem[T]) Emit(event CustomEventName, args ...interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if !p.isInitialized {
		Log.Error("must create peripheral system by ecs.NewPeripheralSystem[IPeripheralSystem]()")
		return
	}
	p.eventQueue <- CustomEvent{
		Event: event,
		Args:  args,
	}
}

func (p *PeripheralSystem[T]) eventDispatch(ctx context.Context) {
	var e CustomEvent
	for {
		select {
		case <-ctx.Done():
			break
		case e = <-p.eventQueue:
			p.lock.RLock()
			if fn, ok := p.events[e.Event]; ok {
				err := TryAndReport(func() {
					fn(e.Args)
				})
				if err != nil {
					Log.Error(err)
				}
			} else {
				Log.Errorf("event not found: %p", e.Event)
			}
			p.lock.RUnlock()
		}
	}
}
