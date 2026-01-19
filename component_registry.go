package ecs

import "sync"

// ComponentSetFactory is a function that creates a new ComponentSet
// and optionally unmarshals data into it
type ComponentSetFactory func() ComponentSet

// ComponentSetUnmarshaler is a function that unmarshals data into a ComponentSet
type ComponentSetUnmarshaler func(data *SerializableSparseArrayData) ComponentSet

// ComponentRegistry manages component type registrations for serialization/deserialization
type ComponentRegistry struct {
	mu           sync.RWMutex
	factories    map[ComponentIntType]ComponentSetFactory
	unmarshalers map[ComponentIntType]ComponentSetUnmarshaler
}

// Global component registry instance
var globalComponentRegistry = &ComponentRegistry{
	factories:    make(map[ComponentIntType]ComponentSetFactory),
	unmarshalers: make(map[ComponentIntType]ComponentSetUnmarshaler),
}

// GetComponentRegistry returns the global component registry
func GetComponentRegistry() *ComponentRegistry {
	return globalComponentRegistry
}

// Register registers a component type with its factory and unmarshaler
func (r *ComponentRegistry) Register(compType ComponentIntType, factory ComponentSetFactory, unmarshaler ComponentSetUnmarshaler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[compType] = factory
	r.unmarshalers[compType] = unmarshaler
}

// RegisterComponent is a convenience function to register a component type
// It automatically creates the factory and unmarshaler from the component type
func RegisterComponent[T ComponentObject, TP ComponentPointer[T]](comp TP) {
	compType := GetIntType[T, TP]()
	globalComponentRegistry.Register(
		compType,
		func() ComponentSet {
			return NewCSet[T]()
		},
		func(data *SerializableSparseArrayData) ComponentSet {
			return NewCSetFromData[T](data)
		},
	)
}

// GetFactory returns the factory for a component type
func (r *ComponentRegistry) GetFactory(compType ComponentIntType) (ComponentSetFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[compType]
	return f, ok
}

// GetUnmarshaler returns the unmarshaler for a component type
func (r *ComponentRegistry) GetUnmarshaler(compType ComponentIntType) (ComponentSetUnmarshaler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.unmarshalers[compType]
	return u, ok
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
	_, ok := r.factories[compType]
	return ok
}

// RegisteredTypes returns all registered component types
func (r *ComponentRegistry) RegisteredTypes() []ComponentIntType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]ComponentIntType, 0, len(r.factories))
	for t := range r.factories {
		types = append(types, t)
	}
	return types
}
