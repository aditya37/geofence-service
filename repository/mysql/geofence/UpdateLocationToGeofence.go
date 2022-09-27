package geofence

import (
	"context"
	"fmt"

	getenv "github.com/aditya37/get-env"
)

func (gm *geofenceManager) UpdateLocationToGeofence(ctx context.Context, location_id int64) error {
	arg := []interface{}{
		&location_id,
	}
	query := fmt.Sprintf(mysqlQueryUpdateLocationToGeofence, getenv.GetString("DB_GEOSPATIAL", "db_geospatial"))
	if _, err := gm.db.ExecContext(ctx, query, arg...); err != nil {
		return err
	}

	return nil
}
