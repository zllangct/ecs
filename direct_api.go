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

func CreateWorld(config *WorldConfig) *World {
	return Runtime.newWorld(config)
}

func DestroyWorld(world *World) {
	Runtime.destroyWorld(world)
}

func AddJob(job func(), hashKey ...uint32) {
	Runtime.addJob(job, hashKey...)
}

// world api

func WorldRun(world *World) {
	world.Run()
}

func GetWorldID(world *World) int64 {
	return world.GetID()
}

func GetWorldStatus(world *World) WorldStatus {
	return world.GetStatus()
}

func RegisterSystem[T ISystemTemplate](world *World) {
	world.registerForT(new(T))
}

func GetSystem[T ISystem](w *World) (ISystem, bool) {
	return w.GetSystem(TypeOf[T]())
}

func NewEntity(world *World) *EntityInfo {
	return newEntityInfo(world)
}

func GetEntity(world *World, entity Entity) *EntityInfo {
	return world.getEntityInfo(entity)
}

// entity api

func GetEntityInfo(world *World, entity Entity) *EntityInfo {
	return world.getEntityInfo(entity)
}

func EntityDestroy(entity *EntityInfo) {
	entity.Destroy()
}

func AddComponent[T Component[T]](entity *EntityInfo) {

}

// system api

func GetInterestedComponents[T any](s ISystem) *Collection[T] {
	typ := GetType[T]()
	if _, ok := s.Requirements()[typ]; !ok {
		Log.Error("not require, typ:", typ)
		return nil
	}
	if s.World() == nil {
		Log.Error("world is nil")
	}
	c := s.World().getComponents(typ)
	if c == nil {
		return nil
	}
	return c.(*Collection[T])
}

func CheckComponent[T any](s ISystem, entity *EntityInfo) *T {
	c := entity.getComponentByType(TypeOf[T]())
	return (*T)(unsafe.Pointer((*iface)(unsafe.Pointer(&c)).data))
}


