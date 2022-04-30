-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX idx_geofence_id ON mst_geofence_area(geofence_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area DROP INDEX `idx_geofence_id`;
-- +goose StatementEnd
