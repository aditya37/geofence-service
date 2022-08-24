package geofence

import (
	"context"
	"database/sql"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/util"
)

func (gm *geofenceManager) GetCountGeofences(ctx context.Context) (*entity.ResultCounter, error) {

	row := gm.db.QueryRowContext(ctx, mysqlQueryGetCounter)

	var record entity.ResultCounter
	if err := row.Scan(
		&record.GeofenceArea,
	); err != nil {
		if err == sql.ErrNoRows {
			util.Logger().Error(err)
			return &entity.ResultCounter{}, nil
		}
		return nil, err
	}
	return &record, nil
}
