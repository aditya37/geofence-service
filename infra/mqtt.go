package infra

import (
	"errors"
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttClientInstance  mqtt.Client = nil
	mqttClientSingleton sync.Once
	err                 error
)

type (
	MQTTConf struct {
		Host     string
		Port     int64
		ClientId string
		Username string
		Password string
	}
)

func NewMqttClientInstance(param MQTTConf) error {
	mqttClientSingleton.Do(func() {
		brokerURL := fmt.Sprintf("tcp://%s:%d", param.Host, param.Port)
		opts := mqtt.NewClientOptions()
		opts.AddBroker(brokerURL)
		opts.SetClientID(param.ClientId)
		opts.SetUsername(param.Username)
		opts.SetPassword(param.Password)

		client := mqtt.NewClient(opts)
		if !client.IsConnected() && !client.Connect().Wait() {
			err = errors.New(fmt.Sprintf("MQTT Broker not connected"))
		}
		mqttClientInstance = client
	})
	if err != nil {
		return err
	}
	return nil
}

func GetMqttClientInstance() mqtt.Client {
	return mqttClientInstance
}
