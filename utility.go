package ecs

type IUtility[T SystemObject] interface {
}

type Utility[T SystemObject] struct {
	w   IWorld
	sys ISystem
}

func (u *Utility[T]) getSystem() ISystem {
	if u.sys == nil {
		s, ok := u.w.GetSystem(TypeOf[T]())
		if ok {
			u.sys = s
		}
	}
	return u.sys
}

func (u *Utility[T]) Emit(event CustomEventName, args ...interface{}) {
	sys := u.getSystem()
	sys.Emit(event, args...)
}

func GetUtility[T SystemObject](w IWorld) (IUtility[T], bool) {
	sys, ok := w.GetSystem(TypeOf[T]())
	if !ok {
		return nil, false
	}
	return (*System[T])(sys.getPointer()).GetUtility(), ok
}

// Test code

type TestUtilitySystem struct {
	System[TestUtilitySystem]
}

type TestUtility struct {
	Utility[TestUtilitySystem]
}

func (t *TestUtility) Hello() {
}
