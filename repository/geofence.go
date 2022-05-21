package repository

import (
	"context"
	"io"

	"github.com/aditya37/geofence-service/entity"
)

type GeofenceManager interface {
	io.Closer
	InsertGeofenceArea(ctx context.Context, data entity.GeofenceArea) error
	DetailGeofenceAreaByName(ctx context.Context, name string) (*entity.GeofenceArea, error)
	DetailGeofenceAreaById(ctx context.Context, id int64) (*entity.GeofenceArea, error)
}
