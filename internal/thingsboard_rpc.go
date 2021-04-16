package internal

import (
	bootstrapContainer "github.com/edgexfoundry/go-mod-bootstrap/bootstrap/container"
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/mqtt_server"
	"github.com/inspii/edgex-thingsboard/internal/thingsboard"
	"github.com/inspii/edgex-thingsboard/internal/utils"
	"time"
)

func serveThingsboardRPC(dic *di.Container) error {
	thingsboardMQTTConfig := container.ConfigurationFrom(dic.Get).ThingsBoardMQTT
	client := container.MQTTFrom(dic.Get)
	logger := bootstrapContainer.LoggingClientFrom(dic.Get)

	mqttTimeout := time.Duration(thingsboardMQTTConfig.Timeout) * time.Millisecond
	handler := newThingsboardRPCHandler(dic)
	server := mqtt_server.New(client, 0, mqttTimeout, logger)
	return server.HandleFunc(thingsboardMQTTConfig.RPCRequestTopic, thingsboardMQTTConfig.RPCResponseTopic, handler.handleRPC)
}

type thingsboardRPCHandler struct {
	dic *di.Container
}

func newThingsboardRPCHandler(dic *di.Container) *thingsboardRPCHandler {
	return &thingsboardRPCHandler{dic}
}

func (p thingsboardRPCHandler) handleRPC(req []byte) (resp []byte, err error) {
	r := &thingsboard.RPCRequestMessage{}
	err = r.FromBytes(req)
	if err != nil {
		return nil, err
	}
	return p.forwardHTTP(r).Bytes(), nil
}

func (p thingsboardRPCHandler) forwardHTTP(req *thingsboard.RPCRequestMessage) *thingsboard.RPCResponseMessage {
	serviceRoutes := container.ServiceRoutesFrom(p.dic.Get)
	if serviceAddr, ok := serviceRoutes.Get(req.Data.Service); ok {
		req.Data.Service = serviceAddr
	}

	var result interface{}
	timeout := time.Duration(req.Data.Timeout) * time.Millisecond
	url := req.Data.Service + req.Data.URI
	err := utils.RequestJSON(req.Data.Method, url, timeout, req.Data.Params, &result)

	httpStatus := utils.GetHTTPStatus(err)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return &thingsboard.RPCResponseMessage{
		ID:     req.Data.ID,
		Device: req.Device,
		Data: thingsboard.RPCResponseData{
			HTTPStatus: httpStatus,
			Message:    errMsg,
			Result:     result,
		},
	}
}
