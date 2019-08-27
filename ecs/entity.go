package ecs

import "sync"

type Entity struct {
	lock sync.RWMutex

	ID         string
	components []IComponent
}

func (p *Entity) AddComponent(component IComponent) {
	p.lock.Lock()
	p.components = append(p.components, component)
	p.lock.Unlock()
}

func (p *Entity) Components() []IComponent {
	//p.lock.RLock()
	//defer p.lock.RUnlock()

	return p.components
}

