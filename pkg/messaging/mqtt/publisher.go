package mqtt

import (
	"errors"
	"github.com/inspii/edgex-thingsboard/pkg/messaging"
)

var errPublishTimeout = errors.New("failed to publish due to timeout reached")

type publisher struct {
	client *Client
}

func NewPublisher(client *Client) messaging.Publisher {
	return publisher{
		client: client,
	}
}

func (p publisher) Publish(topic string, msg []byte) error {
	token := p.client.Publish(topic, qos, false, msg)
	if token.Error() != nil {
		return token.Error()
	}
	ok := token.WaitTimeout(p.client.timeout)
	if ok && token.Error() != nil {
		return token.Error()
	}
	if !ok {
		return errPublishTimeout
	}
	return nil
}

func (p publisher) Close() error {
	return p.client.Close()
}
