package ecs

import (
	"reflect"
	runtime2 "runtime"
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

type ComponentOperate = CollectionOperate

const (
	COMPONENT_OPERATE_ADD    = COLLECTION_OPERATE_ADD    //add component operation
	COMPONENT_OPERATE_DELETE = COLLECTION_OPERATE_DELETE //delete component operation
)

type CollectionOperateInfo struct {
	com IComponent
	op  CollectionOperate
}

func NewCollectionOperateInfo(com IComponent, op CollectionOperate) CollectionOperateInfo {
	return CollectionOperateInfo{com: com, op: op}
}

type ComponentCollection struct {
	collection map[reflect.Type]*componentData
	//new component cache
	base           uint64
	locks          []sync.Mutex
	componentsTemp [][]CollectionOperateInfo
	componentsNew  map[CollectionOperate][]CollectionOperateInfo
}

func NewComponentCollection() *ComponentCollection {
	cc := &ComponentCollection{
		collection:    map[reflect.Type]*componentData{},
		componentsNew: make(map[CollectionOperate][]CollectionOperateInfo),
	}

	numCpu := runtime2.NumCPU()

	for i := 1; ; i++ {
		if c := uint64(1 << i); uint64(numCpu*4) < c {
			cc.base = c - 1
			break
		}
	}

	cc.locks = make([]sync.Mutex, cc.base)
	cc.componentsTemp = make([][]CollectionOperateInfo, cc.base)
	for index := range cc.componentsTemp {
		cc.componentsTemp[index] = make([]CollectionOperateInfo, 0)
		cc.locks[index] = sync.Mutex{}
	}
	return cc
}

//new component temp
func (p *ComponentCollection) TempComponentOperate(com IComponent, op CollectionOperate) {
	hash := com.GetOwner().ID() & p.base
	p.locks[hash].Lock()
	p.componentsTemp[hash] = append(p.componentsTemp[hash], NewCollectionOperateInfo(com, op))
	p.locks[hash].Unlock()
}

//handle and flush new components,should be called before destroy period
func (p *ComponentCollection) TempFlush() {
	var temp []CollectionOperateInfo
	for index, item := range p.componentsTemp {
		p.locks[index].Lock()
		temp = append(temp, item...)
		p.componentsTemp[index] = p.componentsTemp[index][0:0]
		p.locks[index].Unlock()
	}
	tempNew := map[CollectionOperate][]CollectionOperateInfo{}
	for _, operate := range temp{
		if _, ok := tempNew[operate.op]; !ok {
			tempNew[operate.op] = make([]CollectionOperateInfo, 0)
		}
		tempNew[operate.op] = append(tempNew[operate.op], operate)
	}
	p.componentsNew = tempNew

}

func (p *ComponentCollection) Push(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	if v, ok := p.collection[typ]; ok {
		v.push(com, id)
	} else {
		cd := newComponentData()
		cd.push(com, id)
		p.collection[typ] = cd
	}
}

func (p *ComponentCollection) Pop(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	if v, ok := p.collection[typ]; ok {
		v.pop(id)
	}
}

func (p *ComponentCollection) GetComponentsNew() map[CollectionOperate][]CollectionOperateInfo {
	return p.componentsNew
}

func (p *ComponentCollection) GetComponents(com IComponent) []IComponent {
	v, ok := p.collection[reflect.TypeOf(com)]
	if ok {
		return v.data
	}
	return []IComponent{}
}

func (p *ComponentCollection) GetAllComponents() []IComponent {
	length := 0
	for _, value := range p.collection {
		length += len(value.data)
	}
	components := make([]IComponent, length)
	index := 0
	for _, value := range p.collection {
		l := len(value.data)
		copy(components[index:index+l], value.data)
		index += l
	}
	return components
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
	ls := make([]*componentData, len(p.collection))
	i := 0
	for _, value := range p.collection {
		ls[i] = value
		i += 1
	}
	return NewComponentCollectionIter(ls)
}
