package writerreader

import (
	tile38_entity "github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/xjem/t38c"
)

type writerReader struct {
	tile *t38c.Client
}

func NewTile38WriterReader(tile *t38c.Client) (repository.Tile38ReaderWriter, error) {
	return &writerReader{
		tile: tile,
	}, nil
}

//SetGeofencingKey...
func (wr *writerReader) SetGeofencingKey(payload tile38_entity.SetKey) error {
	if err := wr.tile.Keys.
		Set(payload.Key, payload.ObjectId).
		Point(payload.Lat, payload.Long).
		Field("speed", payload.Fields.Speed).
		Field("timestamp", payload.Fields.Timestamp).
		Field("device_type", payload.Fields.DeviceType).
		Field("id", payload.Fields.Id).
		Do(); err != nil {
		return err
	}
	return nil
}

// GetLastGeofencingDetect...
func (wr *writerReader) GetLastGeofencingDetect(key, objId string, withfield bool) (*t38c.GetResponse, error) {
	res, err := wr.tile.Keys.Get(key, objId, withfield)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// close...
func (wr *writerReader) Close() error {
	return wr.tile.Close()
}
