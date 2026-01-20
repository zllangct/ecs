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

type ItDependency struct {
	componentIntType ComponentIntType
	isReadonly       Writable
}

func NewItDependency(it ComponentIntType, writable ...Writable) ItDependency {
	if len(writable) > 0 {
		return ItDependency{it, writable[0]}
	}

	return ItDependency{it, ReadOnly}
}

func (r ItDependency) readonly() bool {
	return r.isReadonly
}

func (r ItDependency) intType() ComponentIntType {
	return r.componentIntType
}

type Dependency[T any, TP ComponentPointer[T]] bool

func (r Dependency[T, TP]) readonly() bool {
	return bool(r)
}

func (r Dependency[T, TP]) intType() ComponentIntType {
	return GetIntType[T, TP]()
}

func Dep[T any, TP ComponentPointer[T]](writable ...Writable) ComponentDependency {
	ro := ReadOnly
	if len(writable) > 0 {
		ro = writable[0]
	}
	return (Dependency[T, TP])(!ro)
}
