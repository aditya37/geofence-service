-- +goose Up
-- +goose StatementBegin
ALTER TABLE 
	mst_geofence_area 
ADD CONSTRAINT 
	`fk_geofence_area_type`  
FOREIGN KEY (`geofence_type`) REFERENCES `mst_geofence_type` (`id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area DROP CONSTRAINT fk_geofence_area_type;
-- +goose StatementEnd
