package geofencing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetGeofenceByLocationId(ctx context.Context, location_id int64) (usecase.ResponseGetGeofenceById, error) {
	resp, err := gu.geofenceManager.GetGeofenceAreaByLocationId(ctx, location_id)
	if err != nil {
		if err.Error() == "geofence area not found" {
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
