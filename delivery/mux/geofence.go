package mux

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/usecase/geofencing"
	"github.com/aditya37/geofence-service/util"
)

type GeofenceDelivery struct {
	geofenceCase *geofencing.GeofencingUsecase
}

func NewGeofencingDeliver(gu *geofencing.GeofencingUsecase) *GeofenceDelivery {
	return &GeofenceDelivery{
		geofenceCase: gu,
	}
}

func (gd *GeofenceDelivery) AddGeofence(w http.ResponseWriter, r *http.Request) {
	var request RequestAddGeofence
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		})
		return
	}

	convertedGeojson := util.ReplaceGeojsonSingleQuote(request.Geojson)

	resp, err := gd.geofenceCase.AddGeofence(
		r.Context(),
		usecase.RequestAddGeofence{
			LocationId:   request.LocationId,
			Name:         request.Name,
			LocationType: request.LocationType,
			Detect:       request.Detect,
			GeofenceType: request.GeofenceType,
			Geojson:      convertedGeojson,
		},
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusCreated, resp)
}

func (gd *GeofenceDelivery) GetGeofenceTypeDetail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query()["id"]
	name := r.URL.Query()["name"]

	intId, _ := strconv.Atoi(id[0])

	resp, err := gd.geofenceCase.GetGeofenceTypeDetail(
		r.Context(),
		usecase.RequestGetGeofenceTypeDetail{
			TypeName: name[0],
			TypeId:   int64(intId),
		},
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
	return
}
