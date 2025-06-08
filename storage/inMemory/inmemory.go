package inMemory

import (
	"sync"
	"time"
)

// InMemoryStore provides a thread-safe in-memory key-value store using sync.Map.
type InMemoryStore struct {
	Store       sync.Map
	expirations map[string]time.Time
	mutex       sync.RWMutex // Mutex to protect the expirations map
}

var storage *InMemoryStore
var once sync.Once

func GetInMemoryStore() *InMemoryStore {
	if storage == nil {
		storage = &InMemoryStore{
			expirations: make(map[string]time.Time),
		}
		// Use sync.Once to ensure the cleanup goroutine is started only once
		once.Do(func() {
			go storage.cleanupExpiredKeys()
		})
	}
	return storage
}

// Set stores a value for a given key.
func (m *InMemoryStore) Set(key, value string) {
	m.Store.Store(key, value)

	// Remove any expiration for this key
	m.mutex.Lock()
	delete(m.expirations, key)
	m.mutex.Unlock()
}

// SetWithExpiration stores a value for a given key with an expiration time.
func (m *InMemoryStore) SetWithExpiration(key, value string, expiration time.Duration) {
	m.Store.Store(key, value)

	if expiration > 0 {
		m.mutex.Lock()
		m.expirations[key] = time.Now().Add(expiration)
		m.mutex.Unlock()
	}
}

// Returns the value and a boolean indicating if the key exists.
func (m *InMemoryStore) Get(key string) (string, bool) {
	// Check if key has expired
	m.mutex.RLock()
	expTime, hasExpiration := m.expirations[key]
	m.mutex.RUnlock()

	if hasExpiration && time.Now().After(expTime) {
		// Key has expired, delete it
		m.Delete(key)
		return "", false
	}

	val, ok := m.Store.Load(key)
	if ok {
		return val.(string), true
	}
	return "", false
}

// Delete removes a key from the store.
func (m *InMemoryStore) Delete(key string) {
	m.Store.Delete(key)

	m.mutex.Lock()
	delete(m.expirations, key)
	m.mutex.Unlock()
}

// Clear clears the entire store and expirations.
func (m *InMemoryStore) Clear() {
	m.Store = sync.Map{}

	m.mutex.Lock()
	m.expirations = make(map[string]time.Time)
	m.mutex.Unlock()
}

// Cleanup expired keys periodically
func (m *InMemoryStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		var keysToDelete []string

		// Find expired keys
		m.mutex.RLock()
		for key, expTime := range m.expirations {
			if now.After(expTime) {
				keysToDelete = append(keysToDelete, key)
			}
		}
		m.mutex.RUnlock()

		// Delete expired keys
		for _, key := range keysToDelete {
			m.Delete(key)
		}
	}
}

// GetExpirations iterates through all expiration entries and calls the provided function
func (m *InMemoryStore) GetExpirations(fn func(key string, expTime time.Time) bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for key, expTime := range m.expirations {
		if !fn(key, expTime) {
			break
		}
	}
}

// SetExpiration directly sets an expiration time for a key
func (m *InMemoryStore) SetExpiration(key string, expTime time.Time) {
	m.mutex.Lock()
	m.expirations[key] = expTime
	m.mutex.Unlock()
}

func (m *InMemoryStore) Exists(key string) bool {
	_, exists := m.Get(key)
	return exists
}

func (m *InMemoryStore) DeleteExpiration(key string) {
	m.mutex.Lock()
	delete(m.expirations, key)
	m.mutex.Unlock()
}

func (m *InMemoryStore) GetTTL(key string) int64 {
	m.mutex.RLock()
	expTime, exists := m.expirations[key]
	m.mutex.RUnlock()

	if !exists {
		return -1 // Key does not exist or has no expiration
	}

	ttl := expTime.Sub(time.Now())
	if ttl < 0 {
		return -1 // Key has expired
	}
	return int64(ttl.Seconds())
}
