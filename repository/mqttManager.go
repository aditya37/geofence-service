package repository

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MqttManager interface {
	Subscribe(topic string, qos byte, f func(c mqtt.Client, m mqtt.Message)) error
	Publish(topic string, qos byte, retain bool, message []byte) error
}
