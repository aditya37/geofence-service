-- +goose Up
-- +goose StatementBegin
CREATE TABLE `mst_geofence_type` (
	id BIGINT NOT NULL AUTO_INCREMENT,
	type_name VARCHAR(32),
	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `mst_geofence_type`;
-- +goose StatementEnd
