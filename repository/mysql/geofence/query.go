package geofence

const (
	// Write....
	mysqlQueryInsertGeofence = `INSERT INTO mst_geofence_area (
		location_id,
		name,
		location_type,
		detect,
		geojson,
		channel_name,
		geofence_type
	) VALUES(?,?,?,?,ST_GeomFromGeoJSON(?),?,?)`
	mysqlQueryUpdateLocationToGeofence = `UPDATE %s.mst_location SET is_geofence = 1 WHERE id = ?`
	// Read...
	mysqlQueryGetGeofenceByChannelName = `SELECT id,location_id FROM mst_geofence_area WHERE channel_name = ?`
	mysqlQueryGetGeofenceTypeById      = `SELECT id,type_name FROM mst_geofence_type WHERE id = ?`
	mysqlQueryGetGeofenceTypeByName    = `SELECT id,type_name FROM mst_geofence_type WHERE type_name = ?`
	mysqlQueryGetCounter               = `SELECT COUNT(id) AS geofence_area FROM mst_geofence_area`

	mysqlQueryGetGeofenceByGeofenceId = `SELECT
			 mga.id,
			 mga.location_id,
			 mga.name,
			 mga.detect,
			 mga.channel_name,
			 ST_ASGEOJSON(mga.geojson) AS geojson,
			 mgt.type_name,
			 COALESCE(
			 	ROUND(AVG(mam.inside+mam.enter+mam.` + `exit` + `)),0
			 ) AS avg_mobility
		FROM mst_geofence_area mga
		INNER JOIN mst_geofence_type mgt ON mga.geofence_type = mgt.id
		LEFT JOIN mst_aggregate_mobility mam ON mam.geofence_id = mga.id
		WHERE mga.id = ? GROUP BY mga.id`
	mysqlQueryGetGeofenceByLocationId = `SELECT
			 mga.id,
			 mga.location_id,
			 mga.name,
			 mga.detect,
			 mga.channel_name,
			 ST_ASGEOJSON(mga.geojson) AS geojson,
			 mgt.type_name,
			 COALESCE(
			 	ROUND(AVG(mam.inside+mam.enter+mam.` + `exit` + `)),0
			 ) AS avg_mobility
		FROM mst_geofence_area mga
		INNER JOIN mst_geofence_type mgt ON mga.geofence_type = mgt.id
		LEFT JOIN mst_aggregate_mobility mam ON mam.geofence_id = mga.id
		WHERE mga.location_id = ? GROUP BY mga.id`
	mysqlQueryGetGeofenceAreaByType = `SELECT
			 mga.id,
			 mga.location_id,
			 mga.name,
			 mga.detect,
			 mga.channel_name,
			 ST_ASGEOJSON(mga.geojson) AS geojson,
			 mgt.type_name
		FROM mst_geofence_area mga
		INNER JOIN mst_geofence_type mgt ON mga.geofence_type = mgt.id WHERE mgt.type_name = ? LIMIT ?,?`
)
