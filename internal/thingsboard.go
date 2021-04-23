package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	msgType "github.com/edgexfoundry/go-mod-messaging/pkg/types"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/device"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"github.com/inspii/edgex-thingsboard/internal/thingsboard"
	"github.com/inspii/edgex-thingsboard/internal/utils"
	"time"
)

const (
	defaultAPITimeout = 60 * time.Second
)

type ThingsboardGateway struct {
	dic *di.Container
}

func NewThingsboardGateway(dic *di.Container) *ThingsboardGateway {
	return &ThingsboardGateway{dic}
}

func (t *ThingsboardGateway) Serve() error {
	lc := bootstrapContainer.LoggingClientFrom(t.dic.Get)

	if err := t.serveConnectDevices(); err != nil {
		lc.Error(fmt.Sprintf("connect thingsboard devices: %s", err.Error()))
		return err
	}
	if err := t.serveRPC(); err != nil {
		lc.Error(fmt.Sprintf("serve thingsboard rpc: %s", err.Error()))
		return err
	}
	if err := t.serveTelemetry(); err != nil {
		lc.Error(fmt.Sprintf("serve thingsboard telemetry: %s", err.Error()))
		return err
	}
	return nil
}

func (t *ThingsboardGateway) serveConnectDevices() error {
	thingsboardMQTTConfig := container.ConfigurationFrom(t.dic.Get).ThingsBoardMQTT
	client := container.MQTTFrom(t.dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(t.dic.Get)

	mqttTimeout := time.Duration(thingsboardMQTTConfig.Timeout) * time.Millisecond
	pubsub := mqtt_server.NewPubSubClient(client, 0, mqttTimeout, logger)

	devices, err := device.ListDevices(t.dic)
	if err != nil {
		logger.Error("list device names")
		return err
	}

	for _, d := range devices {
		m := thingsboard.GatewayConnectMessage{
			DeviceName: d.Name,
		}
		if err := pubsub.Publish(thingsboard.GatewayConnectTopic, m.Bytes()); err != nil {
			logger.Error(fmt.Sprintf("connect device: %s", err))
			continue
		}
		logger.Info(fmt.Sprintf("device %s connected", d.Name))
	}
	return nil
}

func (t *ThingsboardGateway) serveRPC() error {
	thingsboardMQTTConfig := container.ConfigurationFrom(t.dic.Get).ThingsBoardMQTT
	client := container.MQTTFrom(t.dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(t.dic.Get)

	mqttTimeout := time.Duration(thingsboardMQTTConfig.Timeout) * time.Millisecond
	server := mqtt_server.New(client, 0, mqttTimeout, logger)
	return server.HandleFunc(thingsboard.GatewayRPCTopic, thingsboard.GatewayRPCTopic, t.handleRPC)
}

func (t *ThingsboardGateway) handleRPC(req []byte) ([]byte, error) {
	logger := bootstrapContainer.LoggingClientFrom(t.dic.Get)
	logger.Debug(fmt.Sprintf("receive rpc message: %s", req))

	r := &thingsboard.GatewayRPCRequestMessage{}
	err := r.FromBytes(req)
	if err != nil {
		logger.Warn("bad rpc request: %s", req)
		return nil, err
	}
	if r.Data.Service == "" || r.Data.Method == "" || r.Data.URI == "" {
		logger.Warn("bad rpc request: %s", req)
		return nil, errors.New("bad rpc request")
	}

	resp := t.forwardHTTP(r).Bytes()
	logger.Debug(fmt.Sprintf("rpc message response: %s", resp))
	return resp, nil
}

func (t *ThingsboardGateway) forwardHTTP(req *thingsboard.GatewayRPCRequestMessage) *thingsboard.GatewayRPCResponseMessage {
	serviceRoutes := container.ServiceRoutesFrom(t.dic.Get)
	serviceURL := req.Data.Service
	if serviceAddr, ok := serviceRoutes.Get(req.Data.Service); ok {
		serviceURL = serviceAddr
	}

	var result interface{}
	apiURL := serviceURL + req.Data.URI

	timeout := defaultAPITimeout
	if req.Data.APITimeout > 0 {
		timeout = time.Duration(req.Data.APITimeout) * time.Millisecond
	}
	err := utils.RequestJSON(req.Data.Method, apiURL, timeout, req.Data.Params, &result)

	httpStatus := utils.GetHTTPStatus(err)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return &thingsboard.GatewayRPCResponseMessage{
		ID:     req.Data.ID,
		Device: req.Device,
		Data: thingsboard.GatewayRPCResponseData{
			Success:    err == nil,
			HTTPStatus: httpStatus,
			Message:    errMsg,
			Result:     result,
		},
	}
}

func (t *ThingsboardGateway) serveTelemetry() error {
	go t.forwardTelemetry()
	return nil
}

func (t *ThingsboardGateway) forwardTelemetry() {
	conf := container.ConfigurationFrom(t.dic.Get)
	messagingClient := container.MessagingFrom(t.dic.Get)
	thingsboardClient := container.MQTTFrom(t.dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(t.dic.Get)

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

				msg, err := thingsboard.AdapterGatewayTelemetryMessage(event)
				if err != nil {
					logger.Error(fmt.Sprintf("adapter telemetry message error: %s", err))
					continue
				}

				if err := pubsubClient.Publish(thingsboard.GatewayTelemetryTopic, msg.Bytes()); err != nil {
					logger.Error(fmt.Sprintf("publish telemetry message error: %s", err))
					continue
				}
				logger.Debug(fmt.Sprintf("telemetry reported: %s", msg.Bytes()))
			}
		case err := <-errCh:
			if err != nil {
				logger.Error(fmt.Sprintf("message error: %s", err))
			}
		}
	}
}
