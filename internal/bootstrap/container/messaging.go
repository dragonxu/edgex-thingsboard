package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/edgexfoundry/go-mod-messaging/messaging"
)

var MessagingName = di.TypeInstanceToName((*messaging.MessageClient)(nil))

func MessagingFrom(get di.Get) messaging.MessageClient {
	return get(MessagingName).(messaging.MessageClient)
}
