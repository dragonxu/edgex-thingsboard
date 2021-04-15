package interfaces

import "github.com/inspii/edgex-thingsboard/internal/bootstrap"

type MQTTInfo interface {
	GetMQTTInfo() bootstrap.MQTTInfo
}
