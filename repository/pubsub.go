package repository

import (
	"context"

	"github.com/aditya37/geofence-service/entity"
)

type (
	PubsubMessage interface {
		GetMessage() []byte
		Ack()
	}
	// pubsub manager
	Pubsub interface {
		Subscribe(ctx context.Context, topic, servicename string, Callback interface{}) error
		Publish(ctx context.Context, param entity.PublishParam) error
		Close() error
	}
)
