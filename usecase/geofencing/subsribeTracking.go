package geofencing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	env "github.com/aditya37/get-env"
)

func (gu *geofencingUsecase) SubscribeLocationTracking(ctx context.Context, topicname, servicename string) error {
	if err := gu.gcppubsub.Subscribe(
		ctx,
		topicname,
		servicename,
		func(ctx context.Context, msg repository.PubsubMessage) {
			var payload usecase.TrackingPayload
			if err := json.Unmarshal(msg.GetMessage(), &payload); err != nil {
				util.Logger().Error(err)
				return
			}
			// set geofencing key
			objId := fmt.Sprintf("geofencing:%s", payload.Device.DeviceId)
			if err := gu.tile38Manager.SetGeofencingKey(
				tile38.SetKey{
					Key:      env.GetString("GEOFENCE_KEY", "geofencing"),
					ObjectId: objId,
					Lat:      payload.Lat,
					Long:     payload.Long,
					Fields: tile38.Field{
						Timestamp: float64(payload.Timestamp),
						Speed:     float64(payload.Speed),
					},
				},
			); err != nil {
				util.Logger().Error(err)
				return
			}
		},
	); err != nil {
		return err
	}
	return nil
}
