package mux

import (
	"net/http"

	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/usecase/eventstate"
	"github.com/gorilla/mux"
)

type EventstateDelivery struct {
	eventStateCase *eventstate.EventstateUsecase
}

func NewEventStateDelivery(eventstateCase *eventstate.EventstateUsecase) *EventstateDelivery {
	return &EventstateDelivery{
		eventStateCase: eventstateCase,
	}
}

// GetServiceEventState...
func (ed *EventstateDelivery) GetServiceEventState(w http.ResponseWriter, r *http.Request) {
	parseParam := mux.Vars(r)

	svcParam := parseParam["service_name"]
	eventIdParam := parseParam["event_id"]

	resp, err := ed.eventStateCase.GetServiceEventState(
		r.Context(),
		usecase.GetServiceEventStateRequest{
			ServiceName: svcParam,
			EventId:     eventIdParam,
		},
	)
	if err != nil {
		EncodeErrorResponse(w, err)
		return
	}
	EncodeResponse(w, http.StatusOK, resp)
	return
}
