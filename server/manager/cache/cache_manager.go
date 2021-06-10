package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type Manager struct {
	redis *redis.Client
}

func NewCacheManager() *Manager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return &Manager{redis: rdb}
}

func (m *Manager) Ping(ctx context.Context) (string, error) {
	return m.redis.Ping(ctx).Result()
}

func (m *Manager) Get(ctx context.Context, key string) *redis.StringCmd {
	return m.redis.Get(ctx, key)
}

func (m *Manager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return m.redis.Set(ctx, key, value, expiration)
}

func (m *Manager) Del(ctx context.Context, key string) *redis.IntCmd {
	return m.redis.Del(ctx, key)
}
