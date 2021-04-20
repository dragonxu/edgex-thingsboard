package thingsboard

import "encoding/json"

const GatewayConnectTopic = "v1/gateway/connect"

type GatewayConnectMessage struct {
	DeviceName string `json:"device"`
}

func (m GatewayConnectMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
