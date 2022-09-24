package geofencing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetAvgMobililtyByArea(ctx context.Context, request usecase.RequestGetAvgMobililtyByArea) (usecase.ResponseGetAvgMobililtyByArea, error) {

	// validate is available area...
	area, err := gu.geofenceManager.DetailGeofenceAreaById(ctx, request.GeofenceId)
	if err != nil {
		if err.Error() == "Geofence area not found" {
			return usecase.ResponseGetAvgMobililtyByArea{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  err.Error(),
			}

		}
		return usecase.ResponseGetAvgMobililtyByArea{}, &util.ErrorMsg{
			HttpRespCode: http.StatusUnprocessableEntity,
			Description:  err.Error(),
		}
	}
	// only show data if goefence area tourist...
	if isgeofTourist := gu.isTouristArea(area); !isgeofTourist {
		return usecase.ResponseGetAvgMobililtyByArea{
			Average: []usecase.Detect{},
		}, nil
	}
	resp, err := gu.mobilityManager.GetMobilityAverageByArea(ctx, request.Interval, request.GeofenceId)
	if err != nil {
		util.Logger().Error(err)
		return usecase.ResponseGetAvgMobililtyByArea{}, &util.ErrorMsg{
			HttpRespCode: http.StatusUnprocessableEntity,
			Description:  err.Error(),
		}
	}
	var item []usecase.Detect
	for _, v := range resp {
		y, m, d := v.ModifiedAt.Date()
		item = append(item, usecase.Detect{
			Enter:  float64(v.Enter),
			Exit:   float64(v.Exit),
			Inside: float64(v.Inside),
			Date:   fmt.Sprintf("%d-%d-%d", d, m, y),
		})
	}

	return usecase.ResponseGetAvgMobililtyByArea{
		IsTourist: true,
		Average:   item,
	}, nil
}

// check area type...
func (gu *GeofencingUsecase) isTouristArea(data *entity.ResultGetGeofenceById) bool {
	if data.TypeName != "tourist" {
		return false
	} else {
		return true
	}
}
