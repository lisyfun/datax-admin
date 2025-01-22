package cache

import (
	"sync"
	"time"
)

// Item 缓存项
type Item[T any] struct {
	Value      T
	ExpireTime time.Time
}

// Cache 缓存实现
type Cache[T any] struct {
	items    map[string]Item[T]
	mu       sync.RWMutex
	interval time.Duration
}

// NewCache 创建一个新的缓存
func NewCache[T any](cleanupInterval time.Duration) *Cache[T] {
	cache := &Cache[T]{
		items:    make(map[string]Item[T]),
		interval: cleanupInterval,
	}
	go cache.startCleanup()
	return cache
}

// Set 设置缓存项
func (c *Cache[T]) Set(key string, value T, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item[T]{
		Value:      value,
		ExpireTime: time.Now().Add(duration),
	}
}

// Get 获取缓存项
func (c *Cache[T]) Get(key string) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		var zero T
		return zero, false
	}

	if time.Now().After(item.ExpireTime) {
		var zero T
		return zero, false
	}

	return item.Value, true
}

// Delete 删除缓存项
func (c *Cache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Clear 清空所有缓存项
func (c *Cache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]Item[T])
}

// Cleanup 清理过期的缓存项
func (c *Cache[T]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.ExpireTime) {
			delete(c.items, key)
		}
	}
}

// startCleanup 启动定期清理
func (c *Cache[T]) startCleanup() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.Cleanup()
	}
}

// GetItemCount 获取缓存项数量
func (c *Cache[T]) GetItemCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}
