package ecs

func GetInterestedComponents[T any](s ISystem) *Collection[T] {
	typ := GetType[T]()
	if _, ok := s.Requirements()[typ]; !ok {
		return nil
	}
	return s.World().getComponents(typ).(*Collection[T])
}
