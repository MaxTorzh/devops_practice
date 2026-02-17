package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheRepository(client *redis.Client, ttl time.Duration) *CacheRepository {
	return &CacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *CacheRepository) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *CacheRepository) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *CacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *CacheRepository) Flush(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}

// Вспомогательные методы для формирования ключей
func UserKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

func UserListKey(page, limit int) string {
	return fmt.Sprintf("users:list:%d:%d", page, limit)
}

func ProductKey(id int) string {
	return fmt.Sprintf("product:%d", id)
}

func ProductListKey(page, limit int) string {
	return fmt.Sprintf("products:list:%d:%d", page, limit)
}
