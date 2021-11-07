package ecs

import "unsafe"

// runtime api

func RuntimeConfigure(config *RuntimeConfig) {
	Runtime.Configure(config)
}

func Run() {
	Runtime.run()
}

func Stop() {
	Runtime.stop()
}

func CreateWorld(config *WorldConfig) IWorld {
	return Runtime.newWorld(config)
}

func DestroyWorld(world IWorld) {
	Runtime.destroyWorld(world.(*ecsWorld))
}

func AddJob(job func(), hashKey ...uint32) {
	Runtime.addJob(job, hashKey...)
}

// world api

func WorldRun(world IWorld) {
	world.Run()
}

func GetWorldID(world IWorld) int64 {
	return world.GetID()
}

func GetWorldStatus(world IWorld) WorldStatus {
	return world.GetStatus()
}

func RegisterSystem[T ISystemTemplate](world IWorld, order ...Order) {
	world.(*ecsWorld).registerForT(new(T), order...)
}

func GetSystem[T ISystem](w IWorld) (ISystem, bool) {
	return w.getSystem(TypeOf[T]())
}

func GetEntityInfo(world IWorld, entity Entity) *EntityInfo {
	return world.getEntityInfo(entity)
}

// entity api

func NewEntity(world IWorld) *EntityInfo {
	return newEntityInfo(world.(*ecsWorld))
}

func EntityDestroyByInfo(info *EntityInfo) {
	info.Destroy()
}

func EntityDestroy(world IWorld, entity Entity) {
	info := world.getEntityInfo(entity)
	if info != nil {
		info.Destroy()
	}
}

// system api

func AddRequireComponent[T IComponent](sys ISystem) {
	sys.setRequirementsByType(TypeOf[T]())
}

func GetInterestedComponents[T any](sys ISystem) Iterator[T] {
	typ := GetType[T]()
	if _, ok := sys.Requirements()[typ]; !ok {
		Log.Error("not require, typ:", typ)
		return nil
	}
	if sys.World() == nil {
		Log.Error("world is nil")
	}
	c := sys.World().getComponents(typ)
	if c == nil {
		return nil
	}
	return NewIterator(c.(*Collection[T]))
}

func CheckComponent[T any](sys ISystem, entity *EntityInfo) *T {
	typ := TypeOf[T]()
	isRequire := sys.isRequire(typ)
	if !isRequire {
		return nil
	}
	c := entity.getComponentByType(typ)
	return (*T)(unsafe.Pointer((*iface)(unsafe.Pointer(&c)).data))
}

