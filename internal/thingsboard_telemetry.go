package internal

import (
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"time"
)

func serveThingsboardTelemetry(dic *di.Container) error {
	thingsboardMQTTConfig := container.ConfigurationFrom(dic.Get).ThingsBoardMQTT
	client := container.MQTTFrom(dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(dic.Get)

	mqttTimeout := time.Duration(thingsboardMQTTConfig.Timeout) * time.Millisecond
	handler := newThingsboardRPCHandler(dic)
	server := mqtt_server.New(client, 0, mqttTimeout, logger)
	return server.HandleFunc(thingsboardMQTTConfig.RPCRequestTopic, thingsboardMQTTConfig.RPCResponseTopic, handler.handleRPC)
}