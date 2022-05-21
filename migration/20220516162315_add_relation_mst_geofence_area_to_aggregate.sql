-- +goose Up
-- +goose StatementBegin
ALTER TABLE 
	mst_aggregate_mobility 
ADD CONSTRAINT 
	`fk_geofence_area_aggregate`  
FOREIGN KEY (`geofence_id`) REFERENCES `mst_geofence_area` (`id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_aggregate_mobility DROP CONSTRAINT fk_geofence_area_aggregate;
-- +goose StatementEnd
