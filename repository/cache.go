package repository

import (
	"errors"
	"io"
	"time"
)

var (
	ErrRedisNil = errors.New("Cache Nil")
)

type CacheManager interface {
	io.Closer
	Set(key string, ttl time.Duration, value []byte) error
	Get(key string) (string, error)
}
