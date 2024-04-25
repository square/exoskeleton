package exoskeleton

import (
	"encoding/json"
	"os"
)

type summaryCache struct {
	Path    string
	data    *cache
	onError func(error)
}

type cache struct {
	Summary map[string]cachedValue `json:"summary"`
}

type cachedValue struct {
	ModTime int64  `json:"modTime"`
	Value   string `json:"value"`
}

func (c *summaryCache) load() {
	c.data = &cache{}

	// cache is not configured
	if c.Path == "" {
		return
	}

	if b, err := os.ReadFile(c.Path); os.IsNotExist(err) {
		// cache file just does not exist yet
		return
	} else if err != nil {
		c.onError(CacheError{Cause: err, Message: "could not load cache"})
	} else if err := json.Unmarshal(b, &c.data); err != nil {
		c.onError(CacheError{Cause: err, Message: "could not load cache"})
	}
}

func (c *summaryCache) dump() {
	if c.Path == "" {
		return
	}

	if b, err := json.Marshal(c.data); err != nil {
		c.onError(CacheError{Cause: err, Message: "could not write cache"})
	} else if err = os.WriteFile(c.Path, b, 0644); err != nil {
		c.onError(CacheError{Cause: err, Message: "could not write cache"})
	}
}

func (c *summaryCache) Read(cmd Command) (string, error) {
	if _, ok := cmd.(*builtinCommand); ok {
		return cmd.Summary()
	}

	if c.data == nil {
		c.load()
	}

	cacheKey := Usage(cmd)

	modTime, err := modTime(cmd)
	if err != nil {
		c.onError(CacheError{Cause: err, Message: "skipping cache for " + cmd.Path()})
	} else if item, ok := c.data.Summary[cacheKey]; ok && item.ModTime == modTime {
		return item.Value, nil
	}

	summary, err := cmd.Summary()
	if err == nil {
		if c.data.Summary == nil {
			c.data.Summary = make(map[string]cachedValue)
		}
		c.data.Summary[cacheKey] = cachedValue{ModTime: modTime, Value: summary}
		c.dump()
	}
	return summary, err
}

func modTime(cmd Command) (int64, error) {
	if info, err := os.Stat(cmd.Path()); err != nil {
		return -1, err
	} else {
		return info.ModTime().Unix(), nil
	}
}
