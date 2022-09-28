package repository

import (
	"context"
	"io"

	"github.com/aditya37/geofence-service/entity"
)

type GeofenceManager interface {
	io.Closer
	// Writer....
	InsertGeofenceArea(ctx context.Context, data entity.GeofenceArea) error
	UpdateLocationToGeofence(ctx context.Context, location_id int64) error
	// Reader....
	DetailGeofenceAreaByChannelName(ctx context.Context, name string) (*entity.GeofenceArea, error)
	DetailGeofenceAreaById(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error)
	GetGeofenceTypeById(ctx context.Context, model entity.GeofenceType) (*entity.GeofenceType, error)
	GetGeofenceTypeByName(ctx context.Context, model entity.GeofenceType) (*entity.GeofenceType, error)
	GetGeofenceAreaByLocationId(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error)
	GetGeofenceAreaByType(ctx context.Context, model entity.GeofenceArea) ([]*entity.GeofenceArea, error)

	// get count geofence area etc...
	GetCountGeofences(ctx context.Context) (*entity.ResultCounter, error)
}
