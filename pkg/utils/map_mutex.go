package utils

import "sync"

type MapMutex[K comparable, V any] struct {
	mmap map[K]V
	lock sync.RWMutex
}

func NewMapMutex[K comparable, V any]() *MapMutex[K, V] {
	return &MapMutex[K, V]{
		mmap: make(map[K]V),
	}
}

func (m *MapMutex[K, V]) ComputeIfAbsent(k K, fun func() V) V {
	m.lock.Lock()
	defer m.lock.Unlock()
	if v, ok := m.mmap[k]; ok {
		return v
	}
	v := fun()
	m.mmap[k] = v
	return v
}

func (m *MapMutex[K, V]) Get(k K) (v V, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	ret, ok := m.mmap[k]
	return ret, ok
}

func (m *MapMutex[K, V]) PutAll(mm map[K]V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range mm {
		m.mmap[k] = v
	}
}
func (m *MapMutex[K, V]) Put(k K, v V) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.mmap[k] = v
}
