package interfaces

import (
	msgTypes "github.com/edgexfoundry/go-mod-messaging/pkg/types"
)

type MessagingInfo interface {
	GetMessagingConfig() msgTypes.MessageBusConfig
}
