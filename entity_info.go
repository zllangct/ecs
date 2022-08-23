package ecs

import (
	"reflect"
)

type TypeList []reflect.Type

func newTypeList(cap int) TypeList {
	return make(TypeList, 0, cap)
}

func (tl *TypeList) Contains(t reflect.Type) bool {
	for _, t2 := range *tl {
		if t2 == t {
			return true
		}
	}
	return false
}

func (tl *TypeList) Find(t reflect.Type) (int, bool) {
	for i, t2 := range *tl {
		if t2 == t {
			return i, true
		}
	}
	return 0, false
}

func (tl *TypeList) Remove(t reflect.Type) {
	i, ok := tl.Find(t)
	if !ok {
		return
	}
	//*tl = append((*tl)[:i], (*tl)[i+1:]...)
	(*tl)[i], (*tl)[len(*tl)-1] = (*tl)[len(*tl)-1], (*tl)[i]
	*tl = (*tl)[:len(*tl)-1]
}

func (tl *TypeList) Append(t ...reflect.Type) {
	/* todo
	大量分配对象, TypeList 频繁修改, 考虑链表, map的插入效率,find、delete效率高
	*/
	*tl = append(*tl, t...)
}

// 干掉entity info，没有太大的存在意义，components的icomponent会随着set中slice的扩容地址发生变化，导致错误，不能持久引用，也不能持久引用索引
// set中删除操作会改变索引。在访问兄弟组件时，如果是内存整理完整状态，可以直接试探访问下一个元素，如果不是，则需要通过id重新通过Set的Get方法获取，
// 访问优化应该置于迭代器中
type EntityInfo struct {
	world    *ecsWorld
	entity   Entity
	compound Compound
}

func (e *EntityInfo) init(world *ecsWorld) *EntityInfo {
	info := &EntityInfo{
		world:  world,
		entity: newEntity(),
	}
	world.addEntity(info.entity)
	return info
}

func (e *EntityInfo) Destroy() {
	//TODO 删除entity
	//e.world.deleteEntity(e.entity)
}

func (e *EntityInfo) Entity() Entity {
	return e.entity
}

func (e *EntityInfo) hashKey() int64 {
	return int64(e.entity)
}

func (e *EntityInfo) Add(components ...IComponent) []error {
	for _, c := range components {
		e.world.addComponent(e.entity, c)
	}
	return nil
}

func (e *EntityInfo) Remove(components ...IComponent) {
	for _, c := range components {
		e.world.deleteComponent(e.entity, c)
	}
}
