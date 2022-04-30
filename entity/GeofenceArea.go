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
