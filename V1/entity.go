package main

import (
	"reflect"
	"sync"
)

type IEntity interface {
}

type Entity struct {
	sync.RWMutex

	index      *int
	components []IComponent

	ID string
}

func (p *Entity)Has(typ reflect.Type) bool  {
	p.RLock()
	for _, value := range p.components {
		if reflect.TypeOf(value) == typ {
			return true
		}
	}
	p.RUnlock()

	return false
}
