package ecs

type Writable = bool

const (
	ReadWrite Writable = true
	ReadOnly  Writable = false
)

type ComponentDependency interface {
	readonly() bool
	intType() ComponentIntType
}

type Dependency uint32

func NewDependency(it ComponentIntType, writable ...Writable) Dependency {
	w := ReadOnly
	if len(writable) > 0 {
		w = writable[0]
	}
	if w {
		return Dependency(it) | (1 << 16)
	} else {
		return Dependency(it)
	}
}

func (r Dependency) readonly() bool {
	return r>>16&1 == 0
}

func (r Dependency) intType() ComponentIntType {
	return ComponentIntType(r)
}

type readonly[T ComponentObject, TP ComponentPointer[T]] struct{}

func (r readonly[T, TP]) readonly() bool {
	return true
}

func (r readonly[T, TP]) intType() ComponentIntType {
	return GetIntType[T, TP]()
}

type readwrite[T ComponentObject, TP ComponentPointer[T]] struct{}

func (r readwrite[T, TP]) readonly() bool {
	return false
}

func (r readwrite[T, TP]) intType() ComponentIntType {
	return GetIntType[T, TP]()
}

func Dep2[T ComponentObject, TP ComponentPointer[T]](writable ...Writable) ComponentDependency {
	ro := ReadOnly
	if len(writable) > 0 {
		ro = writable[0]
	}
	if ro {
		return &readwrite[T, TP]{}
	} else {
		return &readonly[T, TP]{}
	}
}

func Dep[T ComponentObject, TP ComponentPointer[T]](writable ...Writable) Dependency {
	it := GetIntType[T, TP]()
	return NewDependency(it, writable...)
}
