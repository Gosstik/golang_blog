package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const mainPageCacheKey = "cache:posts:main_page"

// CacheRepository caches the main page posts list in Redis.
type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

func (r *CacheRepository) GetPostsListCache(ctx context.Context) ([]byte, error) {
	data, err := r.client.Get(ctx, mainPageCacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return data, err
}

func (r *CacheRepository) SetPostsListCache(ctx context.Context, data []byte, ttl time.Duration) error {
	return r.client.Set(ctx, mainPageCacheKey, data, ttl).Err()
}

func (r *CacheRepository) InvalidatePostsListCache(ctx context.Context) error {
	return r.client.Del(ctx, mainPageCacheKey).Err()
}
