-- +goose Up
-- +goose StatementBegin
ALTER TABLE mst_geofence_area ADD COLUMN channel_name VARCHAR(300);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area DROP COLUMN channel_name;
-- +goose StatementEnd
