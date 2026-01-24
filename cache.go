package exoskeleton

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// Cache provides caching for expensive string-producing operations.
// Implementations control their own expiry/invalidation logic.
//
// Cache implementations must be safe for concurrent use.
type Cache interface {
	// Fetch returns a cached value for the given command and key, or calls compute()
	// to generate it. If compute() returns an error, the error is returned and the
	// value is not cached.
	//
	// The key identifies the operation (e.g., "summary", "describe-commands").
	// Implementations typically use cmd.Path() combined with key to form a cache key,
	// and may use the file's modification time for cache invalidation.
	Fetch(cmd Command, key string, compute func() (string, error)) (string, error)
}

// nullCache is the default cache that performs no caching.
// It simply invokes the compute function on every call.
type nullCache struct{}

func (nullCache) Fetch(_ Command, _ string, compute func() (string, error)) (string, error) {
	return compute()
}

// FileCache is a file-backed cache with mtime-based invalidation and optional TTL expiration.
// It persists cache entries to a JSON file and invalidates entries when the source file's
// mtime changes or the TTL expires.
//
// FileCache is safe for concurrent use. It uses singleflight to deduplicate concurrent
// calls with the same key.
type FileCache struct {
	// Path is the location of the cache file (e.g., "~/.myapp/cache.json").
	Path string

	// ExpiresAfter is the maximum age of a cache entry before it's considered stale.
	// If zero, entries never expire based on age (only mtime changes invalidate).
	ExpiresAfter time.Duration

	data   map[string]fileCacheEntry
	mu     sync.RWMutex
	sf     singleflight.Group
	loaded bool
}

type fileCacheEntry struct {
	Value    string `json:"value"`
	ModTime  int64  `json:"modTime"`
	CachedAt int64  `json:"cachedAt"`
}

func (c *FileCache) Fetch(cmd Command, key string, compute func() (string, error)) (string, error) {
	cacheKey := key + ":" + cmd.Path()

	result, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		return c.fetchOnce(cmd.Path(), cacheKey, compute)
	})

	if err != nil {
		return "", err
	}
	return result.(string), nil
}

func (c *FileCache) fetchOnce(path, cacheKey string, compute func() (string, error)) (string, error) {
	c.ensureLoaded()

	modTime := c.modTime(path)
	now := time.Now().Unix()

	// Check cache
	c.mu.RLock()
	entry, ok := c.data[cacheKey]
	c.mu.RUnlock()

	if ok && c.isValid(entry, modTime, now) {
		return entry.Value, nil
	}

	// Cache miss or stale
	value, err := compute()
	if err != nil {
		return "", err
	}

	// Update cache
	c.mu.Lock()
	if c.data == nil {
		c.data = make(map[string]fileCacheEntry)
	}
	c.data[cacheKey] = fileCacheEntry{
		Value:    value,
		ModTime:  modTime,
		CachedAt: now,
	}
	c.mu.Unlock()

	c.persist()
	return value, nil
}

func (c *FileCache) isValid(entry fileCacheEntry, currentModTime, now int64) bool {
	// Mtime changed = stale
	if entry.ModTime != currentModTime {
		return false
	}
	// TTL expired = stale (if TTL is configured)
	if c.ExpiresAfter > 0 && now-entry.CachedAt > int64(c.ExpiresAfter.Seconds()) {
		return false
	}
	return true
}

func (c *FileCache) modTime(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func (c *FileCache) ensureLoaded() {
	c.mu.RLock()
	loaded := c.loaded
	c.mu.RUnlock()
	if loaded {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.loaded {
		return
	}

	c.data = make(map[string]fileCacheEntry)
	if b, err := os.ReadFile(c.Path); err == nil {
		json.Unmarshal(b, &c.data) // Ignore errors, start fresh
	}
	c.loaded = true
}

func (c *FileCache) persist() {
	c.mu.RLock()
	b, _ := json.Marshal(c.data)
	c.mu.RUnlock()
	os.WriteFile(c.Path, b, 0644) // Ignore errors
}
