package internal

import (
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/device"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"github.com/inspii/edgex-thingsboard/internal/thingsboard"
	"time"
)

func connectGatewayDevices(dic *di.Container) error {
	thingsboardMQTTConfig := container.ConfigurationFrom(dic.Get).ThingsBoardMQTT
	client := container.MQTTFrom(dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(dic.Get)

	mqttTimeout := time.Duration(thingsboardMQTTConfig.Timeout) * time.Millisecond
	pubsub := mqtt_server.NewPubSubClient(client, 0, mqttTimeout, logger)

	devices, err := device.ListDevices(dic)
	if err != nil {
		logger.Error("list device names")
		return err
	}

	for _, d := range devices {
		m := thingsboard.ConnectMessage{
			DeviceName: d.Name,
		}
		if err := pubsub.Publish(thingsboard.ConnectTopic, m.Bytes()); err != nil {
			logger.Error(fmt.Sprintf("connect device: %s", err))
			continue
		}
		logger.Info(fmt.Sprintf("device %s connected", d.Name))
	}
	return nil
}
