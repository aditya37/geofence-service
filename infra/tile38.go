package infra

import (
	"fmt"
	"sync"

	"github.com/xjem/t38c"
)

type Tile38ClientParam struct {
	Host string
	Port int
}

var (
	tile38ClientInstance *t38c.Client = nil
	tile38Singleton      sync.Once
)

func NewTile38Client(config Tile38ClientParam) error {
	var returnErr error
	tile38Singleton.Do(func() {
		host := fmt.Sprintf("%s:%d", config.Host, config.Port)
		client, err := t38c.New(host)
		if err != nil {
			returnErr = err
			return
		}
		// healthCheck..
		if client.HealthZ(); err != nil {
			returnErr = err
			return
		}
		tile38ClientInstance = client
	})
	if returnErr != nil {
		return returnErr

	}
	return nil
}

func GetTile38ClientInstance() *t38c.Client {
	return tile38ClientInstance
}
