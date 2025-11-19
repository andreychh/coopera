package immutable

import (
	"iter"
	"maps"
)

type immutableMap[K comparable, V any] struct {
	items map[K]V
}

func (m immutableMap[K, V]) Get(key K) (V, bool) {
	value, exists := m.items[key]
	return value, exists
}

func (m immutableMap[K, V]) With(key K, value V) Map[K, V] {
	copied := make(map[K]V, len(m.items)+1)
	maps.Copy(copied, m.items)
	copied[key] = value
	return immutableMap[K, V]{items: copied}
}

func (m immutableMap[K, V]) Remove(key K) Map[K, V] {
	if _, exists := m.items[key]; !exists {
		return m
	}
	copied := make(map[K]V, len(m.items)-1)
	maps.Copy(copied, m.items)
	delete(copied, key)
	return immutableMap[K, V]{items: copied}
}

func (m immutableMap[K, V]) All() iter.Seq2[K, V] {
	return maps.All(m.items)
}

func (m immutableMap[K, V]) Len() int {
	return len(m.items)
}

func EmptyMap[K comparable, V any]() Map[K, V] {
	return immutableMap[K, V]{items: make(map[K]V)}
}

func MapOf[K comparable, V any](origin map[K]V) Map[K, V] {
	copied := make(map[K]V, len(origin))
	maps.Copy(copied, origin)
	return immutableMap[K, V]{items: copied}
}
