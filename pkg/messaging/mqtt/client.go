package mqtt

import (
	"errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"time"
)

const (
	qos = 2
)

var errConnect = errors.New("failed to connect to MQTT broker")

type Client struct {
	mqtt.Client
	timeout time.Duration
	logger  logger.LoggingClient
}

func NewClient(address string, username string, timeout time.Duration, logger logger.LoggingClient) (*Client, error) {
	opts := mqtt.NewClientOptions().
		SetUsername(username).
		AddBroker(address)
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

	c := &Client{
		Client:  client,
		timeout: timeout,
		logger:  logger,
	}
	return c, nil
}

func (c Client) Close() error {
	c.Disconnect(10 * 1000) // 10秒
	return nil
}
