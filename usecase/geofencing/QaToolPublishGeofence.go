package geofencing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
)

func (gu *GeofencingUsecase) QAToolPublishGeofence(ctx context.Context, request usecase.MQTTRespTracking) (usecase.ResponseQAToolPublishGeofence, error) {
	data, _ := json.Marshal(request)
	util.Logger().Info("Send Data to =>", string(data))
	if err := gu.gcppubsub.Publish(ctx, entity.PublishParam{
		TopicName: "tracking-forward-topic",
		Retained:  false,
		Message:   data,
	}); err != nil {
		return usecase.ResponseQAToolPublishGeofence{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return usecase.ResponseQAToolPublishGeofence{
		Message: "Success Publish",
	}, nil
}
