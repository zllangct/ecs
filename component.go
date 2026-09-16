package ecs

import (
	"unsafe"

	rockmem "github.com/zllangct/rockmem/golang"
)

type ComponentIntType uint64

type ComponentPointer[T any] interface {
	Component
	*T
}

// ReadOnlyView 约束 codegen 生成的组件只读视图（TReadOnly）。
// 视图自描述全部类型信息：ComponentPacketIdentifier 与源组件的
// PacketIdentifier 一致（用于依赖/组件集查找），FromPtr 从组件内存指针
// 构造视图；因此只读 API 仅需显式书写视图类型一个类型参数。
// 手写组件的只读视图需自行实现这两个方法方可使用 View API。
type ReadOnlyView[TR any] interface {
	ComponentPacketIdentifier() rockmem.PacketIdentifier
	FromPtr(p unsafe.Pointer) TR
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

func GetIntType[T any, TP ComponentPointer[T]]() ComponentIntType {
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
