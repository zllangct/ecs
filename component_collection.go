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
	target *Entity
	com    IComponent
	op     CollectionOperate
}

func NewCollectionOperateInfo(entity *Entity, com IComponent, op CollectionOperate) CollectionOperateInfo {
	return CollectionOperateInfo{target: entity, com: com, op: op}
}

type ComponentCollection struct {
	collection map[reflect.Type]interface{}
	//new component cache
	base           uint64
	locks          []sync.Mutex
	componentsTemp [][]CollectionOperateInfo
	componentsNew  map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collection:    map[reflect.Type]interface{}{},
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
func (p *ComponentCollection) TempComponentOperate(entity *Entity, com IComponent, op CollectionOperate) {
	hash := entity.ID() & p.base
	p.locks[hash].Lock()
	p.componentsTemp[hash] = append(p.componentsTemp[hash], NewCollectionOperateInfo(entity, com, op))
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
		typ := operate.com.GetRealType()
		//set component owner
		operate.com.setOwner(operate.target)
		//add to component container
		ret := Push(p, typ, operate.com, operate.target.ID())
		//add to entity
		operate.target.componentAdded(typ, ret)

		//add to new component list
		if _, ok := tempNew[operate.op][typ]; !ok {
			tempNew[operate.op][typ] = make([]CollectionOperateInfo, 0)
		}
		tempNew[operate.op][typ] = append(tempNew[operate.op][typ], operate)
	}
	p.componentsNew = tempNew

}

func Push[T IComponent](c *ComponentCollection, com *T, id int64) IComponent {
	typ := reflect.TypeOf(*com)
	ifaceStruct := (*iface)(unsafe.Pointer(com))
	var v *ContainerWithId
	var ok bool
	v, ok = c.collection[typ]
	if !ok {
		v = NewContainerWithIdByte[T]()
		c.collection[typ] = v
	}
	_, ptr := v.Add(com, id)
	ifaceStruct.data = unsafe.Pointer(ptr)
	return com
}

func Pop[T IComponent](c *ComponentCollection, id int64) {
	var ins T
	typ := reflect.TypeOf(ins)
	if v, ok := c.collection[typ]; ok {
		v.(ContainerWithId[T]).RemoveById(id)
	}
}

func GetNewComponentsAll(c *ComponentCollection) []CollectionOperateInfo {
	size := 0
	for _, m := range c.componentsNew {
		for _, mm := range m {
			size += len(mm)
		}
	}
	temp := make([]CollectionOperateInfo, 0, size)
	for _, m := range c.componentsNew {
		for _, mm := range m {
			temp = append(temp, mm...)
		}
	}
	return temp
}

func GetNewComponents[T IComponent](c *ComponentCollection, op CollectionOperate) []CollectionOperateInfo {
	var ins T
	typ := reflect.TypeOf(ins)
	return c.componentsNew[op][typ]
}

func (p *ComponentCollection) GetComponents[T IComponent]() *iterator {
	var ins T
	v, ok := p.collection[reflect.TypeOf(ins)]
	if ok {
		return v.(ContainerWithId[T]).GetIterator()
	}
	return EmptyIterator()
}

//TODO need to refactor
func (p *ComponentCollection) GetAllComponents() ComponentCollectionIter {
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

func (p *ComponentCollection) GetComponent[T IComponent](id int64) unsafe.Pointer {
	var ins T
	v, ok := p.collection[reflect.TypeOf(ins)]
	if ok {
		if c := v.(ContainerWithId[T]).GetById(id); c != nil {
			return c
		}
		return nil
	}
	return nil
}

//TODO need to refactor
func (p *ComponentCollection) GetIterator() *componentCollectionIter {
	ls := make([]*ContainerWithId, len(p.collection))
	i := 0
	for _, value := range p.collection {
		ls[i] = value
		i += 1
	}
	return NewComponentCollectionIter(ls)
}
