package cache

import (
	"sync"
	"time"
)

var (
	// 全局缓存管理器实例
	globalManager *Manager
	once          sync.Once
)

// CacheInterface 缓存接口
type CacheInterface interface {
	Clear()
	Cleanup()
	GetItemCount() int
}

// Manager 缓存管理器
type Manager struct {
	caches   map[string]CacheInterface
	mu       sync.RWMutex
	interval time.Duration
}

// GetManager 获取全局缓存管理器实例
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			caches:   make(map[string]CacheInterface),
			interval: 5 * time.Minute, // 默认清理间隔
		}
		go globalManager.startCleanup()
	})
	return globalManager
}

// GetOrCreate 获取或创建缓存
func GetOrCreate[T any](name string) *Cache[T] {
	mgr := GetManager()
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	if cache, ok := mgr.caches[name]; ok {
		if typedCache, ok := cache.(*Cache[T]); ok {
			return typedCache
		}
		// 如果类型不匹配，删除旧缓存
		delete(mgr.caches, name)
	}

	cache := NewCache[T](mgr.interval)
	mgr.caches[name] = cache
	return cache
}

// SetCleanupInterval 设置清理间隔
func (m *Manager) SetCleanupInterval(interval time.Duration) {
	m.interval = interval
}

// Clear 清空所有缓存
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, cache := range m.caches {
		cache.Clear()
	}
}

// startCleanup 启动定期清理
func (m *Manager) startCleanup() {
	ticker := time.NewTicker(m.interval)
	for range ticker.C {
		m.cleanup()
	}
}

// cleanup 清理所有缓存中的过期项
func (m *Manager) cleanup() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, cache := range m.caches {
		cache.Cleanup()
	}
}

// Stats 缓存统计信息
type Stats struct {
	CacheCount  int            // 缓存实例数量
	ItemsCount  map[string]int // 每个缓存中的项目数量
	TotalItems  int            // 所有缓存中的总项目数量
	MemoryUsage int64          // 预估内存使用量（字节）
	CreatedTime time.Time      // 缓存管理器创建时间
	LastCleanup time.Time      // 最后一次清理时间
}

// GetStats 获取缓存统计信息
func (m *Manager) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := Stats{
		CacheCount:  len(m.caches),
		ItemsCount:  make(map[string]int),
		CreatedTime: time.Now(),
	}

	for name, cache := range m.caches {
		itemCount := cache.GetItemCount()
		stats.ItemsCount[name] = itemCount
		stats.TotalItems += itemCount
	}

	return stats
}
