package mqtt_server

import (
	"errors"
	"fmt"
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
	client        mqtt.Client
	qos           byte
	timeout       time.Duration
	logger        logger.LoggingClient
	subscriptions []string
}

func New(client mqtt.Client, qos byte, timeout time.Duration, logger logger.LoggingClient) *Server {
	return &Server{
		client:  client,
		qos:     qos,
		timeout: timeout,
		logger:  logger,
	}
}

func (f Server) HandleFunc(requestTopic, replyTopic string, handler Handler) error {
	err := f.subscribe(requestTopic, func(topic string, msg []byte) error {
		res, err := handler(msg)
		if err != nil {
			return err
		}

		return f.publish(replyTopic, res)
	})
	if err != nil {
		return err
	}

	f.subscriptions = append(f.subscriptions, requestTopic)
	return nil
}

func (f Server) Close() error {
	for _, t := range f.subscriptions {
		if err := f.unsubscribe(t); err != nil {
			return err
		}
	}

	f.subscriptions = nil
	return nil
}

func (f *Server) publish(topic string, msg []byte) error {
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

func (f *Server) subscribe(topic string, handler func(topic string, msg []byte) error) error {
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

func (f *Server) unsubscribe(topic string) error {
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
