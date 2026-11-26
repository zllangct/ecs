package ecs

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

// ComponentSetFactory is a function that creates a new ComponentSet
// and optionally unmarshals data into it
type ComponentSetFactory func() ComponentSet

// ComponentSetUnmarshaler is a function that unmarshals data into a ComponentSet.
// 返回 error 表示快照数据非法/不兼容（R19：恢复失败必须可传播，不得 panic 或静默跳过）
type ComponentSetUnmarshaler func(data *SerializableSparseArrayData) (ComponentSet, error)

type ComponentDefineInfo struct {
	Group       string
	Type        ComponentIntType
	Factory     ComponentSetFactory
	Unmarshaler ComponentSetUnmarshaler
	// Size 组件定长字节数（unsafe.Sizeof(T)），组表列 stride 用
	Size int
	// HasPointers 组件类型是否包含 Go 指针（含 string/slice/map/chan/func/interface/*T）
	HasPointers bool
	// Proto 组件原型构造器，注册期读取 IsDisposable/IsNomadic 等标记
	Proto  func() Component
	goType reflect.Type
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
// 不同类型组或不同 Go 类型注册同一 compType 属于配置/哈希碰撞错误，直接 panic。
func (r *ComponentRegistry) Register(group string, compType ComponentIntType, factory ComponentSetFactory, unmarshaler ComponentSetUnmarshaler, size int, proto func() Component) {
	r.registerWithMeta(group, compType, factory, unmarshaler, size, true, nil, proto)
}

func (r *ComponentRegistry) registerWithMeta(group string, compType ComponentIntType, factory ComponentSetFactory, unmarshaler ComponentSetUnmarshaler, size int, hasPointers bool, goType reflect.Type, proto func() Component) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// conflict check
	if info, ok := r.infos[compType]; ok {
		if info.Group != group {
			panic(fmt.Sprintf("component type %d is conflict, group %s vs %s", compType, info.Group, group))
		}
		if info.goType != nil && goType != nil && info.goType != goType {
			panic(fmt.Sprintf("component type %d is conflict in group %s: %v vs %v", compType, group, info.goType, goType))
		}
		return
	}

	r.infos[compType] = &ComponentDefineInfo{
		Group:       group,
		Type:        compType,
		Factory:     factory,
		Unmarshaler: unmarshaler,
		Size:        size,
		HasPointers: hasPointers,
		Proto:       proto,
		goType:      goType,
	}
}

// typeHasPointers 递归检查 Go 类型是否包含任何 GC 追踪的指针引用
func typeHasPointers(t reflect.Type) bool {
	if t == nil {
		return true
	}
	switch t.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.String, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return true
	case reflect.Array:
		if t.Len() == 0 {
			return false
		}
		return typeHasPointers(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if typeHasPointers(t.Field(i).Type) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// RegisterComponent is a convenience function to register a component type
// It automatically creates the factory and unmarshaler from the component type
func RegisterComponent[T any, TP ComponentPointer[T]](group string) {
	compType := GetIntType[T, TP]()
	typ := TypeOf[T]()
	globalComponentRegistry.registerWithMeta(
		group,
		compType,
		func() ComponentSet {
			return NewCSet[T, TP]()
		},
		func(data *SerializableSparseArrayData) (ComponentSet, error) {
			return NewCSetFromData[T, TP](data)
		},
		int(unsafe.Sizeof(*new(T))),
		typeHasPointers(typ),
		typ,
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
	if !ok {
		return nil, false
	}
	return info.Factory, true
}

// GetUnmarshaler returns the unmarshaler for a component type
func (r *ComponentRegistry) GetUnmarshaler(compType ComponentIntType) (ComponentSetUnmarshaler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.infos[compType]
	if !ok {
		return nil, false
	}
	return info.Unmarshaler, true
}

// CreateComponentSet creates a new ComponentSet for the given type
func (r *ComponentRegistry) CreateComponentSet(compType ComponentIntType) (ComponentSet, bool) {
	f, ok := r.GetFactory(compType)
	if !ok {
		return nil, false
	}
	return f(), true
}

// UnmarshalComponentSet creates and unmarshals a ComponentSet from serialized data.
// 未注册类型返回 error（快照兼容性失败必须显式传播，不得静默丢弃业务状态）
func (r *ComponentRegistry) UnmarshalComponentSet(compType ComponentIntType, data *SerializableSparseArrayData) (ComponentSet, error) {
	u, ok := r.GetUnmarshaler(compType)
	if !ok {
		return nil, fmt.Errorf("ecs: component type %d not registered", compType)
	}
	return u(data)
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
