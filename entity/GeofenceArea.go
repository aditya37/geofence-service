package entity

import (
	"encoding/json"
	"time"
)

type GeofenceArea struct {
	Id         int64
	GeofenceId string
	Name       string
	AreaType   string
	Detect     json.RawMessage
	Geojson    json.RawMessage
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// redis struct/payload data
type ResultGetLocationDetailByGeofenceId struct {
	LocatioName string `json:"location_name"`
	LocationId  int64  `json:"location_id"`
}
