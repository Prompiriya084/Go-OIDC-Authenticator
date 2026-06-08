package adapters_caching

import (
	ports_caching "OIDCAuthenticator/internal/core/ports/caching"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheRepositoryImpl struct {
	redisClient *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) ports_caching.CacheRepository {
	return &cacheRepositoryImpl{
		redisClient: redisClient,
	}
}

func (r *cacheRepositoryImpl) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// 👈 เปลี่ยนจาก 5*time.Minute มาใช้ expiration ที่ส่งมาจาก Usecase ตัวจริง
	return r.redisClient.Set(ctx, key, data, expiration).Err()
}

func (r *cacheRepositoryImpl) Get(
	ctx context.Context,
	key string,
	target interface{}, // 👈 ส่ง Pointer ของ struct ที่ต้องการให้ data ไปลงเข้ามา
) error {
	data, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	// 👈 Unmarshal ลงตัวแปร target ตรงๆ (เหมือนกับสไตล์ของ json.Unmarshal หรือก๊อปปี้ค่าใส่ช่องว่าง)
	return json.Unmarshal(data, target)
}

func (r *cacheRepositoryImpl) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
