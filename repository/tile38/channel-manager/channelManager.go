package channelmanager

import (
	"context"

	"github.com/aditya37/geofence-service/repository"
	"github.com/xjem/t38c"
)

type chanManager struct {
	tile *t38c.Client
}

func NewChannelManager(tile38 *t38c.Client) (repository.Tile38ChannelManager, error) {
	return &chanManager{
		tile: tile38,
	}, nil
}
func (cm *chanManager) Subscribe(ctx context.Context, pattern string, f func(ge *t38c.GeofenceEvent)) error {
	return cm.tile.Channels.PSubscribe(ctx, f, pattern)
}
func (cm *chanManager) Close() error {
	return cm.tile.Close()
}
