-- +goose Up
-- +goose StatementBegin
CREATE TABLE `mst_aggregate_mobility`(
	id BIGINT NOT NULL AUTO_INCREMENT,
	enter INT NOT NULL DEFAULT '0',
	inside INT NOT NULL DEFAULT '0',
	`exit` INT NOT NULL DEFAULT '0',
	geofence_id BIGINT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	modified_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `mst_aggregate_mobility`;
-- +goose StatementEnd
