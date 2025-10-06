package docs

import (
	"sync"
	"time"
)

type cacheEntry struct {
	content   string
	timestamp time.Time
}

// Cache manages document caching with TTL
type Cache struct {
	fileCache   map[string]cacheEntry
	searchCache map[string]cacheEntry
	mu          sync.RWMutex
	ttl         time.Duration
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		fileCache:   make(map[string]cacheEntry),
		searchCache: make(map[string]cacheEntry),
		ttl:         5 * time.Minute,
	}
}

// Get retrieves content from cache
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.fileCache[key]
	if !exists {
		return "", false
	}

	// Check if expired
	if time.Since(entry.timestamp) > c.ttl {
		return "", false
	}

	return entry.content, true
}

// Set stores content in cache
func (c *Cache) Set(key, content string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.fileCache[key] = cacheEntry{
		content:   content,
		timestamp: time.Now(),
	}

	// Limit cache size
	if len(c.fileCache) > 200 {
		c.evictOldest()
	}
}

// GetSearch retrieves search results from cache
func (c *Cache) GetSearch(query string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.searchCache[query]
	if !exists {
		return "", false
	}

	if time.Since(entry.timestamp) > c.ttl {
		return "", false
	}

	return entry.content, true
}

// SetSearch stores search results in cache
func (c *Cache) SetSearch(query, results string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.searchCache[query] = cacheEntry{
		content:   results,
		timestamp: time.Now(),
	}

	// Limit cache size
	if len(c.searchCache) > 100 {
		c.evictOldestSearch()
	}
}

// Clear removes all cached entries
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.fileCache = make(map[string]cacheEntry)
	c.searchCache = make(map[string]cacheEntry)
}

// evictOldest removes oldest 20% of file cache entries
func (c *Cache) evictOldest() {
	toRemove := len(c.fileCache) / 5
	count := 0

	for key := range c.fileCache {
		delete(c.fileCache, key)
		count++
		if count >= toRemove {
			break
		}
	}
}

// evictOldestSearch removes oldest 20% of search cache entries
func (c *Cache) evictOldestSearch() {
	toRemove := len(c.searchCache) / 5
	count := 0

	for key := range c.searchCache {
		delete(c.searchCache, key)
		count++
		if count >= toRemove {
			break
		}
	}
}
