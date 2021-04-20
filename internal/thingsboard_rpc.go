package internal

import (
	"errors"
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
	return server.HandleFunc(thingsboard.RPCTopic, thingsboard.RPCTopic, handler.handleRPC)
}

type thingsboardRPCHandler struct {
	dic *di.Container
}

func newThingsboardRPCHandler(dic *di.Container) *thingsboardRPCHandler {
	return &thingsboardRPCHandler{dic}
}

func (p *thingsboardRPCHandler) handleRPC(req []byte) (resp []byte, err error) {
	logger := bootstrapContainer.LoggingClientFrom(p.dic.Get)

	r := &thingsboard.RPCRequestMessage{}
	err = r.FromBytes(req)
	if err != nil {
		logger.Warn("bad rpc request: %s", req)
		return nil, err
	}
	if r.Data.Service == "" || r.Data.Method == "" || r.Data.URI == "" {
		logger.Warn("bad rpc request: %s", req)
		return nil, errors.New("bad rpc request")
	}

	return p.forwardHTTP(r).Bytes(), nil
}

func (p thingsboardRPCHandler) forwardHTTP(req *thingsboard.RPCRequestMessage) *thingsboard.RPCResponseMessage {
	serviceRoutes := container.ServiceRoutesFrom(p.dic.Get)
	serviceURL := req.Data.Service
	if serviceAddr, ok := serviceRoutes.Get(req.Data.Service); ok {
		serviceURL = serviceAddr
	}

	var result interface{}
	apiURL := serviceURL + req.Data.URI
	timeout := time.Duration(req.Data.Timeout) * time.Millisecond
	err := utils.RequestJSON(req.Data.Method, apiURL, timeout, req.Data.Params, &result)

	httpStatus := utils.GetHTTPStatus(err)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return &thingsboard.RPCResponseMessage{
		ID:     req.Data.ID,
		Device: req.Device,
		Data: thingsboard.RPCResponseData{
			Success:    err == nil,
			HTTPStatus: httpStatus,
			Message:    errMsg,
			Result:     result,
		},
	}
}
