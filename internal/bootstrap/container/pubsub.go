package container

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/pkg/messaging"
)

var PubSubName = di.TypeInstanceToName((*messaging.PubSub)(nil))

func PubSubFrom(get di.Get) messaging.PubSub {
	return get(PubSubName).(messaging.PubSub)
}
