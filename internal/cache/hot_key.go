package cache

import (
	"sync"
	"time"
)

// cacheItem represents a value and its expiration time.
type cacheItem struct {
	value      string
	expiration time.Time
}

// HotKeyCache is a simple TTL-based in-memory cache using sync.Map
type HotKeyCache struct {
	data sync.Map
	ttl  time.Duration
}

// NewHotKeyCache creates a new cache with given TTL.
func NewHotKeyCache(ttl time.Duration) *HotKeyCache {
	return &HotKeyCache{
		ttl: ttl,
	}
}

// Get returns the cached value if not expired.
func (c *HotKeyCache) Get(key string) (string, bool) {
	itemRaw, ok := c.data.Load(key)
	if !ok {
		return "", false
	}

	item := itemRaw.(cacheItem)
	if time.Now().After(item.expiration) {
		c.data.Delete(key)
		return "", false
	}

	return item.value, true
}

// Set stores a value with expiration.
func (c *HotKeyCache) Set(key string, value string) {
	c.data.Store(key, cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl),
	})
}
