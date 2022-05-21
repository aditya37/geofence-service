package repository

import (
	"context"
	"io"

	"github.com/aditya37/geofence-service/entity"
)

type MobilityManager interface {
	io.Closer
	GetLastAggregateFieldValue(ctx context.Context, field string, geofence_id int64) (int, error)
	UpdateAggregateFieldValue(ctx context.Context, field string, geofence_id int64, value int) error
	InsertDefaultValueAggregateField(ctx context.Context, field string, geofence_id int64, value int) error
	// Get all daily average
	GetDailyAverage(ctx context.Context, interval int) (*entity.ResultGetDailyAvg, error)
	GetAllAreaDailyAverage(ctx context.Context, interval int) ([]*entity.ResultGetDailyAvg, error)
}
