package repository

import (
	"context"
	"fmt"
	"io"

	tileEntity "github.com/aditya37/geofence-service/entity/tile38"
	"github.com/xjem/t38c"
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

// Detect...
type Detect string

const (
	Inside  Detect = "inside"
	Outside Detect = "outside"
	Enter   Detect = "enter"
	Exit    Detect = "exit"
	Cross   Detect = "cross"
)

func (d Detect) ToString() string {
	return fmt.Sprintf("%s", d)
}

type Tile38ChannelManager interface {
	io.Closer
	SetGeofenceChannel(param tileEntity.Geofence) error
	Subscribe(ctx context.Context, pattern string, f func(ge *t38c.GeofenceEvent)) error
	GetChannelDetail(pattern string) ([]t38c.Chan, error)
}
