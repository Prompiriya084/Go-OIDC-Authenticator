package adapters_caching

import (
	ports_caching "OIDCAuthenticator/internal/core/ports/caching"
	"context"
	"time"
)

type cacheRepositoryImpl struct {
}

func NewCacheRepository() ports_caching.CacheRepository {
	return &cacheRepositoryImpl{}
}

func (r *cacheRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return nil
}
func (r *cacheRepositoryImpl) Get(
	ctx context.Context,
	key string,
	target interface{},
) error {
	return nil
}
