package interfaces

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTOption interface {
	GetMQTTOption() *mqtt.ClientOptions
}
