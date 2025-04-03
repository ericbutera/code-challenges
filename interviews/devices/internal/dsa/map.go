package dsa

import "sync"

// TODO: test coverage

// A map enhanced with generics built upon sync.Map
type Map[K comparable, V any] struct {
	m sync.Map
}

func (gm *Map[K, V]) Store(key K, value V) {
	gm.m.Store(key, value)
}

func (gm *Map[K, V]) Load(key K) (V, bool) { //nolint:ireturn
	v, ok := gm.m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true //nolint:forcetypeassert
}

func (gm *Map[K, V]) Delete(key K) {
	gm.m.Delete(key)
}
