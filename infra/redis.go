package infra

import (
	"errors"
	"fmt"
	"log"
	"sync"

	cache "github.com/go-redis/redis/v7"
)

type RedisConfigParam struct {
	Port     int
	Host     string
	Database int
	Password string
}

var (
	redisInstance          *cache.Client = nil
	redisInstanceSingelton sync.Once
)

func NewRedisInstance(param RedisConfigParam) {
	redisInstanceSingelton.Do(func() {
		client := cache.NewClient(&cache.Options{
			Addr:     fmt.Sprintf("%s:%d", param.Host, param.Port),
			DB:       param.Database,
			Password: param.Password,
		})
		ping, err := client.Ping().Result()
		if err != nil {
			log.Fatal(err)
		}
		if ping != "PONG" {
			log.Fatal(errors.New("Failed to connect redis"))
		}
		redisInstance = client
	})
}

func GetRedisInstance() *cache.Client {
	return redisInstance
}
