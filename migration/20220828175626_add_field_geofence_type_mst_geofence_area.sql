-- +goose Up
-- +goose StatementBegin
ALTER TABLE mst_geofence_area ADD COLUMN geofence_type BIGINT DEFAULT 3;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area DROP COLUMN geofence_type;
-- +goose StatementEnd
