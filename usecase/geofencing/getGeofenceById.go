package geofencing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetGeofenceAreaById(ctx context.Context, geofenceId int64) (usecase.ResponseGetGeofenceById, error) {
	resp, err := gu.geofenceManager.DetailGeofenceAreaById(
		ctx,
		geofenceId,
	)
	if err != nil {
		if err.Error() == "Geofence area not found" {
			return usecase.ResponseGetGeofenceById{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  err.Error(),
			}
		}
		return usecase.ResponseGetGeofenceById{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	var detect []string
	if err := json.Unmarshal(resp.Detect, &detect); err != nil {
		util.Logger().Error(err)
		return usecase.ResponseGetGeofenceById{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}

	return usecase.ResponseGetGeofenceById{
		Id:          resp.Id,
		LocationId:  resp.LocationId,
		Name:        resp.Name,
		ChannelName: resp.ChannelName,
		Detect:      detect,
		Geojson:     string(resp.Geojson),
		TypeName:    resp.TypeName,
		AvgMobility: resp.AvgMobility,
	}, nil
}
