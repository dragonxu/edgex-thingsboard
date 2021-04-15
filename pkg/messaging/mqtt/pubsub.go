package mqtt

import "github.com/inspii/edgex-thingsboard/pkg/messaging"

type pubsub struct {
	client *Client
	messaging.Publisher
	messaging.Subscriber
}

func NewPubSub(client *Client) messaging.PubSub {
	pub := NewPublisher(client)
	sub := NewSubscriber(client)
	return &pubsub{
		client:     client,
		Publisher:  pub,
		Subscriber: sub,
	}
}

func (p pubsub) Close() error {
	return p.client.Close()
}
