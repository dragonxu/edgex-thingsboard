package thingsboard

import (
	"encoding/json"
)

const GatewayRPCTopic = "v1/gateway/rpc"

type GatewayRPCRequestData struct {
	ID      int         `json:"id"`
	Service string      `json:"service"`
	URI     string      `json:"uri"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"` // 毫秒
}

type GatewayRPCRequestMessage struct {
	Device string                `json:"device"`
	Data   GatewayRPCRequestData `json:"data"`
}

func (m *GatewayRPCRequestMessage) FromBytes(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m GatewayRPCRequestMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

type GatewayRPCResponseData struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"http_status"`
	Result     interface{} `json:"result"`
}

type GatewayRPCResponseMessage struct {
	ID     int                    `json:"id"`
	Device string                 `json:"device"`
	Data   GatewayRPCResponseData `json:"data"`
}

func (m *GatewayRPCResponseMessage) FromBytes(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m GatewayRPCResponseMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
