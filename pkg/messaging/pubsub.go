package messaging

type Publisher interface {
	Publish(topic string, msg []byte) error
}

type MessageHandler func(topic string, msg []byte) error

type Subscriber interface {
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string) error
}

type PubSub interface {
	Publisher
	Subscriber
}
