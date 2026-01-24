package exoskeleton

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/square/exoskeleton/pkg/shellcomp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCommand implements Command interface for testing
type mockCommand struct {
	path string
}

func (m *mockCommand) Parent() Module                             { return nil }
func (m *mockCommand) Path() string                               { return m.path }
func (m *mockCommand) Name() string                               { return filepath.Base(m.path) }
func (m *mockCommand) Summary() (string, error)                   { return "", nil }
func (m *mockCommand) Help() (string, error)                      { return "", nil }
func (m *mockCommand) Exec(*Entrypoint, []string, []string) error { return nil }
func (m *mockCommand) Complete(*Entrypoint, []string, []string) ([]string, shellcomp.Directive, error) {
	return nil, 0, nil
}

func TestNullCache(t *testing.T) {
	cache := nullCache{}
	cmd := &mockCommand{path: "/test/path"}

	result, err := cache.Fetch(cmd, "summary", compute("computed"))

	assert.NoError(t, err)
	assert.Equal(t, "computed", result)

	// Second call should also compute
	result, err = cache.Fetch(cmd, "summary", func() (string, error) {
		return "computed again", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "computed again", result)
}

func TestFileCacheReadFromCache(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	info, _ := os.Stat(cmdFile.Name())
	modTime := info.ModTime().Unix()

	// Create cache file with pre-populated entry
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())

	cacheKey := "summary:" + cmdFile.Name()
	cacheData := map[string]fileCacheEntry{
		cacheKey: {
			Value:    "CACHED VALUE",
			ModTime:  modTime,
			CachedAt: time.Now().Unix(),
		},
	}
	b, _ := json.Marshal(cacheData)
	cacheFile.Write(b)
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	result, err := cache.Fetch(cmd, "summary", compute("SHOULD NOT BE CALLED"))

	assert.NoError(t, err)
	assert.Equal(t, "CACHED VALUE", result)
}

func TestFileCacheWriteOnMiss(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	// Create empty cache file
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	result, err := cache.Fetch(cmd, "summary", compute("COMPUTED VALUE"))

	assert.NoError(t, err)
	assert.Equal(t, "COMPUTED VALUE", result)

	// Verify cache was written
	b, _ := os.ReadFile(cacheFile.Name())
	var data map[string]fileCacheEntry
	json.Unmarshal(b, &data)

	cacheKey := "summary:" + cmdFile.Name()
	assert.Contains(t, data, cacheKey)
	assert.Equal(t, "COMPUTED VALUE", data[cacheKey].Value)
}

func TestFileCacheInvalidatesOnMtimeChange(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	// Create cache file with OLD mtime
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())

	cacheKey := "summary:" + cmdFile.Name()
	oldModTime := time.Now().Add(-1 * time.Hour).Unix()
	cacheData := map[string]fileCacheEntry{
		cacheKey: {
			Value:    "STALE VALUE",
			ModTime:  oldModTime, // Different from actual file mtime
			CachedAt: time.Now().Unix(),
		},
	}
	b, _ := json.Marshal(cacheData)
	cacheFile.Write(b)
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	result, err := cache.Fetch(cmd, "summary", compute("FRESH VALUE"))

	assert.NoError(t, err)
	assert.Equal(t, "FRESH VALUE", result)
}

func TestFileCacheTTLExpiration(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	info, _ := os.Stat(cmdFile.Name())
	modTime := info.ModTime().Unix()

	// Create cache file with entry cached "long ago"
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())

	cacheKey := "summary:" + cmdFile.Name()
	cacheData := map[string]fileCacheEntry{
		cacheKey: {
			Value:    "EXPIRED VALUE",
			ModTime:  modTime,
			CachedAt: time.Now().Add(-2 * time.Hour).Unix(), // 2 hours ago
		},
	}
	b, _ := json.Marshal(cacheData)
	cacheFile.Write(b)
	cacheFile.Close()

	cache := &FileCache{
		Path:         cacheFile.Name(),
		ExpiresAfter: 1 * time.Hour, // 1 hour TTL
	}
	cmd := &mockCommand{path: cmdFile.Name()}

	result, err := cache.Fetch(cmd, "summary", compute("FRESH VALUE"))

	assert.NoError(t, err)
	assert.Equal(t, "FRESH VALUE", result)
}

func TestFileCacheNoTTLExpiration(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	info, _ := os.Stat(cmdFile.Name())
	modTime := info.ModTime().Unix()

	// Create cache file with entry cached "long ago"
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())

	cacheKey := "summary:" + cmdFile.Name()
	cacheData := map[string]fileCacheEntry{
		cacheKey: {
			Value:    "OLD BUT VALID",
			ModTime:  modTime,
			CachedAt: time.Now().Add(-24 * time.Hour).Unix(), // 24 hours ago
		},
	}
	b, _ := json.Marshal(cacheData)
	cacheFile.Write(b)
	cacheFile.Close()

	cache := &FileCache{
		Path:         cacheFile.Name(),
		ExpiresAfter: 0, // No TTL
	}
	cmd := &mockCommand{path: cmdFile.Name()}

	result, err := cache.Fetch(cmd, "summary", compute("SHOULD NOT BE CALLED"))

	assert.NoError(t, err)
	assert.Equal(t, "OLD BUT VALID", result)
}

func TestFileCacheSingleflight(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	// Create empty cache file
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	var callCount int64
	var wg sync.WaitGroup

	// Launch 10 concurrent requests
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cache.Fetch(cmd, "summary", func() (string, error) {
				atomic.AddInt64(&callCount, 1)
				time.Sleep(50 * time.Millisecond) // Simulate slow computation
				return "COMPUTED", nil
			})
		}()
	}

	wg.Wait()
	assert.Equal(t, int64(1), callCount, "singleflight should deduplicate concurrent calls")
}

func TestFileCacheComputeError(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	// Create empty cache file
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	expectedErr := os.ErrNotExist
	result, err := cache.Fetch(cmd, "summary", func() (string, error) {
		return "", expectedErr
	})

	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, "", result)

	// Verify nothing was cached
	b, _ := os.ReadFile(cacheFile.Name())
	var data map[string]fileCacheEntry
	json.Unmarshal(b, &data)
	assert.Empty(t, data)
}

func TestFileCacheDifferentKeys(t *testing.T) {
	// Create a temporary file to act as the "command"
	cmdFile, err := os.CreateTemp("", "cache-test-cmd")
	require.NoError(t, err)
	defer os.Remove(cmdFile.Name())
	cmdFile.WriteString("content")
	cmdFile.Close()

	// Create empty cache file
	cacheFile, err := os.CreateTemp("", "cache-test.json")
	require.NoError(t, err)
	defer os.Remove(cacheFile.Name())
	cacheFile.Close()

	cache := &FileCache{Path: cacheFile.Name()}
	cmd := &mockCommand{path: cmdFile.Name()}

	result1, _ := cache.Fetch(cmd, "summary", compute("SUMMARY VALUE"))
	result2, _ := cache.Fetch(cmd, "describe-commands", compute("DESCRIBE VALUE"))

	assert.Equal(t, "SUMMARY VALUE", result1)
	assert.Equal(t, "DESCRIBE VALUE", result2)

	// Verify both are cached separately
	b, _ := os.ReadFile(cacheFile.Name())
	var data map[string]fileCacheEntry
	json.Unmarshal(b, &data)

	assert.Contains(t, data, "summary:"+cmdFile.Name())
	assert.Contains(t, data, "describe-commands:"+cmdFile.Name())
}

func compute(value string) func() (string, error) {
	return func() (string, error) {
		return value, nil
	}
}
