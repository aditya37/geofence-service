package geofencing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
	paulgeojson "github.com/paulmach/go.geojson"
	tidwal_gjson "github.com/tidwall/geojson"
	"github.com/tidwall/gjson"
)

var (
	// prefix for channel name in tile38...
	// geofencing.geofence_type.name
	prefixChanGeofencing = "geofencing.%s.%s"
)

func (gu *GeofencingUsecase) AddGeofence(ctx context.Context, request usecase.RequestAddGeofence) (usecase.ResponseAddGeofence, error) {
	// parse and validate geojson...
	parseGjson, err := tidwal_gjson.Parse(request.Geojson, nil)
	if err != nil {
		util.Logger().Error(err)
		return usecase.ResponseAddGeofence{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	if validGeoJson := parseGjson.Valid(); !validGeoJson {
		return usecase.ResponseAddGeofence{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "Geojson not valid",
		}
	}

	geofType, err := gu.geofenceManager.GetGeofenceTypeById(
		ctx,
		entity.GeofenceType{
			Id: request.GeofenceType,
		},
	)
	if err != nil {
		if typeNotFound := strings.Contains(err.Error(), ""); typeNotFound {
			return usecase.ResponseAddGeofence{}, &util.ErrorMsg{
				HttpRespCode: http.StatusNotFound,
				Description:  err.Error(),
			}
		}
		return usecase.ResponseAddGeofence{}, &util.ErrorMsg{
			HttpRespCode: http.StatusInternalServerError,
			Description:  err.Error(),
		}
	}

	// tile38 channel name...
	channelName := fmt.Sprintf(prefixChanGeofencing, geofType.TypeName, request.Name)

	if err := gu.insertGeofenceArea(
		ctx,
		usecase.RequestAddGeofence{
			LocationId:   request.LocationId,
			Name:         request.Name,
			LocationType: request.LocationType,
			ChannelName:  channelName,
			Detect:       request.Detect,
			Geojson:      request.Geojson,
			GeofenceType: geofType.Id,
		},
	); err != nil {
		return usecase.ResponseAddGeofence{}, err
	}

	// create channel to tile38...
	if err := gu.createGeofenceChannel(
		ctx, usecase.RequestAddGeofence{
			ChannelName: channelName,
			Detect:      request.Detect,
			Geojson:     request.Geojson,
		},
	); err != nil {
		return usecase.ResponseAddGeofence{}, err
	}

	return usecase.ResponseAddGeofence{
		Message:    "Success Add Geofence Area",
		Name:       request.Name,
		LocationId: request.LocationId,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}, nil
}

// parseGeojsonGeometry...
func (gu *GeofencingUsecase) parseGeojsonGeometry(request usecase.RequestAddGeofence) (string, error) {
	// parse data
	parsegeo := gjson.Parse(request.Geojson)
	geotype := parsegeo.Get("geometry.type")
	if geotype.String() == usecase.LineString.ToString() {
		return parsegeo.Get("geometry").String(), nil
	} else if geotype.String() == usecase.Polygon.ToString() {
		return parsegeo.Get("geometry").String(), nil
	} else if geotype.String() == usecase.Point.ToString() {
		return parsegeo.Get("geometry").String(), nil
	} else {
		return "", &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  fmt.Sprintf("Geometry Type %s not supported", geotype),
		}
	}
}

// insertGeofenceArea...
// Insert geofence Area to database
func (gu *GeofencingUsecase) insertGeofenceArea(ctx context.Context, request usecase.RequestAddGeofence) error {
	byteDetect, _ := json.Marshal(request.Detect)

	geoVal, err := gu.parseGeojsonGeometry(request)
	if err != nil {
		return err
	}
	jsonRaw := json.RawMessage(geoVal)
	byteJson, _ := json.Marshal(&jsonRaw)

	if err := gu.geofenceManager.InsertGeofenceArea(
		ctx,
		entity.GeofenceArea{
			LocationId:   request.LocationId,
			Name:         request.Name,
			LocationType: request.LocationType,
			Detect:       byteDetect,
			Geojson:      byteJson,
			GeofenceType: request.GeofenceType,
			ChannelName:  request.ChannelName,
		},
	); err != nil {
		util.Logger().Error(err)
		return &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	// update location to geofence
	if err := gu.geofenceManager.UpdateLocationToGeofence(ctx, request.LocationId); err != nil {
		util.Logger().Error(err)
		return &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return nil
}

// Create Channel to tile38...
func (gu *GeofencingUsecase) createGeofenceChannel(ctx context.Context, request usecase.RequestAddGeofence) error {
	var feature *paulgeojson.Feature
	if err := json.Unmarshal([]byte(request.Geojson), &feature); err != nil {
		return err
	}
	if err := gu.tile38Chan.SetGeofenceChannel(
		tile38.Geofence{
			Name:    request.ChannelName,
			Key:     getenv.GetString("GEOFENCING_KEY", "geofencing"),
			Detect:  request.Detect,
			Action:  repository.Within.ToString(),
			Feature: feature,
		},
	); err != nil {
		util.Logger().Error(err)
		return &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return nil
}
