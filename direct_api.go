package ecs

import (
	"reflect"
	"unsafe"
)

func RegisterSystem[T SystemObject](world iWorldBase, order ...Order) {
	world.registerForT(new(T), order...)
}

func AddFreeComponent[T FreeComponentObject, TP FreeComponentPointer[T]](world iWorldBase, component *T) {
	world.AddFreeComponent(TP(component))
}

func GetInterestedComponent[T ComponentObject](sys ISystem, entity Entity) *T {
	return GetRelatedComponent[T](sys, entity)
}

func BindGate[T GateObject](world *AsyncWorld) {
	world.BindGate(any(new(T)).(IGate))
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

func GetGate[T GateObject](world *AsyncWorld) *T {
	gate := any(world.GetGate())
	g, ok := gate.(*T)
	if ok {
		return nil
	}
	return g
}

func TypeOf[T any]() reflect.Type {
	ins := (*T)(nil)
	return reflect.TypeOf(ins).Elem()
}
