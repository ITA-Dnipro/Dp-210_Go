package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RestoreCodeCache struct {
	redis      *redis.Client
	Expiration time.Duration
	typ        string
}

func NewRestoreCodeCache(redis *redis.Client, expiration time.Duration, typ string) *RestoreCodeCache {
	return &RestoreCodeCache{
		redis:      redis,
		Expiration: expiration,
		typ:        typ + ":",
	}
}

func (sc *RestoreCodeCache) Get(ctx context.Context, key string) (string, error) {
	return sc.redis.Get(ctx, sc.typ+key).Result()
}

func (sc *RestoreCodeCache) Set(ctx context.Context, key, value string) error {
	return sc.redis.SetNX(ctx, sc.typ+key, value, sc.Expiration).Err()
}

func (sc *RestoreCodeCache) Del(ctx context.Context, key string) error {
	return sc.redis.Del(ctx, sc.typ+key).Err()
}
