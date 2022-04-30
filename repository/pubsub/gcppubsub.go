package pubsub

import (
	"context"
	"errors"
	"fmt"
	"time"

	gpubsub "cloud.google.com/go/pubsub"
	"github.com/aditya37/geofence-service/entity"
	"github.com/aditya37/geofence-service/repository"
	"github.com/aditya37/geofence-service/util"

	getenv "github.com/aditya37/get-env"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type gcppubsub struct {
	client *gpubsub.Client
}

func NewGcpPubsub(ctx context.Context, projectid string, opts ...option.ClientOption) (repository.Pubsub, error) {
	conn, err := gpubsub.NewClient(ctx, projectid, opts...)
	if err != nil {
		return nil, err
	}

	return &gcppubsub{
		client: conn,
	}, nil

}

// topic..
func (gp *gcppubsub) createTopic(ctx context.Context, topic string) error {
	if _, err := gp.client.CreateTopic(ctx, topic); err != nil {
		return err
	}
	return nil
}

// get topic...
func (gp *gcppubsub) getTopic(ctx context.Context, topicname string) (*gpubsub.Topic, error) {
	topic := gp.client.Topic(topicname)
	ok, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		if err := gp.createTopic(ctx, topicname); err != nil {
			return nil, err
		}
		return nil, errors.New("topic not found")
	}
	return topic, nil
}

// message callback...
type pubsubMessage struct {
	msg *gpubsub.Message
}

func (pm pubsubMessage) GetMessage() []byte {
	return pm.msg.Data
}
func (pm pubsubMessage) Ack() {
	pm.msg.Ack()
}

// createSubscription...
func (gp *gcppubsub) createSubscription(ctx context.Context, servicename, topicname string) (*gpubsub.Subscription, error) {
	topic, err := gp.getTopic(ctx, topicname)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err == nil {
		servicename = servicename + "." + id.String()
	}

	return gp.client.CreateSubscription(
		ctx,
		servicename,
		gpubsub.SubscriptionConfig{
			Topic:               topic,
			RetainAckedMessages: getenv.GetBool("PUBSUB_RETAIN_ACKEDMSG", false),
			RetentionDuration: time.Duration(
				getenv.GetInt("PUBSUB_RETENTION_DURATION", 15) * int(time.Minute),
			),
		},
	)
}

// Subscribe...
func (gp *gcppubsub) Subscribe(ctx context.Context, topic, servicename string, Callback interface{}) error {
	subs, err := gp.createSubscription(ctx, servicename, topic)
	if err != nil {
		return err
	}
	if err := subs.Receive(
		ctx,
		func(ctx context.Context, m *gpubsub.Message) {
			fn := Callback.(func(context.Context, repository.PubsubMessage))
			fn(ctx, pubsubMessage{
				msg: m,
			})
		},
	); err != nil {
		return err
	}
	return nil
}

// Publish...
func (gp *gcppubsub) Publish(ctx context.Context, param entity.PublishParam) error {
	topic, err := gp.getTopic(ctx, param.TopicName)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	resp := topic.Publish(
		ctx,
		&gpubsub.Message{
			Data: param.Message,
		},
	)
	msgId, err := resp.Get(ctx)
	if err != nil {
		util.Logger().Error(err)
		return err
	}
	util.Logger().Info(fmt.Sprintf("Message published with id %s", msgId))
	return nil
}

// closee..
func (gp *gcppubsub) Close() error {
	return gp.client.Close()
}
