-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_geofence_name ON mst_geofence_area(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area DROP INDEX `idx_geofence_name`;
-- +goose StatementEnd
