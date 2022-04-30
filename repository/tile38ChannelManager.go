package repository

import (
	"fmt"
	"io"

	tileEntity "github.com/aditya37/geofence-service/entity/tile38"
)

// ActionType...
type ActionType string

const (
	Within ActionType = "WITHIN"
	Nearby ActionType = "NEARBY"
)

func (at ActionType) ToString() string {
	return fmt.Sprintf("%s", at)
}

type Tile38ChannelManager interface {
	io.Closer
	SetGeofenceChannel(param tileEntity.Geofence) error
}
