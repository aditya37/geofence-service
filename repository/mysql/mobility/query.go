package mobility

const (
	mysqlQueryGetCurrentMobilityCount    = "SELECT `%s` FROM mst_aggregate_mobility WHERE geofence_id = ? AND DATE(modified_at) = DATE(NOW())"
	mysqlQueryUpdateCurrentMobilityCount = "UPDATE mst_aggregate_mobility SET `%s` = ?,modified_at=NOW() WHERE geofence_id =? AND DATE(modified_at) = DATE(NOW())"
	mysqlQueryInsertDefaultFieldValue    = "INSERT INTO mst_aggregate_mobility(`%s`,geofence_id,modified_at) VALUES(?,?,NOW())"
	// query for count,avg or etc
	mysqlQueryGetDailyAvg = `SELECT 
		COALESCE(AVG(enter),0) AS enter,
		COALESCE(AVG("exit"),0) AS 'exit',
		COALESCE(AVG(inside),0) AS inside 
	FROM mst_aggregate_mobility 
	WHERE DATE(modified_at) > DATE(NOW()) - INTERVAL ? DAY`
	mysqlQueryGetAllAreaDailyAvg = `CALL get_all_area_daily_avg(?)`
)

/*
select ma.geofence_id,date(ma.modified_at) as modified_at,(
   select sum(me.enter) from mst_aggregate_mobility me where me.geofence_id = ma.geofence_id
) as avg_enter,(
   select sum(`exit`) from mst_aggregate_mobility me where me.geofence_id = ma.geofence_id
) as avg_exit,(
   select sum(inside) from mst_aggregate_mobility me where me.geofence_id = ma.geofence_id
) as avg_inside from mst_aggregate_mobility ma group by ma.geofence_id,ma.modified_at
*/
