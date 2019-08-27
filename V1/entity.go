package main

import (
	"reflect"
	"sync"
)

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

func (p *Entity)GetComponent(typ reflect.Type) interface{}  {
	p.RLock()
	for _, value := range p.components {
		if reflect.TypeOf(value) == typ {
			return value
		}
	}
	p.RUnlock()
	return nil
}

func (p *Entity)GetComponents(typs ...reflect.Type) []interface{}  {
	cmps:=make([]interface{},0,len(typs))
	p.RLock()
	for index, typ := range typs {
		for _, value := range p.components {
			if reflect.TypeOf(value) == typ {
				cmps[index] = value
			}
		}
	}
	p.RUnlock()
	return cmps
}