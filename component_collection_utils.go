package ecs

//import (
//	"reflect"
//	"unsafe"
//)
//
//func TempComponentOperate[T IComponent](c *ComponentCollection, entity *Entity, com *T, op CollectionOperate) {
//	hash := entity.ID() & c.base
//
//	c.locks[hash].Lock()
//	defer c.locks[hash].Unlock()
//
//	newOpt := NewCollectionOperateInfo(entity, com, op)
//	typ := com.Type()
//	b := c.cTemp[hash]
//	if _, ok := b[typ]; ok {
//		b[typ] = append(b[typ], newOpt)
//	} else {
//		b[typ] = []CollectionOperateInfo{ newOpt }
//	}
//}
//
//func Add[T IComponent](cc *ComponentCollection, com *T, id int64) *T {
//	var c *Collection
//	var ins T
//	typ := reflect.TypeOf(ins)
//	if v, ok := cc.collections[typ]; ok {
//		v = NewCollection(int(unsafe.Sizeof(ins)))
//		cc.collections[typ] = v
//	} else {
//		c = v
//	}
//	_, ptr := c.Add(unsafe.Pointer(com))
//	return (*T)(ptr)
//}
//
//func Remove[T IComponent](c *ComponentCollection, id int64) {
//	var ins T
//	typ := reflect.TypeOf(ins)
//	if v, ok := c.collections[typ]; ok {
//		v.Remove(id)
//	}
//}
//
//func GetNewComponentsAll(c *ComponentCollection) []CollectionOperateInfo {
//	return nil
//}
//
//func GetNewComponents[T IComponent](c *ComponentCollection, op CollectionOperate) []CollectionOperateInfo {
//	var ins T
//	typ := reflect.TypeOf(ins)
//	_=typ
//	return nil
//}

//func GetComponents[T IComponent](cc *ComponentCollection) *iterator {
//	//var ins T
//	//v, ok := cc.collections[reflect.TypeOf(ins)]
//	//if ok {
//	//	return v.(IndexedCollection[T]).GetIterator()
//	//}
//	//return EmptyIterator()
//	return nil
//}