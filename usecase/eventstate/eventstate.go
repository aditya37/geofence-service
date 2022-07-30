package eventstate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/entity/tile38"
	"github.com/aditya37/geofence-service/repository"
	"github.com/aditya37/geofence-service/usecase"
	logger "github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
)

type EventstateUsecase struct {
	pubsub               repository.Pubsub
	eventManager         repository.CacheEventManager
	geofenceManager      repository.GeofenceManager
	tile38ChannelManager repository.Tile38ChannelManager
}

// Deprecated: Not used....
func NewEventStateUsecase(
	pubsub repository.Pubsub,
	eventManager repository.CacheEventManager,
	geofenceManager repository.GeofenceManager,
	tile38ChannelManager repository.Tile38ChannelManager,
) (*EventstateUsecase, error) {
	return &EventstateUsecase{
		pubsub:               pubsub,
		eventManager:         eventManager,
		geofenceManager:      geofenceManager,
		tile38ChannelManager: tile38ChannelManager,
	}, nil
}

// prefix
const (
	msgPrefixSuccess      = "Success %s geofence area"
	msgPrefixInsertFailed = "Failed %s geofence area"
)

// ConsumeEventState...
func (es *EventstateUsecase) ConsumeEventState(ctx context.Context, topicname, servicename string) error {
	result := es.pubsub.Subscribe(
		ctx,
		topicname,
		servicename,
		es.evaluateMessage,
	)
	return result
}

// evaluateMessage..
func (es *EventstateUsecase) evaluateMessage(ctx context.Context, msg repository.PubsubMessage) {
	// load data from json string to struct
	payload, err := es.parseMessage(ctx, msg)
	if err != nil {
		return
	}

	if validService := es.checkServiceSource(payload); !validService {
		logger.Logger().Error("Failed process event,source service must different with destination")
		return
	} else {
		switch payload.State {
		case usecase.EventStateInsert:
			if err := es.insertGeofenceArea(ctx, payload); err != nil {
				// send notify error
				if err := es.eventNotifer(
					ctx,
					payload,
					usecase.EventNotiferTypeFailed,
				); err != nil {
					return
				}
				return
			}

			if err := es.eventNotifer(ctx, payload, usecase.EventNotiferTypeSuccess); err != nil {
				return
			}

			msg.Ack()

		default:
			log.Println(fmt.Sprintf("State not valid State: %s", payload.State))
			msg.Ack()
			return
		}
	}
}

// notifer...
/*
method for publish rollback event if failed/error
or if event success will save log to redis
*/
func (es *EventstateUsecase) eventNotifer(ctx context.Context, data *usecase.GeofenceEventState, notifType string) error {
	if notifType == "SUCCESS" {
		logger.Logger().Info(fmt.Sprintf("Notify %s", usecase.EventNotiferTypeSuccess))

		// set to redis
		data.Metadata.Message = fmt.Sprintf(msgPrefixSuccess, data.State)
		if err := es.storeEventCache(data); err != nil {
			return err
		}

		// convert message to byte
		message := usecase.GeofenceEventState{
			ServiceName: data.ServiceName,
			State:       usecase.EventStateInsertSuccess,
			EventId:     data.EventId,
			Metadata:    data.Metadata,
			GeofenceData: usecase.GeofenceData{
				Name: data.GeofenceData.Name,
			},
		}
		byteMessage, _ := json.Marshal(message)

		// notify to source service if success
		if err := es.pubsub.Publish(
			ctx,
			entity.PublishParam{
				TopicName: getenv.GetString("GEOFENCE_EVENT_STATE_TOPIC", "geofence-event-state"),
				Message:   byteMessage,
			},
		); err != nil {
			logger.Logger().Error(err)
			return err
		}

		return nil
	} else {
		// store event state to redis
		data.Metadata.Message = fmt.Sprintf(msgPrefixInsertFailed, data.State)
		if err := es.storeEventCache(data); err != nil {
			logger.Logger().Error(err)
			return err
		}

		// notify rollback
		// convert message to byte
		message := usecase.GeofenceEventState{
			ServiceName: data.ServiceName,
			State:       usecase.EventStateInsertRollback,
			EventId:     data.EventId,
			Metadata:    data.Metadata,
			GeofenceData: usecase.GeofenceData{
				Name: data.GeofenceData.Name,
			},
		}
		byteMessage, _ := json.Marshal(message)
		if err := es.pubsub.Publish(
			ctx,
			entity.PublishParam{
				TopicName: getenv.GetString("GEOFENCE_EVENT_STATE_TOPIC", "geofence-event-state"),
				Message:   byteMessage,
			},
		); err != nil {
			logger.Logger().Error(err)
			return err
		}
		return nil
	}
}

// storeEventCache...
func (es *EventstateUsecase) storeEventCache(data *usecase.GeofenceEventState) error {
	jsonEventMetaData, _ := json.Marshal(data.Metadata)
	key := fmt.Sprintf("%s:%s", data.ServiceName, data.EventId)

	// set to redis
	if err := es.eventManager.SetEventStateResponse(
		getenv.GetInt("EVENT_STATE_CACHE_DB", 3),
		key,
		jsonEventMetaData,
		time.Duration(
			getenv.GetInt("EVENT_STATE_CACHE_DB", 300)*int(time.Minute),
		),
	); err != nil {
		logger.Logger().Error(err)
		return err
	}
	return nil
}

// insertGeofenceArea...
func (es *EventstateUsecase) insertGeofenceArea(ctx context.Context, data *usecase.GeofenceEventState) error {
	// TODO: Insert to database
	byteDetect, _ := json.Marshal(data.GeofenceData.Detect)
	byteGeojson, _ := data.GeofenceData.Shape.MarshalJSON()
	if err := es.geofenceManager.InsertGeofenceArea(
		ctx,
		entity.GeofenceArea{
			GeofenceId: data.GeofenceData.Id,
			Name:       data.GeofenceData.Name,
			AreaType:   data.GeofenceData.AreaType,
			Detect:     byteDetect,
			Geojson:    byteGeojson,
		},
	); err != nil {
		logger.Logger().Error(err)
		return err
	}
	// TODO: Create Channel in tile38
	if err := es.tile38ChannelManager.SetGeofenceChannel(tile38.Geofence{
		Name:    data.GeofenceData.Name,
		Key:     getenv.GetString("GEOFENCING_KEY", "geofencing"),
		Detect:  data.GeofenceData.Detect,
		Action:  repository.Within.ToString(),
		Feature: &data.GeofenceData.Shape,
	}); err != nil {
		logger.Logger().Error(err)
		return err
	}

	return nil
}

// checkServiceSource...
/*
validate service source, if source message
same with consumer, data not be process
TODO: register check service name from database
*/
func (es *EventstateUsecase) checkServiceSource(data *usecase.GeofenceEventState) bool {
	if data.ServiceName != getenv.GetString("SERVICE_NAME", "geofence-service") {
		return true
	}
	return false
}

// parseMessage...
func (es *EventstateUsecase) parseMessage(ctx context.Context, msg repository.PubsubMessage) (*usecase.GeofenceEventState, error) {
	// TODO: This method for parse or marshal message
	data := msg.GetMessage()
	if data == nil {
		logger.Logger().Error(usecase.ErrMessageDataNil)
		return nil, usecase.ErrMessageDataNil
	}

	var payload usecase.GeofenceEventState
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Logger().Error(err)
		return nil, err
	}

	return &payload, nil
}
