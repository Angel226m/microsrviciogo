// ═══════════════════════════════════════════════════════════════
// Adaptador de Caché Redis – Adaptador secundario para almacenamiento en caché
// Implementa port.CacheRepository usando Redis (go-redis/v9)
// ═══════════════════════════════════════════════════════════════
package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudmart/user-service/internal/domain/port"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

// NewRedisCache crea un nuevo adaptador de caché Redis implementando port.CacheRepository.
func NewRedisCache(client *redis.Client) port.CacheRepository {
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
