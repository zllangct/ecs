package ecs

type ComponentIntType uint16

type ComponentPointer[T ComponentObject] interface {
	Component
	*T
}

type ComponentObject interface {
	ComponentObjectIdentifier()
}

type Component interface {
	NewComponentSet() ComponentSet
	GetComponentSeq() int32
	IsNomadic() bool
	IsDisposable() bool
}

func GetIntTypeByComp(com Component) ComponentIntType {
	return ComponentIntType(com.GetComponentSeq())
}

func GetIntType[T ComponentObject, TP ComponentPointer[T]]() ComponentIntType {
	return ComponentIntType(TP(nil).GetComponentSeq())
}

type dummyComponent struct {
	Seq int32
}

func (d *dummyComponent) NewComponentSet() ComponentSet {
	return NewCSet[dummyComponent]()
}

func (d *dummyComponent) GetComponentSeq() int32 {
	return 65530
}

func (d *dummyComponent) IsNomadic() bool {
	return false
}

func (d *dummyComponent) IsDisposable() bool {
	return false
}

func (d dummyComponent) ComponentObjectIdentifier() {}
