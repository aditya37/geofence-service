package geofence

/*
	&data.GeofenceId,
		&data.Name,
		&data.AreaType,
		&data.Detect,
		&data.Geojson,
*/
const (
	mysqlQueryInsertGeofence = `INSERT INTO mst_geofence_area (
		geofence_id,
		name,
		area_type,
		detect,
		geojson
	) VALUES(?,?,?,?,ST_GeomFromGeoJSON(?))`
)
