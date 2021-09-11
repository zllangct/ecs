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
	collections map[reflect.Type]*Collection
	//new component cache
	base  int64
	locks []sync.Mutex
	cTemp []map[reflect.Type][]CollectionOperateInfo
	cNew map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collections: map[reflect.Type]*Collection{},
	}

	for i := 1; ; i++ {
		if c := int64(1 << i); int64(k) < c {
			cc.base = c - 1
			break
		}
	}

	cc.locks = make([]sync.Mutex, cc.base+1)
	cc.cTemp = make([]map[reflect.Type][]CollectionOperateInfo, cc.base+1)
	for index := range cc.cTemp {
		cc.cTemp[index] = make(map[reflect.Type][]CollectionOperateInfo)
		cc.locks[index] = sync.Mutex{}
	}
	return cc
}

func (c *ComponentCollection) TempComponentOperate(entity *Entity, com IComponent, op CollectionOperate) {
	hash := entity.ID() & c.base

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	newOpt := NewCollectionOperateInfo(entity, com, op)
	typ := com.Type()
	b := c.cTemp[hash]
	if _, ok := b[typ]; ok {
		b[typ] = append(b[typ], newOpt)
	} else {
		b[typ] = []CollectionOperateInfo{ newOpt }
	}
}

func (c *ComponentCollection) GetTempFlushTasks() []func() {
	combination := make(map[reflect.Type][]CollectionOperateInfo)

	for i:=0; i < len(c.cTemp); i++ {
		for typ, op := range c.cTemp[i] {
			if _, ok := combination[typ]; ok {
				combination[typ] = append(combination[typ], op...)
			} else {
				combination[typ] = op
			}
		}
	}

	var tasks []func()
	for typ, opList := range combination {
		typTemp := typ
		oopList := opList
		fn := func(){
			for _, operate := range oopList {
				typTemp = operate.com.Type()
				//set owner
				operate.com.setOwner(operate.target)
				//add to component container
				_, ret := c.add(operate.com)
				//add to entity
				operate.target.componentAdded(typTemp, ret)
			}
		}
		tasks = append(tasks, fn)
	}
	return tasks
}

// TempFlush handle and flush new components,should be called before destroy period
func (c *ComponentCollection) TempFlush() {
	//var temp []CollectionOperateInfo
	//for index, item := range c.cTemp {
	//	c.locks[index].Lock()
	//	temp = append(temp, item...)
	//	c.cTemp[index] = c.cTemp[index][0:0]
	//	c.locks[index].Unlock()
	//}
	//tempNew := map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo{
	//	COLLECTION_OPERATE_ADD:    make(map[reflect.Type][]CollectionOperateInfo),
	//	COLLECTION_OPERATE_DELETE: make(map[reflect.Type][]CollectionOperateInfo),
	//}
	//for _, operate := range temp {
	//	typ := operate.com.Type()
	//	//set owner
	//	operate.com.setOwner(operate.target)
	//	//add to component container
	//	_, ret := c.add(operate.com)
	//	//add to entity
	//	operate.target.componentAdded(typ, ret)
	//	//add to new component list
	//	if _, ok := tempNew[operate.op][typ]; !ok {
	//		tempNew[operate.op][typ] = make([]CollectionOperateInfo, 0)
	//	}
	//	tempNew[operate.op][typ] = append(tempNew[operate.op][typ], operate)
	//}
}

func (c *ComponentCollection) add(com IComponent) (int64, IComponent) {
	typ := com.Type()
	var collection *Collection
	if v, ok := c.collections[typ]; ok {
		nc := NewCollection(int(typ.Size()))
		c.collections[typ] = nc
	} else {
		collection = v
	}
	i := (*iface)(unsafe.Pointer(&com))
	id, ptr := collection.Add(i.data)
	i.data = ptr
	com.setID(id)
	return id, com
}

func (p *ComponentCollection) GetNewComponentsAll() []CollectionOperateInfo {
	size := 0
	for _, m := range p.cNew {
		for _, mm := range m {
			size += len(mm)
		}
	}
	temp := make([]CollectionOperateInfo, 0, size)
	for _, m := range p.cNew {
		for _, mm := range m {
			temp = append(temp, mm...)
		}
	}
	return temp
}

func (p *ComponentCollection) GetNewComponents(op CollectionOperate, typ reflect.Type) []CollectionOperateInfo {
	return p.cNew[op][typ]
}

// GetAllComponents TODO need to refactor
func (c *ComponentCollection) GetAllComponents() ComponentCollectionIter {
	length := 0
	for _, value := range c.collections {
		length += value.Len()
	}
	components := make([]*Collection, 0, length)
	index := 0
	for _, value := range c.collections {
		l := value.Len()
		components = append(components, value)
		index += l
	}
	return NewComponentCollectionIter(components)
}

func (c *ComponentCollection) GetComponent(id int64) unsafe.Pointer {
	//var ins T
	//v, ok := c.collections[reflect.TypeOf(ins)]
	//if ok {
	//	if c := v.(IndexedCollection[T]).Get(id); c != nil {
	//		return c
	//	}
	//	return nil
	//}
	return nil
}

