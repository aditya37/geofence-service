-- +goose Up
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE area_type location_type BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE location_type area_type VARCHAR;
-- +goose StatementEnd
