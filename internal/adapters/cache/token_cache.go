package cache

import (
	"context"
	"time"

	"github.com/kanta/backend-challenge/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

type tokenCache struct {
	client redis.Cmdable
}

func NewTokenCache(client redis.Cmdable) ports.CachePort {
	return &tokenCache{
		client: client,
	}
}

func (r *tokenCache) SetToken(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *tokenCache) GetToken(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *tokenCache) DeleteToken(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
