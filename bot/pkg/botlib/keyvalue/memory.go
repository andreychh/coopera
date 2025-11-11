package keyvalue

import (
	"context"
	"fmt"
	"sync"
)

type memoryStore struct {
	mu         sync.RWMutex
	dataSource map[string]string
}

func (m *memoryStore) Read(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, exists := m.dataSource[key]
	if !exists {
		return "", fmt.Errorf("(%T) key not found", m)
	}
	return value, nil
}

func (m *memoryStore) Write(_ context.Context, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dataSource[key] = value
	return nil
}

func (m *memoryStore) Exists(_ context.Context, key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.dataSource[key]
	return exists, nil
}

func MemoryStore() Store {
	return &memoryStore{dataSource: make(map[string]string)}
}
