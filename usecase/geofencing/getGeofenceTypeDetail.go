package geofencing

import (
	"context"
	"net/http"
	"strings"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) GetGeofenceTypeDetail(ctx context.Context, request usecase.RequestGetGeofenceTypeDetail) (usecase.ResponseGetGeofenceTypeDetail, error) {
	if request.TypeId != 0 && request.TypeName == "" {
		resp, err := gu.processGetGefenceTypeById(ctx, request)
		if err != nil {
			return usecase.ResponseGetGeofenceTypeDetail{}, err
		}
		return resp, nil

	} else if request.TypeId == 0 && request.TypeName != "" {
		resp, err := gu.processGetGefenceTypeByName(ctx, request)
		if err != nil {
			return usecase.ResponseGetGeofenceTypeDetail{}, err
		}
		return resp, nil
	} else {
		return usecase.ResponseGetGeofenceTypeDetail{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "unknown request for get detail geofence type",
		}
	}
}

// processGetGefenceTypeByName
func (gu *GeofencingUsecase) processGetGefenceTypeByName(ctx context.Context, request usecase.RequestGetGeofenceTypeDetail) (usecase.ResponseGetGeofenceTypeDetail, error) {
	resp, err := gu.geofenceManager.GetGeofenceTypeByName(
		ctx,
		entity.GeofenceType{
			TypeName: request.TypeName,
		},
	)
	if err != nil {
		if notFound := strings.Contains(err.Error(), "Geofence type not found"); notFound {
			return usecase.ResponseGetGeofenceTypeDetail{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  "location type not found",
			}
		}
		return usecase.ResponseGetGeofenceTypeDetail{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return usecase.ResponseGetGeofenceTypeDetail{
		Id:       resp.Id,
		TypeName: resp.TypeName,
	}, nil

}

// processGetGefenceTypeById
func (gu *GeofencingUsecase) processGetGefenceTypeById(ctx context.Context, request usecase.RequestGetGeofenceTypeDetail) (usecase.ResponseGetGeofenceTypeDetail, error) {
	resp, err := gu.geofenceManager.GetGeofenceTypeById(
		ctx,
		entity.GeofenceType{
			Id: request.TypeId,
		},
	)
	if err != nil {
		if notFound := strings.Contains(err.Error(), "Geofence type not found"); notFound {
			return usecase.ResponseGetGeofenceTypeDetail{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  "location type not found",
			}
		}
		return usecase.ResponseGetGeofenceTypeDetail{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return usecase.ResponseGetGeofenceTypeDetail{
		Id:       resp.Id,
		TypeName: resp.TypeName,
	}, nil
}
