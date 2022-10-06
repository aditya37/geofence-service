package mux

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/usecase/geofencing"
	"github.com/aditya37/geofence-service/util"
	"github.com/gorilla/mux"
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

// getGeofenceCounts...
func (gd *GeofenceDelivery) GetGeofenceCount(w http.ResponseWriter, r *http.Request) {
	resp, err := gd.geofenceCase.GetCounts(r.Context())
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}

	EncodeResponse(w, http.StatusOK, resp)
	return
}

// getGeofenceById
func (gd *GeofenceDelivery) GetGeofenceById(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	id := param["id"]
	paramId, _ := strconv.Atoi(id)
	resp, err := gd.geofenceCase.GetGeofenceAreaById(
		r.Context(),
		int64(paramId),
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
	return
}

// getGeofenceByLocationId...
func (gd *GeofenceDelivery) GetGeofenceByLocationId(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	locationId := param["location_id"]
	intLocationId, _ := strconv.Atoi(locationId)
	resp, err := gd.geofenceCase.GetGeofenceByLocationId(
		r.Context(),
		int64(intLocationId),
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
}

// getGeofenceByType....
func (gd *GeofenceDelivery) GetGeofenceByType(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	geofenceType, ok := param["type"]
	if !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set geofence type",
		})
		return
	}
	// get query param...
	qry := r.URL.Query()
	if _, ok := qry["page"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set page value",
		})
		return
	}
	if _, ok := qry["itemPerPage"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set item per page value",
		})
		return
	}
	page, _ := strconv.Atoi(qry["page"][0])
	itemPerPage, _ := strconv.Atoi(qry["itemPerPage"][0])

	// usecase
	resp, err := gd.geofenceCase.GetGeofenceByType(
		r.Context(),
		usecase.RequestGetGeofenceByType{
			Type:        geofenceType,
			Page:        page,
			ItemPerPage: itemPerPage,
		},
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)

}

func (gd *GeofenceDelivery) GetAggregateMobilityByArea(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	geofenceId, ok := param["geofence_id"]
	if !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set geofence type",
		})
		return
	}
	qry := r.URL.Query()
	if _, ok := qry["interval"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set interval",
		})
		return
	}
	interval, _ := strconv.Atoi(qry["interval"][0])
	id, _ := strconv.Atoi(geofenceId)
	resp, err := gd.geofenceCase.GetAvgMobililtyByArea(
		r.Context(),
		usecase.RequestGetAvgMobililtyByArea{
			Interval:   int64(interval),
			GeofenceId: int64(id),
		},
	)
	if err != nil {
		util.Logger().Error(err)
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
}

func (gd *GeofenceDelivery) QaToolGeofence(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if _, ok := q["speed"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set speed",
		})
		return
	}
	if _, ok := q["lat"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set lat",
		})
		return
	}
	if _, ok := q["long"]; !ok {
		EncodeErrorResponse(w, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "please set long",
		})
		return
	}

	intSpeed, _ := strconv.Atoi(q["speed"][0])
	fLat, _ := strconv.ParseFloat(q["lat"][0], 64)
	fLong, _ := strconv.ParseFloat(q["long"][0], 64)
	resp, err := gd.geofenceCase.QAToolPublishGeofence(
		context.Background(),
		usecase.TrackingPayload{
			Lat:       fLat,
			Long:      fLong,
			Speed:     int64(intSpeed),
			Timestamp: time.Now().Unix(),
			Device: usecase.DeviceMetadata{
				DeviceId: q["device"][0],
			},
		},
	)
	if err != nil {
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
}
func (gd *GeofenceDelivery) GetGeofenceTypes(w http.ResponseWriter, r *http.Request) {
	resp, err := gd.geofenceCase.GetGeofenceTypes(r.Context())
	if err != nil {
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
}
