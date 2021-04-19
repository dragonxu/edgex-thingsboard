package mqtt_server

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"time"
)

var (
	errPublishTimeout     = errors.New("failed to publish due to timeout reached")
	errSubscribeTimeout   = errors.New("failed to subscribe due to timeout reached")
	errUnsubscribeTimeout = errors.New("failed to unsubscribe due to timeout reached")
)

type Handler func(req []byte) (resp []byte, err error)

type Server struct {
	pubsub        *PubSubClient
	subscriptions []string
}

func New(client mqtt.Client, qos byte, timeout time.Duration, logger logger.LoggingClient) *Server {
	pubsub := NewPubSubClient(client, qos, timeout, logger)
	return &Server{
		pubsub: pubsub,
	}
}

func (f Server) HandleFunc(requestTopic, replyTopic string, handler Handler) error {
	err := f.pubsub.Subscribe(requestTopic, func(topic string, msg []byte) error {
		res, err := handler(msg)
		if err != nil {
			return err
		}

		return f.pubsub.Publish(replyTopic, res)
	})
	if err != nil {
		return err
	}

	f.subscriptions = append(f.subscriptions, requestTopic)
	return nil
}

func (f Server) Close() error {
	for _, t := range f.subscriptions {
		if err := f.pubsub.Unsubscribe(t); err != nil {
			return err
		}
	}

	f.subscriptions = nil
	return nil
}
