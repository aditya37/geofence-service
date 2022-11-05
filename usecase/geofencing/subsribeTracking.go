package geofencing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	env "github.com/aditya37/get-env"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (gu *GeofencingUsecase) SubscribeDeviceTracking(c mqtt.Client, m mqtt.Message) {

	var payload usecase.MQTTRespTracking
	if err := json.Unmarshal(m.Payload(), &payload); err != nil {
		util.Logger().Error(err)
		return
	}
	// mapping device type if value 0
	if payload.DeviceType == 0 {
		payload.DeviceType = usecase.DeviceTypeZero
	}

	objId := fmt.Sprintf("geofencing:%s", payload.DeviceId)
	if err := gu.tile38Manager.SetGeofencingKey(
		tile38.SetKey{
			Key:      env.GetString("GEOFENCE_KEY", "geofencing"),
			ObjectId: objId,
			Lat:      payload.GPSData.Lat,
			Long:     payload.GPSData.Long,
			Fields: tile38.Field{
				Timestamp:  float64(time.Now().Unix()),
				Speed:      float64(payload.GPSData.Speed),
				DeviceType: float64(payload.DeviceType),
				Id:         float64(payload.Id),
			},
		},
	); err != nil {
		util.Logger().Error(err)
		return
	}

}
