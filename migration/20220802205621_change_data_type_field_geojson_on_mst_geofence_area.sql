-- +goose Up
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE geojson geojson GEOMETRY;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mst_geofence_area CHANGE geojson geojson GEOMETRY;
-- +goose StatementEnd
