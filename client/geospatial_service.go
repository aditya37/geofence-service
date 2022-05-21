package client

import (
	"context"
	"fmt"

	geospatialSrv "github.com/aditya37/api-contract/geospatial-service/service"
)

// path
type UrlPath string

const (
	PathServiceHealthCheck UrlPath = "/"
	PathGetLocationById    UrlPath = "/locations/%d"
	PathGetNearbyLocation  UrlPath = "/locations/nearby"
)

func (up UrlPath) ToString() string {
	return fmt.Sprintf("%s", up)
}

type GeospatialServiceClient interface {
	GetLocationById(ctx context.Context, locationid int64) (*ResponseGetLocationById, error)
	GetNearbyLocation(ctx context.Context, request geospatialSrv.GetNearbyLocationRequest) (*geospatialSrv.GetNearbyLocationResponse, error)
}
