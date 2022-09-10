package ecs

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type Map[K comparable, V any] struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[K]*entry[V]
	misses int
}

type readOnly[K comparable, V any] struct {
	m       map[K]*entry[V]
	amended bool
}

var expunged = unsafe.Pointer(new(interface{}))

type entry[V any] struct {
	p unsafe.Pointer
}

func newEntry[V any](i V) *entry[V] {
	return &entry[V]{p: unsafe.Pointer(&i)}
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	read, _ := m.read.Load().(readOnly[K, V])
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly[K, V])
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return *new(V), false
	}
	return e.load()
}

func (e *entry[V]) load() (value V, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expunged {
		return *new(V), false
	}
	return *(*V)(p), true
}

func (m *Map[K, V]) Store(key K, value V) {
	read, _ := m.read.Load().(readOnly[K, V])
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly[K, V])
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
		e.storeLocked(&value)
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnly[K, V]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry[V](value)
	}
	m.mu.Unlock()
}

func (e *entry[V]) tryStore(i *V) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entry[V]) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}

func (e *entry[V]) storeLocked(i *V) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	read, _ := m.read.Load().(readOnly[K, V])
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnly[K, V])
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnly[K, V]{m: read.m, amended: true})
		}
		m.dirty[key] = newEntry(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entry[V]) tryLoadOrStore(i V) (actual V, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expunged {
		return *new(V), false, false
	}
	if p != nil {
		return *(*V)(p), true, true
	}

	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expunged {
			return *new(V), false, false
		}
		if p != nil {
			return *(*V)(p), true, true
		}
	}
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	read, _ := m.read.Load().(readOnly[K, V])
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly[K, V])
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete()
	}
	return *new(V), false
}

func (m *Map[K, V]) Delete(key K) {
	m.LoadAndDelete(key)
}

func (e *entry[V]) delete() (value V, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return *new(V), false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*V)(p), true
		}
	}
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) bool {
	read, _ := m.read.Load().(readOnly[K, V])
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnly[K, V])
		if read.amended {
			read = readOnly[K, V]{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			return false
		}
	}
	return true
}

func (m *Map[K, V]) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnly[K, V]{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *Map[K, V]) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnly[K, V])
	m.dirty = make(map[K]*entry[V], len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entry[V]) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expunged) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expunged
}
