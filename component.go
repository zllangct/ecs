package ecs

import rockmem "github.com/zllangct/rockmem/golang"

type ComponentIntType uint64

type ComponentPointer[T ComponentObject] interface {
	Component
	*T
}

type ComponentObject interface {
	ComponentObjectIdentifier()
}

type Component interface {
	NewComponentSet() ComponentSet
	PacketIdentifier() rockmem.PacketIdentifier
	IsNomadic() bool
	IsDisposable() bool
}

func GetIntTypeByComp(com Component) ComponentIntType {
	return ComponentIntType(com.PacketIdentifier())
}

func GetIntType[T ComponentObject, TP ComponentPointer[T]]() ComponentIntType {
	return ComponentIntType(TP(nil).PacketIdentifier())
}

type dummyComponent struct {
	Seq int32
}

func (d *dummyComponent) NewComponentSet() ComponentSet {
	return NewCSet[dummyComponent]()
}

func (d *dummyComponent) PacketIdentifier() rockmem.PacketIdentifier {
	return 65530
}

func (d *dummyComponent) IsNomadic() bool {
	return false
}

func (d *dummyComponent) IsDisposable() bool {
	return false
}

func (d dummyComponent) ComponentObjectIdentifier() {}
