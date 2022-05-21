package cache

import (
	"time"

	"github.com/aditya37/geofence-service/repository"
	"github.com/go-redis/redis/v7"
)

type cacheManager struct {
	client *redis.Client
}

func NewCacheManager(redis *redis.Client) (repository.CacheManager, error) {
	return &cacheManager{
		client: redis,
	}, nil
}

// Get..
func (cm *cacheManager) Get(key string) (string, error) {
	result, err := cm.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", repository.ErrRedisNil
		}
		return "", err
	}
	return result, nil
}

// Set...
func (cm *cacheManager) Set(key string, ttl time.Duration, value []byte) error {
	if err := cm.client.Set(key, string(value), ttl).Err(); err != nil {
		return err
	}
	return nil
}

//Close...
func (cm *cacheManager) Close() error {
	return cm.client.Close()
}
