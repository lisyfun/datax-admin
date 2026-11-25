// Package loginlimit 提供进程内的登录失败限流。
//
// 这不是分布式限流，适用于单实例部署；多实例部署时应替换为 Redis 实现。
package loginlimit

import (
	"sync"
	"time"
)

type entry struct {
	failures int
	resetAt  time.Time
}

// Limiter 基于固定窗口的失败计数器。
type Limiter struct {
	mu        sync.Mutex
	entries   map[string]*entry
	max       int
	window    time.Duration
	lastSweep time.Time
}

// New 创建一个限流器：window 内失败 max 次后拒绝。
func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		entries:   make(map[string]*entry),
		max:       max,
		window:    window,
		lastSweep: time.Now(),
	}
}

// Allowed 返回该 key 当前是否允许尝试登录。
func (l *Limiter) Allowed(key string) bool {
	if l == nil {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sweepLocked()

	e, ok := l.entries[key]
	if !ok {
		return true
	}
	if time.Now().After(e.resetAt) {
		delete(l.entries, key)
		return true
	}
	return e.failures < l.max
}

// Fail 记录一次失败。
func (l *Limiter) Fail(key string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	e, ok := l.entries[key]
	if !ok || now.After(e.resetAt) {
		l.entries[key] = &entry{failures: 1, resetAt: now.Add(l.window)}
		return
	}
	e.failures++
}

// Reset 登录成功后清除计数。
func (l *Limiter) Reset(key string) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entries, key)
}

// Remaining 返回距离解禁的剩余时间。
func (l *Limiter) Remaining(key string) time.Duration {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[key]
	if !ok {
		return 0
	}
	if d := time.Until(e.resetAt); d > 0 {
		return d
	}
	return 0
}

func (l *Limiter) sweepLocked() {
	now := time.Now()
	if now.Sub(l.lastSweep) < l.window {
		return
	}
	for k, e := range l.entries {
		if now.After(e.resetAt) {
			delete(l.entries, k)
		}
	}
	l.lastSweep = now
}
