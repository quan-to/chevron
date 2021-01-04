package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/quan-to/slog"
)

type rediser interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.StatusCmd
	SetXX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd
	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) *redis.BoolCmd

	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

// Driver is a database handler proxy for caching
type Driver struct {
	proxy ProxiedHandler
	log   slog.Instance
	redis rediser
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
func (h *Driver) Setup(client rediser, maxLocalObjects int, localObjectTTL time.Duration) error {
	if maxLocalObjects == 0 {
		return fmt.Errorf("max local objects can't be zero")
	}
	if localObjectTTL == 0 {
		return fmt.Errorf("local object TTL can't be zero")
	}
	if client == nil {
		return fmt.Errorf("you should specify a redis client")
	}
	h.redis = client
	h.cache = cache.New(&cache.Options{
		Redis:      h.redis,
		LocalCache: cache.NewTinyLFU(maxLocalObjects, localObjectTTL),
	})

	return nil
}
