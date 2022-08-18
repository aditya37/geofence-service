-- +goose Up
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE geofence_id location_id BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE location_id geofence_id VARCHAR;
-- +goose StatementEnd
