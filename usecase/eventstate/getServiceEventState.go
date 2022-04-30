package eventstate

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/util"
	getenv "github.com/aditya37/get-env"
)

func (es *EventstateUsecase) GetServiceEventState(ctx context.Context, request usecase.GetServiceEventStateRequest) (usecase.GetServiceEventStateResponse, error) {

	key, err := es.eventStateKey(request)
	if err != nil {
		return usecase.GetServiceEventStateResponse{}, err
	}
	resp, err := es.eventManager.GetEventState(getenv.GetInt("EVENT_STATE_CACHE_DB", 3), key)
	if err != nil {
		return usecase.GetServiceEventStateResponse{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}

	var response usecase.GetServiceEventStateResponse
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return usecase.GetServiceEventStateResponse{}, &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  err.Error(),
		}
	}
	return response, nil
}

func (es *EventstateUsecase) eventStateKey(request usecase.GetServiceEventStateRequest) (string, error) {
	if request.EventId == "" || request.ServiceName == "" {
		return "", &util.ErrorMsg{
			HttpRespCode: http.StatusBadRequest,
			Description:  "Please set all request",
		}
	}
	return fmt.Sprintf("%s:%s", request.ServiceName, request.EventId), nil
}
