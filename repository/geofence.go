package repository

import (
	"context"
	"io"

	"github.com/aditya37/geofence-service/entity"
)

type GeofenceManager interface {
	io.Closer
	InsertGeofenceArea(ctx context.Context, data entity.GeofenceArea) error
}
