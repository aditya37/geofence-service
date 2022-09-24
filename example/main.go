package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/repository/pubsub"
	"github.com/aditya37/geofence-service/usecase"
)

func init() {
	os.Setenv("PUBSUB_EMULATOR_HOST", "37.44.244.196:8085")
}
func main() {
	gp, err := pubsub.NewGcpPubsub(context.Background(), "pubsub-emulator")
	if err != nil {
		log.Println(err)
		return
	}

	tracking := []usecase.TrackingPayload{
		{
			Speed: 5,
			Lat:   -7.21626272274221,
			Long:  111.86905592679977,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.216651223242152,
			Long:  111.86900496482848,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.217199380914439,
			Long:  111.86896741390228,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.217561271052578,
			Long:  111.86894059181213,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.2176065072994975,
			Long:  111.86901569366455,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.212965776993304,
			Long:  111.86944484710693,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.212965776993304,
			Long:  111.86944484710693,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.209857738352886,
			Long:  111.86983108520506,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.212284564870234,
			Long:  111.8712043762207,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.216414397634615,
			Long:  111.87274932861327,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
		{
			Speed: 5,
			Lat:   -7.2210125176034845,
			Long:  111.87438011169434,
			Device: usecase.DeviceMetadata{
				DeviceId: "111",
			},
		},
	}

	for _, v := range tracking {
		v.Timestamp = time.Now().Unix()
		j, _ := json.Marshal(v)
		log.Println(
			fmt.Sprintf("Tracking => %s", string(j)),
		)
		if err := gp.Publish(context.Background(), entity.PublishParam{
			TopicName: "tracking-forward-topic",
			Retained:  false,
			Message:   j,
		}); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(v.Speed) * time.Second)
	}
}
