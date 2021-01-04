package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/quan-to/slog"
)

// Driver is a database handler proxy for caching
type Driver struct {
	proxy ProxiedHandler
	log   slog.Instance
	redis *redis.ClusterClient
	cache *cache.Cache
}

// MakeRedisDriver creates a Redis Caching layer for the specified handler
// It also implements the Token mechanics (it does not pass to the underlying handler)
func MakeRedisDriver(dbh ProxiedHandler, log slog.Instance) *Driver {
	if log == nil {
		log = slog.Scope("REDIS")
	} else {
		log = log.SubScope("REDIS")
	}
	return &Driver{proxy: dbh, log: log}
}

// HealthCheck returns nil if everything is OK with the handler
func (h *Driver) HealthCheck() error {
	// This might deviate the statistics,
	// but its the only way I found out to test the connection
	err := h.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   userTokenPrefix + "__HC__",
		Value: &struct{}{},
		TTL:   time.Second * 4,
	})
	if err == nil {
		// Test the proxy
		return h.proxy.HealthCheck()
	}
	return err
}

// Setup configures the RedisDriver connection and cache ring
func (h *Driver) Setup(opts *redis.ClusterOptions, maxLocalObjects int, localObjectTTL time.Duration) error {
	if maxLocalObjects == 0 {
		return fmt.Errorf("max local objects can't be zero")
	}
	if localObjectTTL == 0 {
		return fmt.Errorf("local object TTL can't be zero")
	}

	h.redis = redis.NewClusterClient(opts)
	h.cache = cache.New(&cache.Options{
		Redis:      h.redis,
		LocalCache: cache.NewTinyLFU(maxLocalObjects, localObjectTTL),
	})

	return nil
}
