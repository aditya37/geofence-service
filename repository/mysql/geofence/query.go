package geofence

const (
	mysqlQueryInsertGeofence = `INSERT INTO mst_geofence_area (
		geofence_id,
		name,
		area_type,
		detect,
		geojson
	) VALUES(?,?,?,?,ST_GeomFromGeoJSON(?))`
	mysqlQueryGetGeofenceByName       = `SELECT id,geofence_id FROM mst_geofence_area WHERE name = ?`
	mysqlQueryGetGeofenceByGeofenceId = `SELECT id,geofence_id FROM mst_geofence_area WHERE id = ?`
)
