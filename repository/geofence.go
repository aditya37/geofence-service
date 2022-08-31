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

	// Reader....
	DetailGeofenceAreaByChannelName(ctx context.Context, name string) (*entity.GeofenceArea, error)
	DetailGeofenceAreaById(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error)
	GetGeofenceTypeById(ctx context.Context, model entity.GeofenceType) (*entity.GeofenceType, error)
	GetGeofenceTypeByName(ctx context.Context, model entity.GeofenceType) (*entity.GeofenceType, error)
	GetGeofenceAreaByLocationId(ctx context.Context, id int64) (*entity.ResultGetGeofenceById, error)
	// get count geofence area etc...
	GetCountGeofences(ctx context.Context) (*entity.ResultCounter, error)
}
