package geofence

import (
	"context"

	"github.com/aditya37/geofence-service/entity"
)

func (gm *geofenceManager) GetGeofenceAreaByType(ctx context.Context, model entity.GeofenceArea) ([]*entity.GeofenceArea, error) {
	offset := (model.Page - 1) * model.ItemPerPage
	arg := []interface{}{
		&model.Type,
		&offset,
		&model.ItemPerPage,
	}
	rows, err := gm.db.QueryContext(ctx, mysqlQueryGetGeofenceAreaByType, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.GeofenceArea
	for rows.Next() {
		var record entity.GeofenceArea
		if err := rows.Scan(
			&record.Id,
			&record.LocationId,
			&record.Name,
			&record.Detect,
			&record.ChannelName,
			&record.Geojson,
			&record.Type,
		); err != nil {
			return nil, err
		}
		result = append(result, &record)
	}
	return result, nil
}
