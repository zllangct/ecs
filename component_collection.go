package ecs

import (
	"reflect"
	"sync"
	"unsafe"
)

type CollectionOperate uint8

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
	collection map[reflect.Type]*ContainerWithId
	//new component cache
	base           uint64
	locks          []sync.Mutex
	componentsTemp [][]CollectionOperateInfo
	componentsNew  map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collection:    map[reflect.Type]*ContainerWithId{},
		componentsNew: make(map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo),
	}

	for i := 1; ; i++ {
		if c := uint64(1 << i); uint64(k) < c {
			cc.base = c - 1
			break
		}
	}

	cc.locks = make([]sync.Mutex, cc.base+1)
	cc.componentsTemp = make([][]CollectionOperateInfo, cc.base+1)
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
	tempNew := map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo{
		COLLECTION_OPERATE_ADD:    make(map[reflect.Type][]CollectionOperateInfo),
		COLLECTION_OPERATE_DELETE: make(map[reflect.Type][]CollectionOperateInfo),
	}
	for _, operate := range temp {
		typ := reflect.TypeOf(operate.com)
		//add to component container
		p.push(typ, operate.com, operate.com.GetOwner().ID())
		//add to new component list
		if _, ok := tempNew[operate.op][typ]; !ok {
			tempNew[operate.op][typ] = make([]CollectionOperateInfo, 0)
		}
		tempNew[operate.op][typ] = append(tempNew[operate.op][typ], operate)
	}
	p.componentsNew = tempNew

}

func (p *ComponentCollection) Push(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	p.push(typ, com, id)
}

func (p *ComponentCollection) push(typ reflect.Type, com IComponent, id uint64) unsafe.Pointer {
	efaceStruct := (*eface)(unsafe.Pointer(&com))
	var v *ContainerWithId
	var ok bool
	v, ok = p.collection[typ]
	if !ok {
		v = NewContainerWithId(typ.Size())
		p.collection[typ] = v
	}
	_, pointer := v.Add(unsafe.Pointer(*(**int)(efaceStruct.data)))
	//TODO 检查正确性、检查是否造成外部内存泄漏
	*(**[]byte)(efaceStruct.data) = (*[]byte)(pointer)
	return pointer
}

func (p *ComponentCollection) Pop(com IComponent, id uint64) {
	typ := reflect.TypeOf(com)
	if v, ok := p.collection[typ]; ok {
		v.RemoveById(id)
	}
}

func (p *ComponentCollection) GetNewComponentsAll() []CollectionOperateInfo {
	size := 0
	for _, m := range p.componentsNew {
		for _, mm := range m {
			size += len(mm)
		}
	}
	temp := make([]CollectionOperateInfo, 0, size)
	for _, m := range p.componentsNew {
		for _, mm := range m {
			temp = append(temp, mm...)
		}
	}
	return temp
}

func (p *ComponentCollection) GetNewComponents(op CollectionOperate, typ reflect.Type) []CollectionOperateInfo {
	return p.componentsNew[op][typ]
}

func (p *ComponentCollection) GetComponents(com IComponent) *Iterator {
	v, ok := p.collection[reflect.TypeOf(com)]
	if ok {
		return v.GetIterator()
	}
	return EmptyIterator()
}

func (p *ComponentCollection) GetAllComponents() *ComponentCollectionIter {
	length := 0
	for _, value := range p.collection {
		length += value.Len()
	}
	components := make([]*ContainerWithId, 0, length)
	index := 0
	for _, value := range p.collection {
		l := value.Len()
		components = append(components, value)
		index += l
	}
	return NewComponentCollectionIter(components)
}

func (p *ComponentCollection) GetComponent(com IComponent, id uint64) IComponent {
	v, ok := p.collection[reflect.TypeOf(com)]
	if ok {
		if c := v.GetById(id); c != nil {
			return com
		}
	}
	return nil
}

func (p *ComponentCollection) GetIterator() *ComponentCollectionIter {
	ls := make([]*ContainerWithId, len(p.collection))
	i := 0
	for _, value := range p.collection {
		ls[i] = value
		i += 1
	}
	return NewComponentCollectionIter(ls)
}
