package entity

import (
	"encoding/json"
	"time"
)

type GeofenceArea struct {
	Id           int64
	LocationId   int64
	Name         string
	LocationType int64
	Detect       json.RawMessage
	Geojson      json.RawMessage
	CreatedAt    time.Time
	ModifiedAt   time.Time
	GeofenceType int64
	ChannelName  string
}
type ResultGetGeofenceById struct {
	Id          int64
	LocationId  int64
	Name        string
	Detect      json.RawMessage
	ChannelName string
	Geojson     json.RawMessage
	TypeName    string
	AvgMobility float64
}

// redis struct/payload data
type ResultGetLocationDetailByGeofenceId struct {
	LocatioName string `json:"location_name"`
	LocationId  int64  `json:"location_id"`
}

// result counter...
type ResultCounter struct {
	GeofenceArea int64
}
