package repository

import (
	"io"

	tile38_entity "github.com/aditya37/geofence-service/entity/tile38"
	"github.com/xjem/t38c"
)

type Tile38ReaderWriter interface {
	io.Closer
	SetGeofencingKey(tile38_entity.SetKey) error
	GetLastGeofencingDetect(key, objId string, withfield bool) (*t38c.GetResponse, error)
}
