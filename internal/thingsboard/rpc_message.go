package thingsboard

import (
	"encoding/json"
)

const RPCTopic = "v1/gateway/rpc"

type RPCRequestData struct {
	ID      int         `json:"id"`
	Service string      `json:"service"`
	URI     string      `json:"uri"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Timeout int         `json:"timeout"` // 毫秒
}

type RPCRequestMessage struct {
	Device string         `json:"device"`
	Data   RPCRequestData `json:"data"`
}

func (m *RPCRequestMessage) FromBytes(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m RPCRequestMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

type RPCResponseData struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"http_status"`
	Result     interface{} `json:"result"`
}

type RPCResponseMessage struct {
	ID     int             `json:"id"`
	Device string          `json:"device"`
	Data   RPCResponseData `json:"data"`
}

func (m *RPCResponseMessage) FromBytes(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m RPCResponseMessage) Bytes() []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}
