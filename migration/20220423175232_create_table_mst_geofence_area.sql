-- +goose Up
-- +goose StatementBegin
CREATE TABLE `mst_geofence_area` (
	id BIGINT NOT NULL AUTO_INCREMENT,
	geofence_id VARCHAR(50) NOT NULL,
	name VARCHAR(50) NOT NULL,
	area_type VARCHAR(50) NOT NULL,
	detect JSON NOT NULL,
	geojson POLYGON,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `mst_geofence_area`
-- +goose StatementEnd
