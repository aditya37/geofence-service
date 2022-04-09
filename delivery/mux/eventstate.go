package mux

import "github.com/aditya37/geofence-service/usecase/eventstate"

type EventstateDelivery struct {
	eventStateCase *eventstate.EventstateUsecase
}

func NewEventStateDelivery(eventstateCase *eventstate.EventstateUsecase) *EventstateDelivery {
	return &EventstateDelivery{
		eventStateCase: eventstateCase,
	}
}
