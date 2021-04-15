package internal

import (
	"github.com/edgexfoundry/go-mod-bootstrap/di"
	"github.com/inspii/edgex-thingsboard/internal/bootstrap/container"
	"github.com/inspii/edgex-thingsboard/internal/thingsboard"
	"github.com/inspii/edgex-thingsboard/internal/utils"
	"net/http"
	"time"
)

func mqttForwardHTTP(dic *di.Container) error {
	config := container.ConfigurationFrom(dic.Get)
	pubsub := container.PubSubFrom(dic.Get)
	serviceRoutes := container.ClientsFrom(dic.Get)

	return pubsub.Subscribe(config.Mqtt.RPCRequestTopic, func(topic string, msg []byte) error {
		req := &thingsboard.RPCRequestMessage{}
		resp := &thingsboard.RPCResponseMessage{}

		if err := req.FromBytes(msg); err == nil {
			if serviceAddr, ok := serviceRoutes.Get(req.Data.Service); ok {
				req.Data.Service = serviceAddr
			}

			resp = forward(req)
		} else {
			resp = &thingsboard.RPCResponseMessage{
				Data: thingsboard.RPCResponseData{
					HTTPStatus: http.StatusInternalServerError,
					Message:    err.Error(),
				},
			}
		}

		resp.ID = req.Data.ID
		resp.Device = req.Device
		return pubsub.Publish(config.Mqtt.RPCResponseTopic, resp.Bytes())
	})
}

func forward(req *thingsboard.RPCRequestMessage) *thingsboard.RPCResponseMessage {
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
		Data: thingsboard.RPCResponseData{
			HTTPStatus: httpStatus,
			Message:    errMsg,
			Result:     result,
		},
	}
}
