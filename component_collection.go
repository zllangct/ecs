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
	base  int64
	locks []sync.Mutex
	cTemp []map[reflect.Type][]func()
}

func NewComponentCollection(k int) *ComponentCollection {
	cc := &ComponentCollection{
		collection: map[reflect.Type]interface{}{},
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

func TempComponentOperate[T IComponent](c *ComponentCollection, entity *Entity, com *T, op CollectionOperate) {
	hash := entity.ID() & c.base

	c.locks[hash].Lock()
	defer c.locks[hash].Unlock()

	newOpt := NewCollectionOperateInfo(entity, com, op)
	typ := com.GetType()
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
		ttyp := typ
		oopList := opList
		fn := func(){
			for _, operate := range oopList {
				//set component owner
				operate.com.setOwner(operate.target)
				//add to component container
				ret := Add(c, operate.com, operate.target.ID())
				//add to entity
				operate.target.componentAdded(typ, ret)

				//add to new component list
				if _, ok := tempNew[operate.op][typ]; !ok {
					tempNew[operate.op][typ] = make([]CollectionOperateInfo, 0)
				}
				tempNew[operate.op][typ] = append(tempNew[operate.op][typ], operate)

			}
		}
		tasks = append(tasks, fn)
	}
	return tasks
}

//handle and flush new components,should be called before destroy period
func (c *ComponentCollection) TempFlush() {
	var temp []CollectionOperateInfo
	for index, item := range c.cTemp {
		c.locks[index].Lock()
		temp = append(temp, item...)
		c.cTemp[index] = c.cTemp[index][0:0]
		c.locks[index].Unlock()
	}
	tempNew := map[CollectionOperate]map[reflect.Type][]CollectionOperateInfo{
		COLLECTION_OPERATE_ADD:    make(map[reflect.Type][]CollectionOperateInfo),
		COLLECTION_OPERATE_DELETE: make(map[reflect.Type][]CollectionOperateInfo),
	}
	for _, operate := range temp {
		typ := operate.com.GetType()
		//set component owner
		operate.com.setOwner(operate.target)
		//add to component container
		ret := (c, operate.com, operate.target.ID())
		//add to entity
		operate.target.componentAdded(typ, ret)

		//add to new component list
		if _, ok := tempNew[operate.op][typ]; !ok {
			tempNew[operate.op][typ] = make([]CollectionOperateInfo, 0)
		}
		tempNew[operate.op][typ] = append(tempNew[operate.op][typ], operate)
	}
}

func Add[T IComponent](cc *ComponentCollection, com *T, id int64) *T {
	var c *IndexedCollection[T]
	var ins T
	typ := reflect.TypeOf(ins)
	if v, ok := cc.collection[typ]; ok {
		v = NewContainerWithId[T]()
		cc.collection[typ] = v
	} else {
		c = v.(*IndexedCollection[T])
	}
	_, ptr := c.Add(com, id)
	return ptr
}

func Remove[T IComponent](c *ComponentCollection, id int64) {
	var ins T
	typ := reflect.TypeOf(ins)
	if v, ok := c.collection[typ]; ok {
		v.(*IndexedCollection[T]).Remove(id)
	}
}

func GetNewComponentsAll(c *ComponentCollection) []CollectionOperateInfo {
	return nil
}

func GetNewComponents[T IComponent](c *ComponentCollection, op CollectionOperate) []CollectionOperateInfo {
	var ins T
	typ := reflect.TypeOf(ins)
	_=typ
	return nil
}

//func GetComponents[T IComponent](cc *ComponentCollection) *iterator {
//	//var ins T
//	//v, ok := cc.collection[reflect.TypeOf(ins)]
//	//if ok {
//	//	return v.(IndexedCollection[T]).GetIterator()
//	//}
//	//return EmptyIterator()
//	return nil
//}

//TODO need to refactor
func (c *ComponentCollection) GetAllComponents() ComponentCollectionIter {
	//length := 0
	//for _, value := range c.collection {
	//	length += value.Len()
	//}
	//components := make([]*IndexedCollection, 0, length)
	//index := 0
	//for _, value := range c.collection {
	//	l := value.Len()
	//	components = append(components, value)
	//	index += l
	//}
	//return NewComponentCollectionIter(components)
	return nil
}

func (c *ComponentCollection) GetComponent(id int64) unsafe.Pointer {
	//var ins T
	//v, ok := c.collection[reflect.TypeOf(ins)]
	//if ok {
	//	if c := v.(IndexedCollection[T]).Get(id); c != nil {
	//		return c
	//	}
	//	return nil
	//}
	return nil
}

//TODO need to refactor
//func (p *ComponentCollection) GetIterator() *componentCollectionIter {
//	//ls := make([]*IndexedCollection, len(p.collection))
//	//i := 0
//	//for _, value := range p.collection {
//	//	ls[i] = value
//	//	i += 1
//	//}
//	//return NewComponentCollectionIter(ls)
//	return nil
//}
