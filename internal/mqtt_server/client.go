package mqtt_server

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
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

type PubSubClient struct {
	client  mqtt.Client
	qos     byte
	timeout time.Duration
	logger  logger.LoggingClient
}

func NewPubSubClient(client mqtt.Client, qos byte, timeout time.Duration, logger logger.LoggingClient) *PubSubClient {
	return &PubSubClient{
		client:  client,
		qos:     qos,
		timeout: timeout,
		logger:  logger,
	}
}

func (f *PubSubClient) Publish(topic string, msg []byte) error {
	token := f.client.Publish(topic, f.qos, false, msg)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(f.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errPublishTimeout
	}
	return nil
}

func (f *PubSubClient) Subscribe(topic string, handler func(topic string, msg []byte) error) error {
	token := f.client.Subscribe(topic, f.qos, func(client mqtt.Client, m mqtt.Message) {
		if err := handler(m.Topic(), m.Payload()); err != nil {
			f.logger.Warn(fmt.Sprintf("failed to handle message: err=%s, payload=%s", err, m.Payload()))
		}
	})
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(f.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errSubscribeTimeout
	}
	return nil
}

func (f *PubSubClient) Unsubscribe(topic string) error {
	token := f.client.Unsubscribe(topic)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(f.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errUnsubscribeTimeout
	}
	return nil
}
