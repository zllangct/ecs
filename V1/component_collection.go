package main

import (
	"reflect"
	"sync"
)

type componentData struct {
	data  []IComponent
	index map[uint64]int
}

func newComponentData() *componentData {
	return &componentData{
		data:  []IComponent{},
		index: map[uint64]int{},
	}
}

func (p *componentData) push(com IComponent, id uint64) {
	p.data = append(p.data, com)
	p.index[id] = len(p.data) - 1
}

func (p *componentData) pop(id uint64) {
	if index, ok := p.index[id]; ok {
		length := len(p.data)
		p.data[index], p.data[length-1] = p.data[length-1], p.data[index]
		if length > 0 {
			p.data = p.data[:length-1]
		}
	}
}

type CollectionOperate int

const (
	COLLECTION_OPERATE_NONE   CollectionOperate = iota
	COLLECTION_OPERATE_ADD                      //add component operation
	COLLECTION_OPERATE_DELETE                   //delete component operation
)

type CollectionOperateInfo struct {
	com IComponent
	op  CollectionOperate
}

func NewCollectionOperateInfo(com IComponent, op CollectionOperate) *CollectionOperateInfo {
	return &CollectionOperateInfo{com: com, op: op}
}

type ComponentCollection struct {
	collection map[reflect.Type]*componentData
	//new component cache
	lockInput        sync.Mutex
	componentsTemp   []*CollectionOperateInfo
	componentsNew   []*CollectionOperateInfo
}

func NewComponentCollection() *ComponentCollection {
	return &ComponentCollection{
		collection:     map[reflect.Type]*componentData{},
		lockInput:      sync.Mutex{},
		componentsTemp: make([]*CollectionOperateInfo, 0, 10),
		componentsNew: make([]*CollectionOperateInfo, 0),
	}
}

//new component temp
func (p *ComponentCollection) TempComponentOperate(com IComponent, op CollectionOperate) {
	p.lockInput.Lock()
	p.componentsTemp = append(p.componentsTemp, NewCollectionOperateInfo(com, op))
	p.lockInput.Unlock()
}

//handle and flush new components,should called before destroy period
func (p *ComponentCollection) TempFlush() {
	p.lockInput.Lock()
	defer p.lockInput.Unlock()
	p.componentsNew = p.componentsNew[0:0]
	p.componentsNew,p.componentsTemp = p.componentsTemp,p.componentsNew
}

func (p *ComponentCollection) push(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	if v, ok := p.collection[typ]; ok {
		v.push(com, id)
	} else {
		cd := newComponentData()
		cd.push(com, id)
		p.collection[typ] = cd
	}
}

func (p *ComponentCollection) pop(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	if v, ok := p.collection[typ]; ok {
		v.pop(id)
	}
}

func (p *ComponentCollection) GetComponentsNew() []*CollectionOperateInfo {
	return p.componentsNew
}

func (p *ComponentCollection) GetComponents(com IComponent) []IComponent {
	v, ok := p.collection[reflect.TypeOf(com)]
	if ok {
		return v.data
	}
	return []IComponent{}
}

func (p *ComponentCollection) GetComponent(com IComponent, id uint64) interface{} {
	v, ok := p.collection[reflect.TypeOf(com)]
	if ok {
		if c, ok := v.index[id]; ok {
			return v.data[c]
		}
	}
	return nil
}

func (p *ComponentCollection) GetIterator() *ComponentCollectionIter {
	ls := make([]*componentData, 0)
	for _, value := range p.collection {
		ls = append(ls, value)
	}
	return newComponentCollectionIter(ls)
}
