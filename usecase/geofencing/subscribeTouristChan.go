package geofencing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
	"github.com/xjem/t38c"
)

func (gu *GeofencingUsecase) SubscribeTouristChan(ge *t38c.GeofenceEvent) {

	switch ge.Detect {
	case repository.Inside.ToString():
		if err := gu.processInsideDetect(ge); err != nil {
			return
		}
	case repository.Enter.ToString():
		if err := gu.processEnterDetect(ge); err != nil {
			return
		}
	default:
		util.Logger().Info("Nothing detect")
	}
}

// setLastDetect
func (gu *GeofencingUsecase) setLastDetect(ge *t38c.GeofenceEvent) error {
	objId := fmt.Sprintf("last:%s", ge.ID)
	speedField, ok := ge.Fields["speed"]
	if !ok {
		return errors.New("fields not found")
	}
	if err := gu.tile38Manager.SetGeofencingKey(tile38.SetKey{
		Key:      getenv.GetString("GEOFENCING_LAST_DETECT_KEY", "lastgeofencing"),
		ObjectId: objId,
		Lat:      ge.Object.Geometry.Point[1],
		Long:     ge.Object.Geometry.Point[0],
		Fields: tile38.Field{
			Timestamp: float64(time.Now().Unix()),
			Speed:     speedField,
		},
	}); err != nil {
		util.Logger().Error(err)
		return err
	}
	return nil
}

// evaluateDetectTime...
func (gu *GeofencingUsecase) evaluateDetectTime(ge *t38c.GeofenceEvent) (bool, error) {
	// get last detect time
	objId := fmt.Sprintf("last:%s", ge.ID)
	last, err := gu.tile38Manager.GetLastGeofencingDetect(
		getenv.GetString("GEOFENCING_LAST_DETECT_KEY", "lastgeofencing"),
		objId,
		true,
	)
	if err != nil {
		util.Logger().Error(err)
		// if last key not available
		if isAvailable := strings.Contains(err.Error(), "key not found"); isAvailable {
			if err := gu.setLastDetect(ge); err != nil {
				return false, err
			}
			return false, nil
		}
		return false, err
	}

	timestampField, ok := last.Fields["timestamp"]
	if !ok {
		return false, errors.New("Field timestamp not found")
	}

	// convert time to timestamp
	trackTime := time.Unix(int64(timestampField), 0)

	// get diff time from track timestamp with detected timestamp
	diffTime := ge.Time.Sub(trackTime)

	// parse duration
	parseDuration, _ := time.ParseDuration(diffTime.String())
	detectDuration := getenv.GetInt("MOBILITY_DETECT_DURATION", 30)
	if int64(parseDuration.Seconds()) >= int64(detectDuration) {
		// set last detect
		if err := gu.setLastDetect(ge); err != nil {
			return false, err
		}
		return true, nil
	} else {
		if err := gu.setLastDetect(ge); err != nil {
			return false, err
		}
		return false, nil
	}

}

// getLastMobilityCountByDetect
func (gu *GeofencingUsecase) getLastMobilityCountByDetect(ctx context.Context, ge *t38c.GeofenceEvent) (int, error) {
	// get geofence name
	detail, err := gu.geofenceManager.DetailGeofenceAreaByChannelName(ctx, ge.Hook)
	if err != nil {
		util.Logger().Error(err)
		return 0, err
	}
	count, err := gu.mobilityManager.GetLastAggregateFieldValue(
		ctx,
		ge.Detect,
		detail.Id,
	)
	if err != nil {
		util.Logger().Error(err)
		if err.Error() == "Field is empty" {
			// insert default value
			if err := gu.mobilityManager.InsertDefaultValueAggregateField(
				ctx,
				ge.Detect,
				detail.Id,
				1,
			); err != nil {
				util.Logger().Error(err)
			}
			return 0, nil
		}
		return 0, err
	}
	return count, nil

}

// updateMobilityCountByDetect
func (gu *GeofencingUsecase) updateMobilityCountByDetect(ctx context.Context, ge *t38c.GeofenceEvent, value int) error {
	// get geofence name
	detail, err := gu.geofenceManager.DetailGeofenceAreaByChannelName(ctx, ge.Hook)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	if err := gu.mobilityManager.UpdateAggregateFieldValue(
		ctx,
		ge.Detect,
		detail.Id,
		value,
	); err != nil {
		util.Logger().Error(err)
		return err
	}

	return nil
}

// insertMobilityCounter
func (gu *GeofencingUsecase) insertMobilityCounter(ctx context.Context, ge *t38c.GeofenceEvent) error {
	count, err := gu.getLastMobilityCountByDetect(ctx, ge)
	if err != nil {
		return err
	}
	// increment count
	increment := count + 1
	if err := gu.updateMobilityCountByDetect(ctx, ge, increment); err != nil {
		return err
	}
	return nil
}

// processInsideDetect...
func (gu *GeofencingUsecase) processInsideDetect(ge *t38c.GeofenceEvent) error {
	ctx := context.Background()
	log.Println(fmt.Sprintf("Name: %s Detect: %s ID: %s Long: %f Lat: %f",
		ge.Hook,
		ge.Detect,
		ge.ID,
		ge.Object.Geometry.Point[0], // long
		ge.Object.Geometry.Point[1], // lat
	))

	detected, err := gu.evaluateDetectTime(ge)
	if err != nil {
		return err
	}
	if !detected {
		go gu.NotifyDetectTourist(ctx, ge)
		return nil
	}

	if err := gu.insertMobilityCounter(ctx, ge); err != nil {
		return err
	}
	go gu.NotifyDetectTourist(ctx, ge)

	return nil

}

// processEnterDetect
func (gu *GeofencingUsecase) processEnterDetect(ge *t38c.GeofenceEvent) error {
	ctx := context.Background()
	log.Println(fmt.Sprintf("Name: %s Detect: %s ID: %s Long: %f Lat: %f",
		ge.Hook,
		ge.Detect,
		ge.ID,
		ge.Object.Geometry.Point[0], // long
		ge.Object.Geometry.Point[1], // lat
	))

	detected, err := gu.evaluateDetectTime(ge)
	if err != nil {
		return err
	}
	if !detected {
		go gu.NotifyDetectTourist(ctx, ge)
		return nil
	}
	if err := gu.insertMobilityCounter(ctx, ge); err != nil {
		return err
	}
	go gu.NotifyDetectTourist(ctx, ge)

	return nil
}
