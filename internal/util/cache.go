package util

import (
	"encoding/hex"
	"sync"
)

// cache is a thread-safe map of keys to passwords.
//
// The idea was to have something that is strongly typed compared to [sync.Map]
type Cache struct {
	mu    sync.RWMutex
	store map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string][]byte),
	}
}

// Loads a key from the cache, returns (nil, false) if not found
func (c *Cache) Load(pass []byte) (key []byte, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	passKey := hex.EncodeToString(pass)
	key, ok = c.store[passKey]
	return key, ok
}

// Stores a key in the cache with the hex of pass as key
func (c *Cache) Store(pass, key []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	passKey := hex.EncodeToString(pass)
	c.store[passKey] = key
}
