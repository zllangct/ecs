package ecs

import (
	"reflect"
	"unsafe"
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
	return w.getSystem(TypeOf[T]())
}

func GetEntityInfo(world IWorld, entity Entity) (*EntityInfo, bool) {
	return world.GetEntityInfo(entity)
}

func AddFreeComponent[T FreeComponentObject, TP FreeComponentPointer[T]](world IWorld, component *T) {
	world.AddFreeComponent(TP(component))
}

// entity api

func NewEntity(world IWorld) *EntityInfo {
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

func GetInterestedComponent[T ComponentObject](sys ISystem, entity Entity) *T {
	return GetRelatedComponent[T](sys, entity)
}

func GetInterestedComponents[T ComponentObject](sys ISystem) Iterator[T] {
	if sys.getState() == SystemStateInvalid {
		return EmptyIter[T]()
	}
	if !sys.isExecuting() {
		return EmptyIter[T]()
	}
	typ := GetType[T]()
	r, ok := sys.GetRequirements()[typ]
	if !ok {
		return EmptyIter[T]()
	}

	c := sys.World().getComponentSet(typ)
	if c == nil {
		return EmptyIter[T]()
	}
	return NewComponentSetIterator[T](c.(*ComponentSet[T]), r.getPermission() == ComponentReadOnly)
}

func GetRelatedComponent[T ComponentObject](sys ISystem, entity Entity) *T {
	typ := TypeOf[T]()
	isRequire := sys.isRequire(typ)
	if !isRequire {
		return nil
	}
	var cache *ComponentGetter[T]
	cacheMap := sys.getGetterCache()
	c := cacheMap.Get(typ)
	if c != nil {
		cache = (*ComponentGetter[T])(c)
	} else {
		cache = NewComponentGetter[T](sys)
		cacheMap.Add(typ, unsafe.Pointer(cache))
	}
	return cache.Get(entity)
}

func GetIntType(world IWorld, typ reflect.Type) uint16 {
	info := world.GetComponentMetaInfo(typ)
	return info.it
}
