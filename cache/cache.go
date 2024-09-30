package cache

import (
	"sync"
	"time"
)

// CacheItem represents an item in the cache
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// Cache is the in-memory cache structure
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCache creates a new Cache instance
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set adds a new item to the cache with an expiration time
func (c *Cache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache, returns nil if not found or expired
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists || item.ExpiresAt.Before(time.Now()) {
		// If the item does not exist or has expired, return nil
		return nil, false
	}

	return item.Data, true
}

// CleanUp removes expired items from the cache
func (c *Cache) CleanUp() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.ExpiresAt.Before(time.Now()) {
			delete(c.items, key)
		}
	}
}
