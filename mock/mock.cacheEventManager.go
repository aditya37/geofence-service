package mock

import (
	"errors"
	"log"
	"time"

	"github.com/aditya37/geofence-service/repository"
)

var (
	ValidEventJson = `{"message":"Success Add geofence area","location_id":1,"location_name":"Nothing","location_type":"Road","published_at":"2022-01-11"}`
)

type mockEventCache struct {
	event map[string]string
}

func NewMockCacheEventManager() (repository.CacheEventManager, error) {
	return &mockEventCache{
		event: make(map[string]string),
	}, nil
}

func (me *mockEventCache) SetEventStateResponse(db int, key string, value []byte, duration time.Duration) error {
	//	log.Println(string(value), key)

	me.event[key] = string(value)
	return nil
}
func (me *mockEventCache) GetEventState(db int, key string) (string, error) {
	resp, ok := me.event[key]
	if !ok {
		return "", errors.New("event nil")
	}
	log.Println(resp)
	return resp, nil
}
func (me *mockEventCache) Close() error {
	me.event[""] = ""
	return nil
}
