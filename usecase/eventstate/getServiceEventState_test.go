package eventstate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aditya37/geofence-service/mock"
	"github.com/aditya37/geofence-service/usecase"
	"github.com/aditya37/geofence-service/usecase/eventstate"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceEventState(t *testing.T) {
	tests := []struct {
		name         string
		requst       []usecase.GetServiceEventStateRequest
		expectedResp usecase.GetServiceEventStateResponse
		err          error
	}{
		{
			name: "get state with empty service name",
			requst: []usecase.GetServiceEventStateRequest{
				{
					ServiceName: "",
					EventId:     "1234",
				},
			},
			expectedResp: usecase.GetServiceEventStateResponse{},
			err:          errors.New("Please set all request"),
		},
		{
			name: "Get state with empty event id",
			requst: []usecase.GetServiceEventStateRequest{
				{
					ServiceName: "Test-service",
					EventId:     "",
				},
			},
			expectedResp: usecase.GetServiceEventStateResponse{},
			err:          errors.New("Please set all request"),
		}, {
			name: "Get event state from service geospatial",
			requst: []usecase.GetServiceEventStateRequest{
				{
					ServiceName: "Geospatial-Service",
					EventId:     "2022",
				},
			},
			expectedResp: usecase.GetServiceEventStateResponse{
				Message:      "Success Add geofence area",
				PublishedAt:  "2022-01-11",
				LocationId:   1,
				LocationName: "Nothing",
				LocationType: "Road",
			},
			err: nil,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			var err error
			ctx := context.Background()

			//dependency
			mockEventCache, _ := mock.NewMockCacheEventManager()

			// initial insert event
			mockEventCache.SetEventStateResponse(
				3,
				"Geospatial-Service:2022",
				[]byte(mock.ValidEventJson),
				300,
			)
			// usecase
			eventstateCase, _ := eventstate.NewEventStateUsecase(nil, mockEventCache)
			actualResp := usecase.GetServiceEventStateResponse{}
			defer mockEventCache.Close()
			// DO Test..
			for _, req := range tt.requst {
				actualResp, err = eventstateCase.GetServiceEventState(ctx, req)
				if err != nil {
					assert.Error(t, tt.err, err)
				}
			}
			assert.Equal(t, tt.expectedResp, actualResp)

		})
	}

}
