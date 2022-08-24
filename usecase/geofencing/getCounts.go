package geofencing

import (
	"context"
	"net/http"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetCounts(ctx context.Context) (usecase.ResponseGetCounts, error) {
	counts, err := gu.geofenceManager.GetCountGeofences(ctx)
	if err != nil {
		return usecase.ResponseGetCounts{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return usecase.ResponseGetCounts{
		GeofenceArea: counts.GeofenceArea,
	}, nil
}
