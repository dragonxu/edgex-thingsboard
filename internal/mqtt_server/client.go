package mqtt_server

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var errConnect = errors.New("failed to connect to MQTT broker")

func NewClient(opts *mqtt.ClientOptions, timeout time.Duration) (mqtt.Client, error) {
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Error() != nil {
		return nil, token.Error()
	}

	ok := token.WaitTimeout(timeout)
	if ok && token.Error() != nil {
		return nil, token.Error()
	}
	if !ok {
		return nil, errConnect
	}

	return client, nil
}
