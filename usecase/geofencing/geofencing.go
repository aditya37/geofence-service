package geofencing

import (
	"errors"

	"github.com/aditya37/geofence-service/client"
	"github.com/aditya37/geofence-service/repository"
)

type GeofencingUsecase struct {
	tile38Chan      repository.Tile38ChannelManager
	gcppubsub       repository.Pubsub
	tile38Manager   repository.Tile38ReaderWriter
	mobilityManager repository.MobilityManager
	geofenceManager repository.GeofenceManager
	cache           repository.CacheManager
	geospatialSvc   client.GeospatialServiceClient
}

func NewGeofencingUsecase(
	tile38Chan repository.Tile38ChannelManager,
	gcppubsub repository.Pubsub,
	tile38Manager repository.Tile38ReaderWriter,
	mobilityManager repository.MobilityManager,
	geofenceManager repository.GeofenceManager,
	cache repository.CacheManager,
	geospatialSvc client.GeospatialServiceClient,
) (*GeofencingUsecase, error) {
	if tile38Chan == nil && gcppubsub == nil {
		return nil, errors.New("Please set dependencies")
	}
	return &GeofencingUsecase{
		tile38Chan:      tile38Chan,
		gcppubsub:       gcppubsub,
		tile38Manager:   tile38Manager,
		mobilityManager: mobilityManager,
		geofenceManager: geofenceManager,
		cache:           cache,
		geospatialSvc:   geospatialSvc,
	}, nil
}
