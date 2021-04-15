package mqtt

import (
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/inspii/edgex-thingsboard/pkg/messaging"
)

var (
	errSubscribeTimeout   = errors.New("failed to subscribe due to timeout reached")
	errUnsubscribeTimeout = errors.New("failed to unsubscribe due to timeout reached")
)

type subscriber struct {
	client *Client
}

func NewSubscriber(client *Client) messaging.Subscriber {
	return &subscriber{
		client: client,
	}
}

func (s subscriber) Subscribe(topic string, handler messaging.MessageHandler) error {
	token := s.client.Subscribe(topic, qos, s.mqttHandler(handler))
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(s.client.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errSubscribeTimeout
	}
	return nil
}

func (s subscriber) Unsubscribe(topic string) error {
	token := s.client.Unsubscribe(topic)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(s.client.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errUnsubscribeTimeout
	}
	return nil
}

func (s subscriber) Close() error {
	return s.client.Close()
}

func (s subscriber) mqttHandler(h messaging.MessageHandler) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		if err := h(m.Topic(), m.Payload()); err != nil {
			s.client.logger.Warn(fmt.Sprintf("failed to handle message: err=%s, payload=%s", err, m.Payload()))
		}
	}
}
