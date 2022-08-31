package geofence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aditya37/geofence-service/entity"
)

func (gm *geofenceManager) GetGeofenceAreaByLocationId(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error) {
	arg := []interface{}{
		&id,
	}
	row := gm.db.QueryRowContext(ctx, mysqlQueryGetGeofenceByLocationId, arg...)

	var result entity.ResultGetGeofenceById
	if err := row.Scan(
		&result.Id,
		&result.LocationId,
		&result.Name,
		&result.Detect,
		&result.ChannelName,
		&result.Geojson,
		&result.TypeName,
		&result.AvgMobility,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("geofence area not found")
		}
		return nil, err
	}
	return &result, nil
}
