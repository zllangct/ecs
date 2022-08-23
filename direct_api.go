package ecs

import (
	"reflect"
)

// world api

func GetWorldID(world IWorld) int64 {
	return world.GetID()
}

func GetWorldStatus(world IWorld) WorldStatus {
	return world.GetStatus()
}

func RegisterSystem[T SystemObject](world IWorld, order ...Order) {
	world.registerForT(new(T), order...)
}

func GetSystem[T SystemObject](w IWorld) (ISystem, bool) {
	return w.GetSystem(TypeOf[T]())
}

func GetEntityInfo(world IWorld, entity Entity) EntityInfo {
	return world.GetEntityInfo(entity)
}

func AddFreeComponent[T FreeComponentObject, TP FreeComponentPointer[T]](world IWorld, component *T) {
	world.AddFreeComponent(TP(component))
}

// entity api

func NewEntity(world IWorld) EntityInfo {
	return world.(*ecsWorld).NewEntity()
}

func Destroy(world IWorld, entity Entity) {
	world.(*ecsWorld).deleteEntity(entity)
}

func AddComponent(world IWorld, entity Entity, components ...IComponent) {
	for _, com := range components {
		world.addComponent(entity, com)
	}
}

// system api

func NewPeripheralSystem[T PeripheralSystemObject, TP PeripheralSystemPointer[T]]() *T {
	var ins T
	psys := IPeripheralSystem(TP(&ins))
	psys.init()
	if i, ok := psys.(InitReceiver); ok {
		i.Init()
	}
	return &ins
}

func GetInterestedComponents[T ComponentObject](sys ISystem) Iterator[T] {
	if sys.getState() == SystemStateInvalid {
		return EmptyIter[T]()
	}
	if !sys.isExecuting() {
		return EmptyIter[T]()
	}
	typ := GetType[T]()
	r, ok := sys.Requirements()[typ]
	if !ok {
		return EmptyIter[T]()
	}

	c := sys.World().getComponents(typ)
	if c == nil {
		return EmptyIter[T]()
	}
	return NewComponentSetIterator[T](c.(*ComponentSet[T]), r.getPermission() == ComponentReadOnly)
}

func GetRelatedComponent[T ComponentObject](sys ISystem, entity Entity) *T {
	typ := TypeOf[T]()
	_, isRequire := sys.isRequire(typ)
	if !isRequire {
		return nil
	}
	var cache *ComponentGetter[T]
	cacheMap := sys.getGetterCache()
	if v, ok := cacheMap[typ]; ok {
		cache = v.(*ComponentGetter[T])
	} else {
		cacheMap[typ] = NewComponentGetter[T](sys)
	}
	return cache.Get(entity)
}

func GetIntType(typ reflect.Type) uint16 {
	info := ComponentMeta.GetComponentMetaInfo(typ)
	return info.it
}
