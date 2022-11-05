package mqttmanager

import (
	"github.com/aditya37/geofence-service/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttManager struct {
	client mqtt.Client
}

func NewTrackingManager(client mqtt.Client) repository.MqttManager {
	return &mqttManager{
		client: client,
	}
}

func (mm *mqttManager) Subscribe(topic string, qos byte, f func(c mqtt.Client, m mqtt.Message)) error {
	if err := mm.client.Subscribe(topic, qos, f).Error(); err != nil {
		return err
	}
	return nil
}

func (nm *mqttManager) Publish(topic string, qos byte, retain bool, message []byte) error {
	if err := nm.client.Publish(topic, qos, retain, message).Error(); err != nil {
		return err
	}
	return nil
}
