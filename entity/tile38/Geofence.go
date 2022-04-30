package tile38

import geojson "github.com/paulmach/go.geojson"

type Geofence struct {
	Name    string
	Key     string
	Detect  []string
	Action  string
	Feature *geojson.Feature
}
