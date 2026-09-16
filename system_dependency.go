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
		return ItDependency{it, !writable[0]}
	}

	return ItDependency{it, false}
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

// Dep 声明组件依赖，默认可写（ReadWrite）；Dependency 底层 bool 存的是 isReadonly。
func Dep[T any, TP ComponentPointer[T]](writable ...Writable) ComponentDependency {
	w := ReadWrite
	if len(writable) > 0 {
		w = writable[0]
	}
	return (Dependency[T, TP])(!w)
}
