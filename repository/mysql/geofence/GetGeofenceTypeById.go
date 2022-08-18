package geofence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aditya37/geofence-service/entity"
)

func (gm *geofenceManager) GetGeofenceTypeById(ctx context.Context, model entity.GeofenceType) (*entity.GeofenceType, error) {
	arg := []interface{}{
		&model.Id,
	}
	row := gm.db.QueryRowContext(ctx, mysqlQueryGetGeofenceTypeById, arg...)

	var record entity.GeofenceType
	if err := row.Scan(
		&record.Id,
		&record.TypeName,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Geofence type not found")
		}
		return nil, err
	}
	return &record, nil
}
