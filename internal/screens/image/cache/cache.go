package cache

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

var ErrCacheExpired = errors.New("cache expired")

type Cache struct {
	dir     string
	ttl     time.Duration
	keyFunc func(url string) string
}

func NewCache(dir string, ttl time.Duration, keyFunc func(url string) string) *Cache {
	return &Cache{
		dir:     dir,
		ttl:     ttl,
		keyFunc: keyFunc,
	}
}

func (c *Cache) Get(key string) ([]byte, error) {
	key = c.keyFunc(key)

	cacheFile := path.Join(c.dir, key)
	fileExists, err := os.Stat(cacheFile)
	if err != nil {
		return nil, ErrCacheExpired
	}

	if time.Since(fileExists.ModTime()) > c.ttl {
		return nil, ErrCacheExpired
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	return data, nil
}

func (c *Cache) Set(key string, data []byte) error {
	key = c.keyFunc(key)

	cacheFile := path.Join(c.dir, key)
	err := os.WriteFile(cacheFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}
