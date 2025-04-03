package memory

import (
	"context"
	"sync"
	"time"

	"github.com/iamsorryprincess/go-project-layout/internal/pkg/database"
)

type cacheItem[T any] struct {
	ttl       time.Duration
	createdAt time.Time
	value     T
}

type Cache[TKey comparable, TValue any] struct {
	mu sync.RWMutex
	wg sync.WaitGroup

	clearInterval time.Duration

	storage map[TKey]cacheItem[TValue]
}

func NewCache[TKey comparable, TValue any](ctx context.Context, config Config) *Cache[TKey, TValue] {
	cache := &Cache[TKey, TValue]{
		clearInterval: config.ClearInterval,
		storage:       make(map[TKey]cacheItem[TValue]),
	}

	cache.wg.Add(1)
	go func(ctx context.Context, cache *Cache[TKey, TValue]) {
		defer cache.wg.Done()
		timer := time.NewTimer(cache.clearInterval)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				if err := cache.clearExpiredItems(ctx); err != nil {
					return
				}
				timer.Reset(cache.clearInterval)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, cache)

	return cache
}

func (c *Cache[TKey, TValue]) Set(_ context.Context, key TKey, ttl time.Duration, value TValue) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage[key] = cacheItem[TValue]{
		ttl:       ttl,
		createdAt: time.Now(),
		value:     value,
	}
	return nil
}

func (c *Cache[TKey, TValue]) Get(_ context.Context, key TKey) (TValue, error) {
	c.mu.RLock()
	item, ok := c.storage[key]
	if !ok {
		c.mu.RUnlock()
		return item.value, database.ErrNotFound
	}

	if time.Since(item.createdAt) >= item.ttl {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.storage, key)
		c.mu.Unlock()
		return item.value, database.ErrNotFound
	}

	c.mu.RUnlock()
	return item.value, nil
}

func (c *Cache[TKey, TValue]) Wait() {
	c.wg.Wait()
}

func (c *Cache[TKey, TValue]) clearExpiredItems(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.storage {
		if err := ctx.Err(); err != nil {
			return err
		}
		if time.Since(v.createdAt) >= v.ttl {
			delete(c.storage, k)
		}
	}
	return nil
}
