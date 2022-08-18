package eventmanager

import (
	"errors"
	"time"

	"github.com/aditya37/geofence-service/repository"
	"github.com/go-redis/redis/v7"
)

type eventmanager struct {
	client *redis.Client
}

// Deprecated: Not used...
func NewCacheEventManager(client *redis.Client) (repository.CacheEventManager, error) {
	return &eventmanager{
		client: client,
	}, nil
}

// [GEOP-42] SetEventStateResponse....
func (em *eventmanager) SetEventStateResponse(db int, key string, value []byte, duration time.Duration) error {

	// switch to cache event database
	cmd := redis.NewStringCmd("SELECT", db)
	if err := em.client.Process(cmd); err != nil {
		return err
	}

	// set..
	if err := em.client.Set(
		key,
		string(value),
		duration,
	).Err(); err != nil {
		return err
	}
	return nil
}

// GetEventState..
func (em *eventmanager) GetEventState(db int, key string) (string, error) {
	// move db
	cmd := redis.NewStringCmd("SELECT", db)
	if err := em.client.Process(cmd); err != nil {
		return "", err
	}
	result, err := em.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("event log not found")
		}
		return "", err
	}
	return result, nil
}
func (em *eventmanager) Close() error {
	return em.client.Close()
}
