package ecs

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
