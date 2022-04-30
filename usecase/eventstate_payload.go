package usecase

import (
	"errors"

	gjson "github.com/paulmach/go.geojson"
)

// error variable
var (
	ErrMessageDataNil = errors.New("message is nil")
)

// state
var (
	EventNotiferTypeSuccess  = "SUCCESS"
	EventNotiferTypeFailed   = "FAILED"
	EventStateInsert         = "INSERT"
	EventStateInsertSuccess  = "INSERT_SUCCESS"
	EventStateInsertRollback = "INSERT_ROLLBACK"
)

type (
	// payload publish add geofence
	GeofenceData struct {
		Id       string        `json:"id"`
		Name     string        `json:"name"`
		Detect   []string      `json:"detect"`
		AreaType string        `json:"area_type"`
		Shape    gjson.Feature `json:"feature,omitempty"`
	}
	Metadata struct {
		PublishedAt     string `json:"published_at"`
		LocationId      int64  `json:"location_id"`
		LocationName    string `json:"location_name"`
		LocationType    string `json:"location_type,omitempty"`
		SubLocationId   int64  `json:"sub_location_id"`
		SubLocationName string `json:"sub_location_name"`
		Message         string `json:"message"`
	}
	GeofenceEventState struct {
		EventId      string       `json:"event_id"`
		State        string       `json:"state"`
		ServiceName  string       `json:"service_name"`
		GeofenceData GeofenceData `json:"geofence_data,omitempty"`
		Metadata     Metadata     `json:"metadata"`
	}

	//GetServiceEventState...
	GetServiceEventStateRequest struct {
		ServiceName string `json:"service_name"`
		EventId     string `json:"event_id"`
	}
	GetServiceEventStateResponse struct {
		PublishedAt     string `json:"published_at"`
		LocationId      int64  `json:"location_id"`
		LocationName    string `json:"location_name"`
		LocationType    string `json:"location_type,omitempty"`
		SubLocationId   int64  `json:"sub_location_id"`
		SubLocationName string `json:"sub_location_name"`
		Message         string `json:"message"`
	}
)
