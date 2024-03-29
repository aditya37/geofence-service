package usecase

import (
	"fmt"

	geospatialSrv "github.com/aditya37/api-contract/geospatial-service/service"
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

//default value if device type 0
//because tile38 fields not support zero value
var DeviceTypeZero = 11

type (

	// MQTTRespTracking....
	// Payload for send to device or mqtt
	Message struct {
		Value  string `json:"value"`
		Reason string `json:"reason"`
	}
	GPSData struct {
		Lat      float64 `json:"lat"`
		Long     float64 `json:"long"`
		Altitude float64 `json:"altitude"`
		Speed    float64 `json:"speed"`
		Angle    float64 `json:"angle"`
	}
	Sensor struct {
		Temp   float64 `json:"temp"`
		Signal float64 `json:"signal"`
	}
	MQTTRespTracking struct {
		DeviceId    string  `json:"device_id"`
		DeviceType  int     `json:"device_type"`
		Id          int64   `json:"id"`
		Status      string  `json:"status"`
		RespMessage Message `json:"message"`
		GPSData     GPSData `json:"gps_data"`
		Sensors     Sensor  `json:"sensor_data"`
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
	NotifyGeofencingPayload struct {
		Type           string                                  `json:"type"`
		Detect         string                                  `json:"detect"`
		Mobility       Mobility                                `json:"mobility,omitempty"`
		NearbyLocation geospatialSrv.GetNearbyLocationResponse `json:"nearby_location,omitempty"`
		ChannelName    string                                  `json:"channel_name"`
		Object         string                                  `json:"object"`
		DeviceId       string                                  `json:"device_id"`
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

	// GetCounts...
	ResponseGetCounts struct {
		GeofenceArea int64 `json:"gefence_area"`
	}

	// ResponseGetGeofenceById...
	ResponseGetGeofenceById struct {
		Id          int64    `json:"id"`
		LocationId  int64    `json:"location_id"`
		Name        string   `json:"name"`
		Detect      []string `json:"detect"`
		ChannelName string   `json:"channel_name"`
		Geojson     string   `json:"geojson"`
		TypeName    string   `json:"type_name"`
		AvgMobility float64  `json:"avg_mobility"`
	}

	// RequestGetGoefenceByType...
	RequestGetGeofenceByType struct {
		Type        string `json:"type"`
		Page        int    `json:"page"`
		ItemPerPage int    `json:"item_perpage"`
	}
	// ResponseGetGeofenceByType...
	ResponseGetGeofenceByType struct {
		GeofenceAreas []ResponseGetGeofenceById `json:"geofence_areas"`
	}

	//RequestGetAvgMobililtyByInterval
	RequestGetAvgMobililtyByArea struct {
		GeofenceId int64 `json:"geofence_id"`
		Interval   int64 `json:"interval"`
	}
	ResponseGetAvgMobililtyByArea struct {
		IsTourist bool     `json:"is_tourist"`
		Average   []Detect `json:"average"`
	}

	// responseQaToolPublishGeofence
	ResponseQAToolPublishGeofence struct {
		Message string `json:"message"`
	}
	PayloadInsertDeviceDetect struct {
		DeviceId   int64   `json:"device_id"`
		Detect     string  `json:"detect"`
		Lat        float64 `json:"lat"`
		Long       float64 `json:"long"`
		DetectTime int64   `json:"detect_time"`
	}
)
