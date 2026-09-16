package ecs

import (
	"fmt"
	"sync"
	"unsafe"
)

// ComponentSetFactory is a function that creates a new ComponentSet
// and optionally unmarshals data into it
type ComponentSetFactory func() ComponentSet

// ComponentSetUnmarshaler is a function that unmarshals data into a ComponentSet
type ComponentSetUnmarshaler func(data *SerializableSparseArrayData) ComponentSet

type ComponentDefineInfo struct {
	Group       string
	Type        ComponentIntType
	Factory     ComponentSetFactory
	Unmarshaler ComponentSetUnmarshaler
	// Size 组件定长字节数（unsafe.Sizeof(T)），组表列 stride 用
	Size int
	// Proto 组件原型构造器，注册期读取 IsDisposable/IsNomadic 等标记
	Proto func() Component
}

// ComponentRegistry manages component type registrations for serialization/deserialization
type ComponentRegistry struct {
	mu    sync.RWMutex
	infos map[ComponentIntType]*ComponentDefineInfo
}

// Global component registry instance
var globalComponentRegistry = &ComponentRegistry{
	infos: make(map[ComponentIntType]*ComponentDefineInfo),
}

// GetComponentRegistry returns the global component registry
func GetComponentRegistry() *ComponentRegistry {
	return globalComponentRegistry
}

// Register registers a component type with its factory and unmarshaler.
// 设计约定：组件类型进程级全局唯一（rockgo 单 world 场景）；
// 不同类型组注册同一 compType 属于配置错误，直接 panic。
func (r *ComponentRegistry) Register(group string, compType ComponentIntType, factory ComponentSetFactory, unmarshaler ComponentSetUnmarshaler, size int, proto func() Component) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// conflict check
	if info, ok := r.infos[compType]; ok {
		if info.Group != group {
			panic(fmt.Sprintf("component type %d is conflict, group %s vs %s", compType, info.Group, group))
		}
		return
	}

	r.infos[compType] = &ComponentDefineInfo{
		Group:       group,
		Type:        compType,
		Factory:     factory,
		Unmarshaler: unmarshaler,
		Size:        size,
		Proto:       proto,
	}
}

// RegisterComponent is a convenience function to register a component type
// It automatically creates the factory and unmarshaler from the component type
func RegisterComponent[T any, TP ComponentPointer[T]](group string) {
	compType := GetIntType[T, TP]()
	globalComponentRegistry.Register(
		group,
		compType,
		func() ComponentSet {
			return NewCSet[T, TP]()
		},
		func(data *SerializableSparseArrayData) ComponentSet {
			return NewCSetFromData[T, TP](data)
		},
		int(unsafe.Sizeof(*new(T))),
		func() Component {
			return TP(new(T))
		},
	)
}

// GetInfo returns the full define info for a component type
func (r *ComponentRegistry) GetInfo(compType ComponentIntType) (*ComponentDefineInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.infos[compType]
	return info, ok
}

// GetFactory returns the factory for a component type
func (r *ComponentRegistry) GetFactory(compType ComponentIntType) (ComponentSetFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.infos[compType]
	return info.Factory, ok
}

// GetUnmarshaler returns the unmarshaler for a component type
func (r *ComponentRegistry) GetUnmarshaler(compType ComponentIntType) (ComponentSetUnmarshaler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.infos[compType]
	return info.Unmarshaler, ok
}

// CreateComponentSet creates a new ComponentSet for the given type
func (r *ComponentRegistry) CreateComponentSet(compType ComponentIntType) (ComponentSet, bool) {
	f, ok := r.GetFactory(compType)
	if !ok {
		return nil, false
	}
	return f(), true
}

// UnmarshalComponentSet creates and unmarshals a ComponentSet from serialized data
func (r *ComponentRegistry) UnmarshalComponentSet(compType ComponentIntType, data *SerializableSparseArrayData) (ComponentSet, bool) {
	u, ok := r.GetUnmarshaler(compType)
	if !ok {
		return nil, false
	}
	return u(data), true
}

// IsRegistered checks if a component type is registered
func (r *ComponentRegistry) IsRegistered(compType ComponentIntType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.infos[compType]
	return ok
}

// RegisteredTypes returns all registered component types
func (r *ComponentRegistry) RegisteredTypes() []ComponentIntType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]ComponentIntType, 0, len(r.infos))
	for t := range r.infos {
		types = append(types, t)
	}
	return types
}
