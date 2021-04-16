package container

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
)

var MQTTName = di.TypeInstanceToName((*mqtt.Client)(nil))

func MQTTFrom(get di.Get) mqtt.Client {
	return get(MQTTName).(mqtt.Client)
}
