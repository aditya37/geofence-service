package eventstate

import (
	"context"

	"github.com/aditya37/geofence-service/usecase"
)

type EventstateUsecase struct {
}

func NewEventStateUsecase() (*EventstateUsecase, error) {
	return &EventstateUsecase{}, nil
}

// getEventById
func (es *EventstateUsecase) GetEventById(ctx context.Context, eventid, sevicename string) (usecase.GetEventStateResponse, error) {
	return usecase.GetEventStateResponse{}, nil
}
