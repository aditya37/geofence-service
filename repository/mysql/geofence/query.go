package geofence

const (
	// Write....
	mysqlQueryInsertGeofence = `INSERT INTO mst_geofence_area (
		location_id,
		name,
		location_type,
		detect,
		geojson,
		channel_name
	) VALUES(?,?,?,?,ST_GeomFromGeoJSON(?),?)`

	// Read...
	mysqlQueryGetGeofenceByChannelName = `SELECT id,location_id FROM mst_geofence_area WHERE channel_name = ?`
	mysqlQueryGetGeofenceByGeofenceId  = `SELECT id,location_id FROM mst_geofence_area WHERE id = ?`
	mysqlQueryGetGeofenceTypeById      = `SELECT id,type_name FROM mst_geofence_type WHERE id = ?`
	mysqlQueryGetGeofenceTypeByName    = `SELECT id,type_name FROM mst_geofence_type WHERE type_name = ?`
)
