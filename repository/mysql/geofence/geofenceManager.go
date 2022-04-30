package geofence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/repository"
)

type geofenceManager struct {
	db *sql.DB
}

func NewGeofenceManager(db *sql.DB) (repository.GeofenceManager, error) {
	return &geofenceManager{
		db: db,
	}, nil
}

// InsertGeofenceArea....
func (gm *geofenceManager) InsertGeofenceArea(ctx context.Context, data entity.GeofenceArea) error {
	args := []interface{}{
		&data.GeofenceId,
		&data.Name,
		&data.AreaType,
		&data.Detect,
		&data.Geojson,
	}
	row, err := gm.db.ExecContext(ctx, mysqlQueryInsertGeofence, args...)
	if err != nil {
		return err
	}
	if isAffacted, _ := row.RowsAffected(); isAffacted == 0 {
		return errors.New("Failed to insert new geofence")
	}
	return nil
}

// Close....
func (gm *geofenceManager) Close() error {
	return gm.db.Close()
}
