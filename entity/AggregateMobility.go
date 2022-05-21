package entity

import "time"

type (
	AggregateMobility struct {
		Id         int64
		Enter      int
		Inside     int
		Exit       int
		GeofenceId int64
		CreatedAt  time.Time
		ModifiedAt time.Time
	}
	ResultGetDailyAvg struct {
		GeofencId int64
		Enter     float64
		Exit      float64
		Inside    float64
	}
)
