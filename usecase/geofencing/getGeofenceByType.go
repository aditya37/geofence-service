package geofencing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	"github.com/tidwall/geojson"
)

func (gu *GeofencingUsecase) GetGeofenceByType(ctx context.Context, request usecase.RequestGetGeofenceByType) (usecase.ResponseGetGeofenceByType, error) {
	if _, err := gu.geofenceManager.GetGeofenceTypeByName(
		ctx,
		entity.GeofenceType{
			TypeName: request.Type,
		},
	); err != nil {
		if err.Error() == "Geofence type not found" {
			return usecase.ResponseGetGeofenceByType{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  err.Error(),
			}
		}
		return usecase.ResponseGetGeofenceByType{}, &util.ErrorMsg{
			HttpRespCode: http.StatusUnprocessableEntity,
			Description:  err.Error(),
		}
	}

	resp, err := gu.geofenceManager.GetGeofenceAreaByType(
		ctx,
		entity.GeofenceArea{
			Page:        request.Page,
			ItemPerPage: request.ItemPerPage,
			Type:        request.Type,
		},
	)
	if err != nil {
		return usecase.ResponseGetGeofenceByType{}, &util.ErrorMsg{
			HttpRespCode: http.StatusUnprocessableEntity,
			Description:  err.Error(),
		}
	}

	var result []usecase.ResponseGetGeofenceById
	for _, v := range resp {
		var d []string
		if err := json.Unmarshal(v.Detect, &d); err != nil {
			util.Logger().Error(err)
			continue
		}

		// generate geojsonfeature...
		gjsonParse, _ := geojson.Parse(string(v.Geojson), nil)
		feature := geojson.NewFeature(gjsonParse, "")
		result = append(result, usecase.ResponseGetGeofenceById{
			Id:          v.Id,
			LocationId:  v.LocationId,
			Name:        v.Name,
			Detect:      d,
			ChannelName: v.ChannelName,
			Geojson:     feature.String(),
			TypeName:    v.Type,
			AvgMobility: 0,
		})
	}
	return usecase.ResponseGetGeofenceByType{
		GeofenceAreas: result,
	}, nil
}
