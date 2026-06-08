package ports_caching

import (
	"context"
	"time"
)

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, target interface{}) error
	Delete(ctx context.Context, key string) error
}
