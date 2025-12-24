package ports

import (
	"context"
	"time"
)

type CachePort interface {
	SetToken(ctx context.Context, key string, value string, expiration time.Duration) error
	GetToken(ctx context.Context, key string) (string, error)
	DeleteToken(ctx context.Context, key string) error
}
