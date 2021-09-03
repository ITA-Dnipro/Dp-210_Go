package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type SessionCache struct {
	redis      *redis.Client
	Expiration time.Duration
	typ        string
}

func NewSessionCache(redis *redis.Client, expiration time.Duration, typ string) *SessionCache {
	return &SessionCache{
		redis:      redis,
		Expiration: expiration,
		typ:        typ + ":",
	}
}

func (sc *SessionCache) Get(ctx context.Context, key string) (string, error) {
	return sc.redis.Get(ctx, sc.typ+key).Result()
}

func (sc *SessionCache) Set(ctx context.Context, key, value string) error {
	return sc.redis.SetNX(ctx, sc.typ+key, value, sc.Expiration).Err()
}

func (sc *SessionCache) Del(ctx context.Context, key string) error {
	return sc.redis.Del(ctx, sc.typ+key).Err()
}
