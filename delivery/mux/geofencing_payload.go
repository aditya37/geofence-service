package mux

type (
	RequestAddGeofence struct {
		LocationId   int64    `json:"location_id"`
		Name         string   `json:"name"`
		LocationType int64    `json:"location_type"`
		Detect       []string `json:"detect"`
		Geojson      string   `json:"geojson"`
		GeofenceType int64    `json:"geofence_type"`
	}
)
