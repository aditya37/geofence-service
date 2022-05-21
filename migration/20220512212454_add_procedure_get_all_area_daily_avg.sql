-- +goose Up
DROP PROCEDURE IF EXISTS `get_all_area_daily_avg`;
-- +goose StatementBegin
CREATE PROCEDURE `get_all_area_daily_avg` (	
	IN interval_day BIGINT
)
BEGIN
SELECT ma.geofence_id,(
	SELECT AVG(me.enter) FROM mst_aggregate_mobility me WHERE me.geofence_id = ma.geofence_id
  ) AS enter,(
    SELECT AVG(`exit`) FROM mst_aggregate_mobility me WHERE me.geofence_id = ma.geofence_id
  ) AS 'exit',(
    SELECT AVG(inside) FROM mst_aggregate_mobility me WHERE me.geofence_id = ma.geofence_id
  ) AS inside
FROM 
	mst_aggregate_mobility ma 
WHERE 
	DATE(ma.modified_at) > DATE(NOW()) - INTERVAL interval_day DAY 
GROUP BY 
	ma.geofence_id;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP PROCEDURE `get_all_area_daily_avg`;
-- +goose StatementEnd
