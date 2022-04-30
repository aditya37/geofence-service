package repository

import (
	"io"
	"time"
)

// manage event cache
type CacheEventManager interface {
	io.Closer
	SetEventStateResponse(db int, key string, value []byte, duration time.Duration) error
	GetEventState(db int, key string) (string, error)
}
