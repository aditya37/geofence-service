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

// DetailGeofenceAreaByName...
func (gm *geofenceManager) DetailGeofenceAreaByChannelName(ctx context.Context, name string) (*entity.GeofenceArea, error) {
	arg := []interface{}{
		&name,
	}
	row := gm.db.QueryRowContext(ctx, mysqlQueryGetGeofenceByChannelName, arg...)
	var result entity.GeofenceArea
	if err := row.Scan(
		&result.Id,
		&result.LocationId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Geofence Area not found")
		}
		return nil, err
	}
	return &result, nil
}

// DetailGeofenceAreaByGeofenceId...
func (gm *geofenceManager) DetailGeofenceAreaById(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error) {
	arg := []interface{}{
		&id,
	}
	row := gm.db.QueryRowContext(ctx, mysqlQueryGetGeofenceByGeofenceId, arg...)
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
			return nil, errors.New("Geofence area not found")
		}
		return nil, err
	}
	return &result, nil
}

// InsertGeofenceArea....
func (gm *geofenceManager) InsertGeofenceArea(ctx context.Context, data entity.GeofenceArea) error {
	args := []interface{}{
		&data.LocationId,
		&data.Name,
		&data.LocationType,
		&data.Detect,
		&data.Geojson,
		&data.ChannelName,
		&data.GeofenceType,
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
