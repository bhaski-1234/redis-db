package inMemory

import (
	"sync"
)

// InMemoryStore provides a thread-safe in-memory key-value store using sync.Map.
type InMemoryStore struct {
	store sync.Map
}

var storage *InMemoryStore

func GetInMemoryStore() *InMemoryStore {
	if storage == nil {
		storage = &InMemoryStore{}
	}
	return storage
}

// Set stores a value for a given key.
func (m *InMemoryStore) Set(key, value string) {
	m.store.Store(key, value)
}

// Returns the value and a boolean indicating if the key exists.
func (m *InMemoryStore) Get(key string) (string, bool) {
	val, ok := m.store.Load(key)
	if ok {
		return val.(string), true
	}
	return "", false
}

// Delete removes a key from the store.
func (m *InMemoryStore) Delete(key string) {
	m.store.Delete(key)
}
