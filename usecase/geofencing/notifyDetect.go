package geofencing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	geospatialSrv "github.com/aditya37/api-contract/geospatial-service/service"
	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
	"github.com/xjem/t38c"
)

const (
	dailyAvgInterval             = 7
	errorFieldPropertiesNotfound = "Field in properties not found"
	cacheLocationDetailPrfx      = "cache:location_detail:%d"
	cacheNearbyLocationPrfx      = "nearby-location:%f:%f:%d:%d"
)

func (gu *GeofencingUsecase) NotifyDetectTourist(ctx context.Context, ge *t38c.GeofenceEvent) error {
	// get daily avg
	dailyAvg, err := gu.mobilityManager.GetDailyAverage(ctx, dailyAvgInterval)
	if err != nil {
		util.Logger().Error(err)
		return err
	}

	// get all area daily avg
	dailyAreaAvg, err := gu.mobilityManager.GetAllAreaDailyAverage(ctx, dailyAvgInterval)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	timeLoc, _ := time.LoadLocation("Asia/Jakarta")

	// response converter
	dailyAreaAvgResp, err := gu.responseConverterDailyAreaAvg(ctx, dailyAreaAvg)
	if err != nil {
		return err
	}

	getNearbyLocation, err := gu.getNearbyLocation(
		ctx,
		ge.Object.Geometry.Point[1],
		ge.Object.Geometry.Point[0],
		int64(getenv.GetInt("RADIUS_NEARBY_LOCATION", 3000)),
		int64(getenv.GetInt("COUNT_NEARBY_LOCATION", 10)),
	)
	if err != nil {
		return err
	}

	// get current Position of object...
	object, _ := json.Marshal(ge.Object.Geometry)
	deviceId := strings.Replace(ge.ID, "geofencing:", "", -1)
	// assert data to struct
	msgPayload := usecase.NotifyGeofencingPayload{
		Type:        "tourist",
		Detect:      ge.Detect,
		ChannelName: ge.Hook,
		DeviceId:    deviceId,
		Mobility: usecase.Mobility{
			DailyAverage: usecase.Detect{
				Enter:  dailyAvg.Enter,
				Inside: dailyAvg.Inside,
				Exit:   dailyAvg.Exit,
				Date:   time.Now().In(timeLoc).Format(time.RFC3339),
			},
			MobilityAreas: dailyAreaAvgResp,
		},
		Object:         string(object),
		NearbyLocation: getNearbyLocation,
	}

	msgByte, err := json.Marshal(&msgPayload)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	// publish...
	if err := gu.publishNotify(
		ctx,
		msgByte,
		getenv.GetString("GEOFENCING_TOPIC_DETECT", "geofencing-detect"),
	); err != nil {
		return err
	}
	return nil
}

// getNearbyLocation....
func (gu *GeofencingUsecase) getNearbyLocation(ctx context.Context, lat, long float64, radius, count int64) (geospatialSrv.GetNearbyLocationResponse, error) {
	key := fmt.Sprintf(cacheNearbyLocationPrfx, lat, long, radius, count)
	resp, err := gu.cache.Get(key)
	if err != nil {
		nearbyLoc, err := gu.geospatialSvc.GetNearbyLocation(
			ctx,
			geospatialSrv.GetNearbyLocationRequest{
				Count:        count,
				ShowDistance: true,
				Filter: geospatialSrv.NearbyFilter{
					WithType: false,
				},
				CurrentPosition: geospatialSrv.Position{
					Lat:    lat,
					Long:   long,
					Radius: radius,
				},
			})
		if err != nil {
			util.Logger().Error(err)
			return geospatialSrv.GetNearbyLocationResponse{}, err
		}
		response := geospatialSrv.GetNearbyLocationResponse{
			Count:          nearbyLoc.Count,
			NearbyLocation: nearbyLoc.NearbyLocation,
		}
		jsonByte, err := json.Marshal(&response)
		if err != nil {
			return geospatialSrv.GetNearbyLocationResponse{}, err
		}
		if err := gu.cache.Set(
			key,
			time.Duration(
				getenv.GetInt("CACHE_TTL", 100)*int(time.Minute),
			),
			jsonByte,
		); err != nil {
			util.Logger().Error(err)
			return geospatialSrv.GetNearbyLocationResponse{}, err
		}
		return response, nil
	}

	var respPayload geospatialSrv.GetNearbyLocationResponse
	if err := json.Unmarshal([]byte(resp), &respPayload); err != nil {
		util.Logger().Error(err)
		return geospatialSrv.GetNearbyLocationResponse{}, err
	}

	return respPayload, nil
}

// responseConverterDailyAreaAvg...
func (gu *GeofencingUsecase) responseConverterDailyAreaAvg(ctx context.Context, data []*entity.ResultGetDailyAvg) ([]usecase.MobilityArea, error) {
	var result []usecase.MobilityArea
	timeLoc, _ := time.LoadLocation("Asia/Jakarta")
	for _, val := range data {

		locDetail, err := gu.getLocationDetailByGeofenceId(ctx, val.GeofencId)
		if err != nil {
			return result, err
		}
		result = append(result, usecase.MobilityArea{
			GeofenceId:   val.GeofencId,
			LocationId:   locDetail.LocationId,
			LocationName: locDetail.LocatioName,
			Average: usecase.Detect{
				Enter:  val.Enter,
				Inside: val.Inside,
				Exit:   val.Exit,
				Date:   time.Now().In(timeLoc).Format(time.RFC3339),
			},
		})

	}
	return result, nil

}

// getLocationDetailByid
func (gu *GeofencingUsecase) getLocationDetailByGeofenceId(ctx context.Context, id int64) (*entity.ResultGetLocationDetailByGeofenceId, error) {
	cache, err := gu.cache.Get(fmt.Sprintf(cacheLocationDetailPrfx, id))
	if err != nil {
		geofArea, err := gu.geofenceManager.DetailGeofenceAreaById(ctx, id)
		if err != nil {
			util.Logger().Error(err)
			return nil, err
		}

		// hit to geospatial service
		locationDetail, err := gu.geospatialSvc.GetLocationById(
			ctx,
			geofArea.LocationId,
		)
		if err != nil {
			util.Logger().Error(err)
			return nil, err
		}

		result := entity.ResultGetLocationDetailByGeofenceId{
			LocationId:  geofArea.LocationId,
			LocatioName: locationDetail.LocationName,
		}
		// convert to json
		jsonData, err := json.Marshal(&result)
		if err != nil {
			return nil, err
		}

		if err := gu.cache.Set(
			fmt.Sprintf(cacheLocationDetailPrfx, id),
			time.Duration(
				getenv.GetInt("CACHE_TTL", 100)*int(time.Second),
			),
			jsonData,
		); err != nil {
			util.Logger().Error(err)
			return nil, err
		}

		return &result, nil
	}

	var result entity.ResultGetLocationDetailByGeofenceId
	if err := json.Unmarshal([]byte(cache), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// publish
func (gu *GeofencingUsecase) publishNotify(ctx context.Context, data []byte, topicname string) error {
	if err := gu.gcppubsub.Publish(
		ctx,
		entity.PublishParam{
			TopicName: topicname,
			Message:   data,
		},
	); err != nil {
		util.Logger().Error(err)
		return err
	}
	return nil
}
