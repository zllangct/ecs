package ecs

import (
	"reflect"
	"unsafe"
)

func RegisterSystem[T SystemObject](world IWorld, order ...Order) {
	world.registerForT(new(T), order...)
}

func AddFreeComponent[T FreeComponentObject, TP FreeComponentPointer[T]](world IWorld, component *T) {
	world.addFreeComponent(TP(component))
}

func GetComponent[T ComponentObject](sys ISystem, entity Entity) *T {
	return GetRelated[T](sys, entity)
}

func GetComponentAll[T ComponentObject](sys ISystem) Iterator[T] {
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

func GetRelated[T ComponentObject](sys ISystem, entity Entity) *T {
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

func BindUtility[T UtilityObject, TP UtilityPointer[T]](si SystemInitConstraint) {
	if si.isValid() {
		panic("out of initialization stage")
	}
	utility := TP(new(T))
	sys := si.getSystem()
	utility.setSystem(sys)
	utility.setWorld(sys.World())
	sys.setUtility(utility)
	sys.World().base().utilities[utility.Type()] = utility
}

func GetUtility[T UtilityObject](getter IUtilityGetter) (*T, bool) {
	w := getter.getWorld()
	if w == nil {
		return nil, false
	}
	u, ok := w.getUtilityForT(TypeOf[T]())
	if !ok {
		return nil, false
	}
	return (*T)(u), true
}

func TypeOf[T any]() reflect.Type {
	ins := (*T)(nil)
	return reflect.TypeOf(ins).Elem()
}
