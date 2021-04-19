package internal

import (
	"encoding/json"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	msgType "github.com/edgexfoundry/go-mod-messaging/pkg/types"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"github.com/inspii/edgex-thingsboard/internal/thingsboard"
	"time"
)

func serveThingsboardTelemetry(dic *di.Container) error {
	go forwardTelemetry(dic)
	return nil
}

func forwardTelemetry(dic *di.Container) {
	conf := container.ConfigurationFrom(dic.Get)
	messagingClient := container.MessagingFrom(dic.Get)
	thingsboardClient := container.MQTTFrom(dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(dic.Get)

	mqttTimeout := time.Duration(conf.ThingsBoardMQTT.Timeout) * time.Millisecond
	pubsubClient := mqtt_server.NewPubSubClient(thingsboardClient, 0, mqttTimeout, logger)

	errCh := make(chan error)
	defer close(errCh)
	msgCh := make(chan msgType.MessageEnvelope)
	defer close(msgCh)
	ch := msgType.TopicChannel{
		Topic:    conf.MessageQueue.Topic,
		Messages: msgCh,
	}
	if err := messagingClient.Subscribe([]msgType.TopicChannel{ch}, errCh); err != nil {
		logger.Error(fmt.Sprintf("subscribe message channel: %s", err))
	}

	for {
		select {
		case msg := <-msgCh:
			if msg.ContentType == "application/json" && len(msg.Payload) > 0 {
				event := &contract.Event{}
				if err := json.Unmarshal(msg.Payload, event); err != nil {
					logger.Error(fmt.Sprintf("unmarshal telemetry message error: %s", err))
					continue
				}

				msg, err := thingsboard.AdapterTelemetryMessage(event)
				if err != nil {
					logger.Error(fmt.Sprintf("adapter telemetry message error: %s", err))
					continue
				}

				if err := pubsubClient.Publish(conf.ThingsBoardMQTT.TelemetryTopic, msg.Bytes()); err != nil {
					logger.Error(fmt.Sprintf("publish telemetry message error: %s", err))
					continue
				}
			}
		case err := <-errCh:
			if err != nil {
				logger.Error(fmt.Sprintf("message error: %s", err))
			}
		}
	}
}
