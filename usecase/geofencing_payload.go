package usecase

import (
	"fmt"

	geospatialSrv "github.com/aditya37/api-contract/geospatial-service/service"
	"github.com/xjem/t38c"
)

type GeometryType string

var (
	// Geometry Type
	LineString GeometryType = "LineString"
	Polygon    GeometryType = "Polygon"
	Point      GeometryType = "Point"
)

func (gt GeometryType) ToString() string {
	return fmt.Sprintf("%s", gt)
}

type (
	// TODO: Refactor this after service geotracking
	DeviceMetadata struct {
		DeviceId string `json:"device_id"`
		Source   string `json:"source"`
	}
	TrackingPayload struct {
		Speed     int64          `json:"speed"`
		Lat       float64        `json:"lat"`
		Long      float64        `json:"long"`
		Timestamp int64          `json:"timestamp"`
		Device    DeviceMetadata `json:"device_metadata"`
	}

	// notify payload
	Detect struct {
		Enter  float64 `json:"enter"`
		Inside float64 `json:"inside"`
		Exit   float64 `json:"exit"`
		Date   string  `json:"date"`
	}
	MobilityArea struct {
		LocationId   int64  `json:"location_id,omitempty"`
		GeofenceId   int64  `json:"geofence_id,omitempty"`
		LocationName string `json:"location_name,omitempty"`
		Average      Detect `json:"average,omitempty"`
	}
	Mobility struct {
		DailyAverage  Detect         `json:"daily_average,omitempty"`
		MobilityAreas []MobilityArea `json:"mobility_area,omitempty"`
	}

	// payload for geofencing type 'tourist'
	NotifyTouristPayload struct {
		Detect         string                                  `json:"detect"`
		Mobility       Mobility                                `json:"mobility,omitempty"`
		NearbyLocation geospatialSrv.GetNearbyLocationResponse `json:"nearby_location,omitempty"`
		ChannelName    string                                  `json:"channel_name"`
		Object         *t38c.Object                            `json:"object"`
	}

	// AddGeofence
	RequestAddGeofence struct {
		LocationId   int64    `json:"location_id"`
		Name         string   `json:"name"`
		LocationType int64    `json:"location_type"`
		Detect       []string `json:"detect"`
		Geojson      string   `json:"geojson"`
		GeofenceType int64    `json:"geofence_type"`
		ChannelName  string   `json:"-"`
	}
	ResponseAddGeofence struct {
		Message    string `json:"message"`
		Name       string `json:"name"`
		LocationId int64  `json:"location_id"`
		CreatedAt  string `json:"created_at"`
	}

	// RequestGetGeofenceTypeDetail....
	RequestGetGeofenceTypeDetail struct {
		TypeName string `json:"type_name"`
		TypeId   int64  `json:"type_id"`
	}

	ResponseGetGeofenceTypeDetail struct {
		Id       int64  `json:"id"`
		TypeName string `json:"type_name"`
	}
)
